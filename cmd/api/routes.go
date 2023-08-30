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
	mux.HandlerFunc("GET", "/user/search", app.searchUser)

	mux.Handler("PUT", "/friend/set/:id", app.requireAuthenticatedUser(http.HandlerFunc(app.addFriend)))
	mux.Handler("PUT", "/friend/delete/:id", app.requireAuthenticatedUser(http.HandlerFunc(app.deleteFriend)))

	mux.Handler("POST", "/post/create", app.requireAuthenticatedUser(http.HandlerFunc(app.createPost)))
	mux.Handler("PUT", "/post/update", app.requireAuthenticatedUser(http.HandlerFunc(app.updatePost)))
	mux.Handler("PUT", "/post/delete/:id", app.requireAuthenticatedUser(http.HandlerFunc(app.deletePost)))
	mux.Handler("GET", "/post/get/:id", http.HandlerFunc(app.getPost))

	mux.Handler("GET", "/post/feed", app.requireAuthenticatedUser(http.HandlerFunc(app.getFeed)))

	mux.Handler("GET", "/protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))
	mux.Handler("GET", "/basic-auth-protected", app.requireBasicAuthentication(http.HandlerFunc(app.protected)))

	return app.recoverPanic(app.authenticate(mux))
}
