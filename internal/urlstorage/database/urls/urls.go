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

// URLRepo is a struct that contains the necessary settings
type URLRepo struct {
	conn   *pgxpool.Pool
	logger *logger.ZapLogger
}

// NewURLRepo creates a new URL repository
func NewURLRepo(conn *pgxpool.Pool, l *logger.ZapLogger) *URLRepo {
	return &URLRepo{conn: conn, logger: l}
}

// SelectURLByShortURL is a query to select URL by short URL
const SelectURLByShortURL = "SELECT url, is_deleted FROM urls WHERE short_url = $1"

// GetURL returns the URL by short URL
func (d *URLRepo) GetURL(shortURL string) (string, error) {
	var url string
	var isDeleted bool
	err := d.conn.QueryRow(context.Background(), SelectURLByShortURL, shortURL).Scan(&url, &isDeleted)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", NewURLIsDeletedError(shortURL)
	}
	return url, nil
}

// SelectShortURLByURL is a query to select short URL by URL
const (
	SelectShortURLByURL = "SELECT short_url FROM urls WHERE url = $1"
)

// GetShortURL returns the short URL by URL
func (d *URLRepo) GetShortURL(url string) (string, error) {
	var shortURL string
	err := d.conn.QueryRow(context.Background(), SelectShortURLByURL, url).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// InsertUrls is a query to insert URLs
const InsertUrls = "INSERT INTO urls (short_url, url, user_id) VALUES ($1, $2, $3)"

// Save saves the URL
func (d *URLRepo) Save(values domain.URLS, userID int) error {
	_, err := d.conn.Exec(context.TODO(), InsertUrls, values.CorrelationID, values.URL, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			shortURL, sErr := d.GetShortURL(values.URL)
			if sErr != nil {
				return sErr
			}
			return NewDuplicateError(values.URL, shortURL)
		}
		return err
	}
	return nil
}

// SaveBatch saves the URL batch
func (d *URLRepo) SaveBatch(
	ctx context.Context, values []domain.URLS, userID int) error {
	batch := &pgx.Batch{}
	for _, value := range values {
		batch.Queue(InsertUrls, value.CorrelationID, value.URL, userID)
	}

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

// UpdateDeleteFlagInBatch is a query to update the delete flag in the batch
const UpdateDeleteFlagInBatch = "UPDATE urls SET is_deleted =true WHERE short_url = ANY($1) AND user_id = $2"

// DeleteURLs deletes the URL batch
func (d *URLRepo) DeleteURLs(correlationIDS []string,
	userID int,
) error {
	_, err := d.conn.Exec(context.TODO(), UpdateDeleteFlagInBatch, correlationIDS, userID)
	return fmt.Errorf("error when deleting urls: %w", err)
}

// UserURLs is a query to select user URLs
const UserURLs = "SELECT short_url, url FROM urls WHERE user_id = $1"

// GetUserURLs returns the user URLs
func (d *URLRepo) GetUserURLs(userID int, serverAdr string) ([]domain.UserURLs, error) {
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
		url.ShortURL = serverAdr + url.ShortURL
		urls = append(urls, url)
	}
	return urls, nil
}

// GetInitialData is not implemented because all data is in the database
func (d *URLRepo) GetInitialData() (domain.URLStorage, error) {
	d.logger.Log.Warn("GetInitialData is not implemented")
	return nil, fmt.Errorf("not implemented")
}
