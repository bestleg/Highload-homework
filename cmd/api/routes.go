package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.HandlerFunc("GET", "/status", app.status)

	mux.HandlerFunc("POST", "/user/register", app.createUser)
	mux.HandlerFunc("POST", "/login", app.createAuthenticationToken)
	mux.HandlerFunc("GET", "/user/get/:id", app.getUserByUserID)

	mux.Handler("GET", "/protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))
	mux.Handler("GET", "/basic-auth-protected", app.requireBasicAuthentication(http.HandlerFunc(app.protected)))

	return app.recoverPanic(app.authenticate(mux))
}
