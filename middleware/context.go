package middleware

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

func AddContext(h http.Handler) http.Handler {
	ctx := context.Background()
	log := hclog.Default()
	log.Error("TEST MIDDLEWARE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
