package database

import (
	"context"
	"shorter/internal/storage/database/repo"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	conn    *pgx.Conn
	URLRepo *repo.URLRepo
}

func (d *Database) Close(ctx context.Context) {
	d.conn.Close(ctx)
}

func NewDatabase(ctx context.Context, databasePath string) (*Database, error) {
	conn, err := pgx.Connect(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	urlRepo := repo.NewURLRepo(conn)
	if err := urlRepo.Init(ctx, conn); err != nil {
		return nil, err
	}
	return &Database{
		conn:    conn,
		URLRepo: urlRepo,
	}, nil
}
