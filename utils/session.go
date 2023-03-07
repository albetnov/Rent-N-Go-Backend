package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SessionStore struct {
	store *session.Store
	c     *fiber.Ctx
}

var Session = SessionStore{}

// InitStore
// Make a new singleton instance of Store
func (s *SessionStore) InitStore() SessionStore {
	s.store = session.New()
	return *s
}

// Provide
// Set the current session store context to given context.
func (s *SessionStore) Provide(c *fiber.Ctx) SessionStore {
	s.c = c
	return *s
}

// SetSession
// Set a session and automatically save it
func (s SessionStore) SetSession(name string, value any) SessionStore {
	sess, err := s.store.Get(s.c)

	if err != nil {
		ShouldPanic(err)
	}

	sess.Set(name, value)

	if err = sess.Save(); err != nil {
		ShouldPanic(err)
	}

	return s
}

// DeleteSession
// Delete a session with given name and automatically save it
func (s SessionStore) DeleteSession(name string) SessionStore {
	sess, err := s.store.Get(s.c)

	if err != nil {
		ShouldPanic(err)
	}

	sess.Delete(name)

	if err = sess.Save(); err != nil {
		ShouldPanic(err)
	}

	return s
}

// GetSession
// Get a session from given name and return the value or nil if not exist.
func (s SessionStore) GetSession(name string) any {
	sess, err := s.store.Get(s.c)

	if err != nil {
		ShouldPanic(err)
	}

	return sess.Get(name)
}

// GetFlash
// Return the current session value, after return the session will be deleted.
func (s SessionStore) GetFlash(name string) any {
	sess, err := s.store.Get(s.c)

	if err != nil {
		ShouldPanic(err)
	}

	current := sess.Get(name)

	sess.Delete(name)

	defer sess.Save()

	return current
}
