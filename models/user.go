package models

import "context"

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) CreateAccount(ctx context.Context, user User) (err error) {
	return s.db.CreateAccount(ctx, user)
}

func (s *Service) IfEmailExists(ctx context.Context, email string) (bool, error) {
	return s.db.IfEmailExists(ctx, email)
}
