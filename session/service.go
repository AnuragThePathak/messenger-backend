package session

import "github.com/gorilla/sessions"

type Service struct {
	store *sessions.CookieStore
}

func NewSessionService(store *sessions.CookieStore) *Service {
	return &Service{store: store}
}
