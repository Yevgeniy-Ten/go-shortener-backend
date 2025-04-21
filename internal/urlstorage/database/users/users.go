package users

import (
	"context"
	"shorter/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UsersRepo is a struct that contains the necessary settings
type UsersRepo struct {
	conn   *pgxpool.Pool
	logger *logger.ZapLogger
}

// New creates a new user repository
const InsertUser = "INSERT INTO users DEFAULT VALUES RETURNING id"

// New creates a new user repository
func New(conn *pgxpool.Pool, l *logger.ZapLogger) *UsersRepo {
	return &UsersRepo{conn: conn, logger: l}
}

// Create creates a new user
func (u *UsersRepo) Create() (id int, err error) {
	err = u.conn.QueryRow(context.Background(), InsertUser).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
