// Package middleware contains the applications' HTTP endpoints and defines how they respond to client requests
package middleware

import (
	"net/http"

	"article-service/infrastructure/appctx"
)

// SetRequestID middleware attaches application relevant header variables to the request custom context
func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := appctx.GenerateRequestID
		ctx := r.Context()
		ctx = appctx.WithReqID(ctx, reqID())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
