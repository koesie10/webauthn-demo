package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var contextKeySession = "webauthn-demo-session"

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   2592000, // 30 days
			HttpOnly: true,
		}
		c.Set(contextKeySession, sess)

		c.Response().Before(func() {
			sess.Save(c.Request(), c.Response())
		})

		err := next(c)

		c.Set(contextKeySession, nil)

		return err
	}
}

func SessionFromContext(c echo.Context) *sessions.Session {
	sess, ok := c.Get(contextKeySession).(*sessions.Session)
	if !ok {
		return nil
	}
	return sess
}
