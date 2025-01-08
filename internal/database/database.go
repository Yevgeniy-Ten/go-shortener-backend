package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewConnection(ctx context.Context, databasePath string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
