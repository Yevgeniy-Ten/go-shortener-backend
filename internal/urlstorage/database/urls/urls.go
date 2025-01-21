package urls

import (
	"context"
	"errors"
	"fmt"
	"shorter/internal/domain"
	"shorter/internal/logger"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	SelectURLByShortURL = "SELECT url FROM urls WHERE short_url = $1"
	SelectShortURLByURL = "SELECT short_url FROM urls WHERE url = $1"
	InsertUrls          = "INSERT INTO urls (short_url, url, user_id) VALUES ($1, $2, $3)"
	UserURLs            = "SELECT short_url, url FROM urls WHERE user_id = $1"
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
	err := d.conn.QueryRow(context.Background(), SelectURLByShortURL, shortURL).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d *URLRepo) GetShortURL(url string) (string, error) {
	var shortURL string

	err := d.conn.QueryRow(context.Background(), SelectShortURLByURL, url).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (d *URLRepo) Save(values domain.URLS, userID int) error {
	_, err := d.conn.Exec(context.TODO(), InsertUrls, values.CorrelationID, values.URL, userID)
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
func (d *URLRepo) SaveBatch(values []domain.URLS, userID int) error {
	batch := &pgx.Batch{}
	for _, value := range values {
		batch.Queue(InsertUrls, value.CorrelationID, value.URL, userID)
	}

	ctx := context.TODO()

	br := d.conn.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			fmt.Printf("error closing batch: %v\n", err)
		}
	}()

	for range values {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *URLRepo) GetUserURLs(userID int) ([]domain.UserURLs, error) {
	rows, err := d.conn.Query(context.Background(), UserURLs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []domain.UserURLs
	for rows.Next() {
		var url domain.UserURLs
		err = rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			return nil, fmt.Errorf("error when getting user rows: %w", err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (d *URLRepo) GetInitialData() (domain.URLStorage, error) {
	d.logger.Log.Warn("GetInitialData is not implemented")
	return nil, fmt.Errorf("not implemented")
}
