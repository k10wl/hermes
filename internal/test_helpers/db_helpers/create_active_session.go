package db_helpers

import (
	"context"
	"database/sql"

	"github.com/k10wl/hermes/internal/models"
)

func CreateActiveSession(
	db *sql.DB,
	ctx context.Context,
	activeSession *models.ActiveSession,
) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO active_sessions (address, database_dns) VALUES (?, ?);`,
		activeSession.Address,
		activeSession.DatabaseDNS,
	)
	return err
}
