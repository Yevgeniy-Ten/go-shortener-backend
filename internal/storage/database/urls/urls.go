package urls

import (
	"context"
	"errors"
	"fmt"
	"shorter/internal/domain"
	"shorter/internal/logger"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type URLRepo struct {
	conn   *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewURLRepo(conn *pgxpool.Pool, l *logger.ZapLogger) *URLRepo {
	return &URLRepo{conn: conn, logger: l}
}

func (d *URLRepo) GetURL(shortURL string) (string, error) {
	var url string
	err := d.conn.QueryRow(context.Background(), "SELECT url FROM urls WHERE short_url = $1", shortURL).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d *URLRepo) GetShortURL(url string) (string, error) {
	var shortURL string

	err := d.conn.QueryRow(context.Background(), "SELECT short_url FROM urls WHERE url = $1", url).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (d *URLRepo) Save(values domain.URLS) error {
	_, err := d.conn.Exec(context.TODO(), "INSERT INTO urls (short_url, url) VALUES ($1, $2)", values.URLId, values.URL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			shortURL, err := d.GetShortURL(values.URL)
			if err != nil {
				return err
			}
			return NewDuplicateError(values.URL, shortURL)
		}
		return err
	}
	return nil
}
func (d *URLRepo) SaveBatch(values []domain.URLS) error {
	ctx := context.Background()
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			d.logger.Log.Error("Error while Rollback transaction", zap.Error(err))
		} else {
			err = tx.Commit(ctx)
			d.logger.Log.Error("Error while committing transaction", zap.Error(err))
		}
	}()
	for _, value := range values {
		_, err = tx.Exec(context.TODO(), "INSERT INTO urls (short_url, url) VALUES ($1, $2)", value.URLId, value.URL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *URLRepo) Init(ctx context.Context, conn *pgxpool.Pool) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE
		);
	`)

	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS urls_short_url_idx ON urls (short_url);
	`)

	if err != nil {
		return err
	}
	return nil
}

func (d *URLRepo) GetInitialData() (domain.Storage, error) {
	d.logger.Log.Warn("GetInitialData is not implemented")
	return nil, fmt.Errorf("not implemented")
}
