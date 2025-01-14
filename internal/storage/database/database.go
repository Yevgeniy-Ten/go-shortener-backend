package database

import (
	"context"
	"shorter/internal/logger"
	"shorter/internal/storage/database/urls"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	conn    *pgxpool.Pool
	URLRepo *urls.URLRepo
}

func (d *Database) Close(ctx context.Context) {
	d.conn.Close()
}

func New(ctx context.Context, l *logger.ZapLogger, databasePath string) (*Database, error) {
	conn, err := pgxpool.New(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	urlRepo := urls.NewURLRepo(conn, l)
	if err := urlRepo.Init(ctx, conn); err != nil {
		return nil, err
	}
	l.Log.Info("Connected to database")
	return &Database{
		conn:    conn,
		URLRepo: urlRepo,
	}, nil
}
