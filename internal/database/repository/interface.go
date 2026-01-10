package repository

import (
	"context"

	"github.com/leopoldhub/royal-api-personal/internal/models"
)

// BattleRepository manages battle data persistence
type BattleRepository interface {
	Insert(ctx context.Context, battle *models.Battle) error
	BatchInsert(ctx context.Context, battles []*models.Battle) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
	Count(ctx context.Context) (int, error)
	GetRecent(ctx context.Context, deckSignature string, limit int) ([]*models.Battle, error)
}

// MetaDeckRepository manages aggregated deck statistics
type MetaDeckRepository interface {
	Recalculate(ctx context.Context) error
	GetTop(ctx context.Context, limit int, sortBy string, minGames int) ([]*models.Deck, error)
	GetBySignature(ctx context.Context, signature string) (*models.Deck, error)
	Count(ctx context.Context) (int, error)
}
