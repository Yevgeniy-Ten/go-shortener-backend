package database

import (
	"context"
	"shorter/internal/logger"
	"shorter/internal/storage/database/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	conn    *pgxpool.Pool
	URLRepo *repo.URLRepo
}

func (d *Database) Close(ctx context.Context) {
	d.conn.Close()
}

func NewDatabase(ctx context.Context, l *logger.ZapLogger, databasePath string) (*Database, error) {
	conn, err := pgxpool.New(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	urlRepo := repo.NewURLRepo(conn, l)
	if err := urlRepo.Init(ctx, conn); err != nil {
		return nil, err
	}
	return &Database{
		conn:    conn,
		URLRepo: urlRepo,
	}, nil
}
