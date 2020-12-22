package security

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrInternal = errors.New("internal error")
var ErrNoSuchUser = errors.New("no such user")
var ErrTokenNotFound = errors.New("token not found")
var ErrPhoneUsed = errors.New("phone already registered")
var ErrInvalidPassword = errors.New("invalid password")
var ErrTokenExpired = errors.New("token expired")

type Service struct {
	pool *pgxpool.Pool
}

type Token struct {
	Token string `json:"token"`
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}

func (s *Service) Auth(login string, password string) bool {
	err := s.pool.QueryRow(context.Background(),
		`SELECT login, password FROM managers WHERE login=$1 and password=$2`,
		login, password).Scan(&login, &password)
	if err == nil {
		return true
	} else {
		return false
	}
}
