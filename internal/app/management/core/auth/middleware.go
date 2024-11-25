package auth

import (
	"net/http"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
	"github.com/canonical/lxd/lxd/response"
	"github.com/gorilla/mux"
)

func AuthMiddleware(rc types.RouteConfig) mux.MiddlewareFunc {
	verifier := rc.Auth.(*Verifier)

	middlewareFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.TLS == nil {
				response.Forbidden(nil).Render(w, r)
				return
			}

			if verifier == nil {
				response.Forbidden(nil).Render(w, r)
				return
			}

			_, err := verifier.Auth(r.Context(), w, r)
			if err != nil {
				response.Forbidden(nil).Render(w, r)
				return
			}

			// If auth is successful, then we can proceed
			next.ServeHTTP(w, r)
		})
	}

	return middlewareFunc
}
