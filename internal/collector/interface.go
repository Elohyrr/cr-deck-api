package collector

import (
	"context"
	"time"
)

// Service orchestrates the collection process
type Service interface {
	Collect(ctx context.Context) (*CollectResult, error)
}

// CollectResult contains statistics about a collection run
type CollectResult struct {
	PlayersProcessed int
	BattlesCollected int
	BattlesFiltered  int
	BattlesStored    int
	Errors           []error
	Duration         time.Duration
	StartedAt        time.Time
	CompletedAt      time.Time
}
