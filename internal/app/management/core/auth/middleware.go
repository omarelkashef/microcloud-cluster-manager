package auth

import (
	"net/http"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
	"github.com/canonical/lxd/lxd/response"
	"github.com/gorilla/mux"
)

// AuthMiddleware is a middleware function that checks if the request is authenticated.
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

			isTestMode := rc.Env.TestMode
			_, err := verifier.Auth(r.Context(), w, r)
			// NOTE: bypass oidc auth if we are running tests
			if err != nil && !isTestMode {
				response.Forbidden(nil).Render(w, r)
				return
			}

			// If auth is successful, then we can proceed
			next.ServeHTTP(w, r)
		})
	}

	return middlewareFunc
}
