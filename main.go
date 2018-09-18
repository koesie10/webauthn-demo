package main

import (
	"log"
	"net/http"

	"github.com/koesie10/webauthn/webauthn"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

var storage = &Storage{
	authenticators: make(map[string]*Authenticator),
	users:          make(map[string]*User),
}

func main() {
	// Create the webauthn authenticator
	w, err := webauthn.New(&webauthn.Config{
		RelyingPartyName:   "webauthn-demo",
		Debug:              true,
		AuthenticatorStore: storage,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create echo and set some settings
	e := echo.New()
	e.Debug = true
	e.HideBanner = true

	// Add logger and recover middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create the cookie store with an insecure key and use the middleware so sessions are saved
	store := sessions.NewCookieStore([]byte("thisisanunsecurecookiestorepassword"))
	e.Use(session.Middleware(store), SessionMiddleware)

	// Register the handlers
	e.GET("/", Index)

	e.POST("/webauthn/registration/start/:name", func(c echo.Context) error {
		name := c.Param("name")
		u, ok := storage.users[name]
		if !ok {
			u = &User{
				Name:           name,
				Authenticators: make(map[string]*Authenticator),
			}
			storage.users[name] = u
		}

		sess := SessionFromContext(c)

		w.StartRegistration(c.Request(), c.Response(), u, webauthn.WrapMap(sess.Values))
		return nil
	})

	e.POST("/webauthn/registration/finish/:name", func(c echo.Context) error {
		name := c.Param("name")
		u, ok := storage.users[name]
		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		sess := SessionFromContext(c)

		w.FinishRegistration(c.Request(), c.Response(), u, webauthn.WrapMap(sess.Values))
		return nil
	})

	e.POST("/webauthn/login/start/:name", func(c echo.Context) error {
		name := c.Param("name")
		u, ok := storage.users[name]

		sess := SessionFromContext(c)

		if ok {
			w.StartLogin(c.Request(), c.Response(), u, webauthn.WrapMap(sess.Values))
		} else {
			w.StartLogin(c.Request(), c.Response(), nil, webauthn.WrapMap(sess.Values))
		}
		return nil
	})

	e.POST("/webauthn/login/finish/:name", func(c echo.Context) error {
		name := c.Param("name")
		u, ok := storage.users[name]

		sess := SessionFromContext(c)

		var authenticator webauthn.Authenticator
		if ok {
			authenticator = w.FinishLogin(c.Request(), c.Response(), u, webauthn.WrapMap(sess.Values))
		} else {
			authenticator = w.FinishLogin(c.Request(), c.Response(), nil, webauthn.WrapMap(sess.Values))
		}
		if authenticator == nil {
			return nil
		}

		authr, ok := authenticator.(*Authenticator)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, authr.User)
	})

	// Start the server
	log.Fatal(e.Start(":9000"))
}
