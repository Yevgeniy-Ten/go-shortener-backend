package database

import (
	"context"
	"shorter/internal/domain"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	conn *pgx.Conn
}

func (d *Database) GetURL(shortURL string) (string, error) {
	var url string
	err := d.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE short_url = $1", shortURL).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
func (d *Database) Save(ctx context.Context, values domain.URLS) error {
	_, err := d.conn.Exec(ctx, "INSERT INTO urls (short_url, url) VALUES ($1, $2)", values.ShortURL, values.URL)

	if err != nil {
		return err
	}
	return nil
}
func (d *Database) Close(ctx context.Context) {
	d.conn.Close(ctx)
}

func NewConnection(ctx context.Context, databasePath string) (*Database, error) {
	conn, err := pgx.Connect(ctx, databasePath)
	if err != nil {
		return nil, err
	}
	if err := initDataBase(ctx, conn); err != nil {
		return nil, err
	}
	return &Database{
		conn: conn,
	}, nil
}
func initDataBase(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url TEXT NOT NULL,
			url TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
