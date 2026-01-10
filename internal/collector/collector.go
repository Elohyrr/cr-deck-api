package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/leopoldhub/royal-api-personal/internal/database/repository"
	"github.com/leopoldhub/royal-api-personal/internal/errors"
	"github.com/leopoldhub/royal-api-personal/internal/models"
	"github.com/leopoldhub/royal-api-personal/pkg/supercell"
)

// CollectorService implements the Service interface
type CollectorService struct {
	supercellClient supercell.Client
	battleRepo      repository.BattleRepository
	metaRepo        repository.MetaDeckRepository
	limit           int
	logger          *log.Logger
}

var _ Service = (*CollectorService)(nil)

// NewService creates a new collector service
func NewService(
	client supercell.Client,
	battleRepo repository.BattleRepository,
	metaRepo repository.MetaDeckRepository,
	limit int,
	logger *log.Logger,
) Service {
	if logger == nil {
		logger = log.Default()
	}
	return &CollectorService{
		supercellClient: client,
		battleRepo:      battleRepo,
		metaRepo:        metaRepo,
		limit:           limit,
		logger:          logger,
	}
}

// Collect performs the full collection process
func (c *CollectorService) Collect(ctx context.Context) (*CollectResult, error) {
	result := &CollectResult{
		StartedAt: time.Now(),
		Errors:    make([]error, 0),
	}

	c.logger.Printf("Starting collection for top %d players", c.limit)

	players, err := c.supercellClient.GetTopPlayers(ctx, c.limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top players: %w", err)
	}

	c.logger.Printf("Fetched %d top players", len(players))

	battles := c.fetchBattlelogsParallel(ctx, players, result)
	c.logger.Printf("Collected %d raw battles from %d players", len(battles), result.PlayersProcessed)

	filtered := FilterPvPLadder(battles)
	result.BattlesCollected = len(battles)
	result.BattlesFiltered = len(filtered)
	c.logger.Printf("Filtered to %d PvP Ladder battles", len(filtered))

	parsed := c.parseBattles(filtered)
	c.logger.Printf("Parsed %d valid battles", len(parsed))

	if err := c.battleRepo.BatchInsert(ctx, parsed); err != nil {
		return nil, fmt.Errorf("failed to insert battles: %w", err)
	}
	result.BattlesStored = len(parsed)
	c.logger.Printf("Stored %d battles in database", result.BattlesStored)

	if err := c.metaRepo.Recalculate(ctx); err != nil {
		return nil, fmt.Errorf("failed to recalculate meta stats: %w", err)
	}
	c.logger.Println("Recalculated meta deck statistics")

	deleted, err := c.battleRepo.DeleteOlderThan(ctx, 7)
	if err != nil {
		c.logger.Printf("Warning: failed to purge old battles: %v", err)
	} else {
		c.logger.Printf("Purged %d old battles (7+ days)", deleted)
	}

	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	c.logger.Printf("Collection completed in %v", result.Duration)

	return result, nil
}

// fetchBattlelogsParallel fetches battlelogs using a worker pool
func (c *CollectorService) fetchBattlelogsParallel(
	ctx context.Context,
	players []supercell.Player,
	result *CollectResult,
) []supercell.BattleRaw {
	const numWorkers = 10

	jobs := make(chan supercell.Player, len(players))
	results := make(chan battlelogResult, len(players))

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for player := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					battlelog, err := c.supercellClient.GetBattlelog(ctx, player.Tag)
					results <- battlelogResult{
						PlayerTag: player.Tag,
						Battles:   battlelog,
						Error:     err,
					}
				}
			}
		}()
	}

	for _, player := range players {
		jobs <- player
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var allBattles []supercell.BattleRaw
	processed := 0
	for res := range results {
		processed++
		if processed%100 == 0 {
			c.logger.Printf("Progress: %d/%d players processed", processed, len(players))
		}

		if res.Error != nil {
			if apiErr, ok := res.Error.(*errors.APIError); ok {
				if apiErr.IsNotFound() {
					continue
				}
			}
			result.Errors = append(result.Errors, res.Error)
			continue
		}

		allBattles = append(allBattles, res.Battles...)
	}

	result.PlayersProcessed = processed

	return allBattles
}

// parseBattles converts raw battles to internal models
func (c *CollectorService) parseBattles(battles []supercell.BattleRaw) []*models.Battle {
	parsed := make([]*models.Battle, 0, len(battles))
	for _, raw := range battles {
		battle, err := ParseBattle(raw)
		if err != nil {
			continue
		}
		if battle != nil {
			parsed = append(parsed, battle)
		}
	}
	return parsed
}

type battlelogResult struct {
	PlayerTag string
	Battles   []supercell.BattleRaw
	Error     error
}
