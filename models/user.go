package models

import "context"

func NewService(db DB) *Service {
	return &Service{db: db}
}

type Service struct {
	db DB
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DB interface {
	CreateAccount(ctx context.Context, user User) error
}

func (s *Service) CreateAccount(ctx context.Context, user User) (error error) {
	return s.db.CreateAccount(ctx, user)
}
