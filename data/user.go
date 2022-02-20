package data

import "context"

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) CreateAccount(ctx context.Context, user User) error {
	return s.db.CreateAccount(ctx, user)
}

func (s *Service) IfEmailOrUsernameExists(ctx context.Context,
	credentialType string, credential string) (bool, error) {
	return s.db.IfEmailOrUsernameExists(ctx, credentialType, credential)
}

func (s *Service) GetHashByCredential(ctx context.Context,
	credentialType string, credential string) (string, error) {
	return s.db.GetHashByCredential(ctx, credentialType, credential)
}
