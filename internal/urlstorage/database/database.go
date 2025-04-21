package database

import (
	"context"
	"shorter/internal/logger"
	"shorter/internal/urlstorage/database/urls"
	"shorter/internal/urlstorage/database/users"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Database is a struct that contains the necessary settings
type Database struct {
	conn      *pgxpool.Pool
	URLRepo   *urls.URLRepo
	UsersRepo *users.UsersRepo
}

// Close closes the connection to the database
func (d *Database) Close(_ context.Context) {
	d.conn.Close()
}

// New creates a new database connection
func New(ctx context.Context, l *logger.ZapLogger, databasePath string) (*Database, error) {
	conn, err := pgxpool.New(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	if err := Init(ctx, conn); err != nil {
		return nil, err
	}
	urlRepo := urls.NewURLRepo(conn, l)
	usersRepo := users.New(conn, l)

	l.Log.Info("Connected to database")
	return &Database{
		conn:      conn,
		URLRepo:   urlRepo,
		UsersRepo: usersRepo,
	}, nil
}

// Init initializes the database with migrations
func Init(_ context.Context, conn *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(conn)
	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
