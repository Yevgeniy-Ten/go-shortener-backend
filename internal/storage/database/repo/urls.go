package repo

import (
	"context"
	"fmt"
	"shorter/internal/domain"

	"github.com/jackc/pgx/v5"
)

type URLRepo struct {
	conn *pgx.Conn
}

func NewURLRepo(conn *pgx.Conn) *URLRepo {
	return &URLRepo{conn: conn}
}

func (d *URLRepo) GetURL(shortURL string) (string, error) {
	var url string
	err := d.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE short_url = $1", shortURL).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d *URLRepo) Save(values domain.URLS) error {
	_, err := d.conn.Exec(context.TODO(), "INSERT INTO urls (short_url, url) VALUES ($1, $2)", values.URLId, values.URL)

	if err != nil {
		return err
	}
	return nil
}

func (d *URLRepo) Init(ctx context.Context, conn *pgx.Conn) error {
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

func (d *URLRepo) GetInitialData() (domain.Storage, error) {
	return nil, fmt.Errorf("not implemented")
}
