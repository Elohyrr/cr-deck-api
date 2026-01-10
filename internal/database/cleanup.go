package database

import (
	"context"
	"database/sql"
	"fmt"
)

// PurgeOldBattles removes battles older than the specified retention period
func PurgeOldBattles(ctx context.Context, db *sql.DB, retentionDays int) (int64, error) {
	query := `
		DELETE FROM battles 
		WHERE battle_time < NOW() - INTERVAL '1 day' * $1
	`

	result, err := db.ExecContext(ctx, query, retentionDays)
	if err != nil {
		return 0, fmt.Errorf("failed to purge old battles: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
