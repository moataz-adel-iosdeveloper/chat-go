package config

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Store interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	New(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error
}

type CookieStore struct {
	store *sessions.CookieStore
}

func NewCookieStore(secret []byte) Store {
	cs := sessions.NewCookieStore(secret)

	cs.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return &CookieStore{store: cs}
}

func (c *CookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return c.store.Get(r, name)
}

func (c *CookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return c.store.New(r, name)
}

func (c *CookieStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return s.Save(r, w)
}
