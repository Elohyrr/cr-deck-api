package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/leopoldhub/royal-api-personal/internal/errors"
	"github.com/leopoldhub/royal-api-personal/internal/models"
)

type PostgresMetaRepo struct {
	db *sql.DB
}

var _ MetaDeckRepository = (*PostgresMetaRepo)(nil)

func NewMetaDeckRepository(db *sql.DB) MetaDeckRepository {
	return &PostgresMetaRepo{db: db}
}

func (r *PostgresMetaRepo) Recalculate(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &errors.DBError{
			Operation: "begin_transaction",
			Table:     "meta_decks",
			Err:       err,
		}
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM meta_decks"); err != nil {
		return &errors.DBError{
			Operation: "delete_all",
			Table:     "meta_decks",
			Err:       err,
		}
	}

	query := `
		INSERT INTO meta_decks (
			deck_signature, cards, total_games, wins, losses, 
			win_rate, first_seen, last_seen
		)
		SELECT 
			deck_signature,
			(array_agg(deck_cards ORDER BY battle_time DESC))[1] as cards,
			COUNT(*) as total_games,
			SUM(CASE WHEN is_victory THEN 1 ELSE 0 END) as wins,
			SUM(CASE WHEN NOT is_victory THEN 1 ELSE 0 END) as losses,
			ROUND((SUM(CASE WHEN is_victory THEN 1 ELSE 0 END)::numeric / COUNT(*)::numeric) * 100, 2) as win_rate,
			MIN(battle_time) as first_seen,
			MAX(battle_time) as last_seen
		FROM battles
		GROUP BY deck_signature
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return &errors.DBError{
			Operation: "recalculate",
			Table:     "meta_decks",
			Err:       err,
		}
	}

	return tx.Commit()
}

func (r *PostgresMetaRepo) GetTop(ctx context.Context, limit int, sortBy string, minGames int) ([]*models.Deck, error) {
	var orderColumn string
	switch sortBy {
	case "frequency":
		orderColumn = "total_games"
	case "win_rate":
		orderColumn = "win_rate"
	default:
		orderColumn = "win_rate"
	}

	query := fmt.Sprintf(`
		SELECT deck_signature, cards, total_games, wins, losses, 
			   win_rate, first_seen, last_seen, updated_at
		FROM meta_decks
		WHERE total_games >= $1
		ORDER BY %s DESC
		LIMIT $2
	`, orderColumn)

	rows, err := r.db.QueryContext(ctx, query, minGames, limit)
	if err != nil {
		return nil, &errors.DBError{
			Operation: "query_top",
			Table:     "meta_decks",
			Err:       err,
		}
	}
	defer rows.Close()

	var decks []*models.Deck
	for rows.Next() {
		var deck models.Deck
		var cardsJSON []byte
		deck.Stats = &models.DeckStats{}

		err := rows.Scan(
			&deck.Signature,
			&cardsJSON,
			&deck.Stats.TotalGames,
			&deck.Stats.Wins,
			&deck.Stats.Losses,
			&deck.Stats.WinRate,
			&deck.Stats.FirstSeen,
			&deck.Stats.LastSeen,
			&deck.Stats.UpdatedAt,
		)
		if err != nil {
			return nil, &errors.DBError{
				Operation: "scan_row",
				Table:     "meta_decks",
				Err:       err,
			}
		}

		if err := json.Unmarshal(cardsJSON, &deck.Cards); err != nil {
			continue
		}

		decks = append(decks, &deck)
	}

	return decks, nil
}

func (r *PostgresMetaRepo) GetBySignature(ctx context.Context, signature string) (*models.Deck, error) {
	query := `
		SELECT deck_signature, cards, total_games, wins, losses, 
			   win_rate, first_seen, last_seen, updated_at
		FROM meta_decks
		WHERE deck_signature = $1
	`

	var deck models.Deck
	var cardsJSON []byte
	deck.Stats = &models.DeckStats{}

	err := r.db.QueryRowContext(ctx, query, signature).Scan(
		&deck.Signature,
		&cardsJSON,
		&deck.Stats.TotalGames,
		&deck.Stats.Wins,
		&deck.Stats.Losses,
		&deck.Stats.WinRate,
		&deck.Stats.FirstSeen,
		&deck.Stats.LastSeen,
		&deck.Stats.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, &errors.DBError{
			Operation: "query_by_signature",
			Table:     "meta_decks",
			Err:       err,
		}
	}

	if err := json.Unmarshal(cardsJSON, &deck.Cards); err != nil {
		return nil, &errors.DBError{
			Operation: "unmarshal_cards",
			Table:     "meta_decks",
			Err:       err,
		}
	}

	return &deck, nil
}

func (r *PostgresMetaRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM meta_decks").Scan(&count)
	if err != nil {
		return 0, &errors.DBError{
			Operation: "count",
			Table:     "meta_decks",
			Err:       err,
		}
	}
	return count, nil
}
