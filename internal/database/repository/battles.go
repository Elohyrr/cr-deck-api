package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/leopoldhub/royal-api-personal/internal/errors"
	"github.com/leopoldhub/royal-api-personal/internal/models"
)

type PostgresBattleRepo struct {
	db *sql.DB
}

var _ BattleRepository = (*PostgresBattleRepo)(nil)

func NewBattleRepository(db *sql.DB) BattleRepository {
	return &PostgresBattleRepo{db: db}
}

func (r *PostgresBattleRepo) Insert(ctx context.Context, battle *models.Battle) error {
	cardsJSON, err := json.Marshal(battle.DeckCards)
	if err != nil {
		return &errors.DBError{
			Operation: "marshal_cards",
			Table:     "battles",
			Err:       err,
		}
	}

	query := `
		INSERT INTO battles (
			battle_time, player_tag, opponent_tag, game_mode,
			player_crowns, opponent_crowns, deck_signature,
			deck_cards, is_victory
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (player_tag, battle_time) DO NOTHING
	`

	_, err = r.db.ExecContext(ctx, query,
		battle.BattleTime,
		battle.PlayerTag,
		battle.OpponentTag,
		battle.GameMode,
		battle.PlayerCrowns,
		battle.OpponentCrowns,
		battle.DeckSignature,
		cardsJSON,
		battle.IsVictory,
	)

	if err != nil {
		return &errors.DBError{
			Operation: "insert",
			Table:     "battles",
			Err:       err,
		}
	}

	return nil
}

func (r *PostgresBattleRepo) BatchInsert(ctx context.Context, battles []*models.Battle) error {
	if len(battles) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &errors.DBError{
			Operation: "begin_transaction",
			Table:     "battles",
			Err:       err,
		}
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO battles (
			battle_time, player_tag, opponent_tag, game_mode,
			player_crowns, opponent_crowns, deck_signature,
			deck_cards, is_victory
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (player_tag, battle_time) DO NOTHING
	`)
	if err != nil {
		return &errors.DBError{
			Operation: "prepare_statement",
			Table:     "battles",
			Err:       err,
		}
	}
	defer stmt.Close()

	const batchSize = 100
	for i, battle := range battles {
		cardsJSON, err := json.Marshal(battle.DeckCards)
		if err != nil {
			continue
		}

		_, err = stmt.ExecContext(ctx,
			battle.BattleTime,
			battle.PlayerTag,
			battle.OpponentTag,
			battle.GameMode,
			battle.PlayerCrowns,
			battle.OpponentCrowns,
			battle.DeckSignature,
			cardsJSON,
			battle.IsVictory,
		)
		if err != nil {
			return &errors.DBError{
				Operation: "exec_insert",
				Table:     "battles",
				Err:       err,
			}
		}

		if (i+1)%batchSize == 0 {
			if err := tx.Commit(); err != nil {
				return &errors.DBError{
					Operation: "commit_batch",
					Table:     "battles",
					Err:       err,
				}
			}
			tx, _ = r.db.BeginTx(ctx, nil)
			stmt, _ = tx.PrepareContext(ctx, `
				INSERT INTO battles (
					battle_time, player_tag, opponent_tag, game_mode,
					player_crowns, opponent_crowns, deck_signature,
					deck_cards, is_victory
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (player_tag, battle_time) DO NOTHING
			`)
		}
	}

	return tx.Commit()
}

func (r *PostgresBattleRepo) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	query := `
		DELETE FROM battles 
		WHERE battle_time < NOW() - INTERVAL '1 day' * $1
	`

	result, err := r.db.ExecContext(ctx, query, days)
	if err != nil {
		return 0, &errors.DBError{
			Operation: "delete_old",
			Table:     "battles",
			Err:       err,
		}
	}

	return result.RowsAffected()
}

func (r *PostgresBattleRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM battles").Scan(&count)
	if err != nil {
		return 0, &errors.DBError{
			Operation: "count",
			Table:     "battles",
			Err:       err,
		}
	}
	return count, nil
}

func (r *PostgresBattleRepo) GetRecent(ctx context.Context, deckSignature string, limit int) ([]*models.Battle, error) {
	query := `
		SELECT id, battle_time, player_tag, opponent_tag, game_mode,
			   player_crowns, opponent_crowns, deck_signature, deck_cards, is_victory
		FROM battles
		WHERE deck_signature = $1
		ORDER BY battle_time DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, deckSignature, limit)
	if err != nil {
		return nil, &errors.DBError{
			Operation: "query_recent",
			Table:     "battles",
			Err:       err,
		}
	}
	defer rows.Close()

	var battles []*models.Battle
	for rows.Next() {
		var battle models.Battle
		var cardsJSON []byte

		err := rows.Scan(
			&battle.ID,
			&battle.BattleTime,
			&battle.PlayerTag,
			&battle.OpponentTag,
			&battle.GameMode,
			&battle.PlayerCrowns,
			&battle.OpponentCrowns,
			&battle.DeckSignature,
			&cardsJSON,
			&battle.IsVictory,
		)
		if err != nil {
			return nil, &errors.DBError{
				Operation: "scan_row",
				Table:     "battles",
				Err:       err,
			}
		}

		if err := json.Unmarshal(cardsJSON, &battle.DeckCards); err != nil {
			continue
		}

		battles = append(battles, &battle)
	}

	return battles, nil
}
