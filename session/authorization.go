package session

import "net/http"

func (s *Service) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session, err := s.store.Get(r, "session-name"); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if session.IsNew {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
