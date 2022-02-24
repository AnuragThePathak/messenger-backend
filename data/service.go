package data

import (
	"context"
)

func NewService(db DB) *Service {
	return &Service{db: db}
}

type Service struct {
	db DB
}

type DB interface {
	CreateAccount(ctx context.Context, user User) error

	IfEmailOrUsernameExists(ctx context.Context, credentialType string,
		credential string) (bool, error)

	GetHashByCredential(ctx context.Context, credentialType string,
		credential string) (string, error)
}
