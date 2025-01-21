package users

import (
	"context"
	"shorter/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	conn   *pgxpool.Pool
	logger *logger.ZapLogger
}

const InsertUser = "INSERT INTO users DEFAULT VALUES RETURNING id"

func New(conn *pgxpool.Pool, l *logger.ZapLogger) *UsersRepo {
	return &UsersRepo{conn: conn, logger: l}
}

func (u *UsersRepo) Create() (id int, err error) {
	err = u.conn.QueryRow(context.Background(), InsertUser).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
