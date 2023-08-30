package main

import (
	"context"
	"net/http"

	"otus-homework/internal/database"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user *database.UserAuth) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *database.UserAuth {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*database.UserAuth)
	if !ok {
		return nil
	}

	return user
}
