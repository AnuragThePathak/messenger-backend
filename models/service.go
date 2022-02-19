package models

import "context"

func NewService(db DB) *Service {
	return &Service{db: db}
}

type Service struct {
	db DB
}

type DB interface {
	CreateAccount(ctx context.Context, user User) error

	IfEmailExists(ctx context.Context, email string) (bool, error)
}