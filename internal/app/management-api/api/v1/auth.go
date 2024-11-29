package v1

import (
	"net/http"

	"github.com/canonical/lxd-cluster-manager/internal/app/management-api/core/auth"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
	"github.com/canonical/lxd/lxd/response"
	"github.com/google/uuid"
)

var Auth = types.RouteGroup{
	IsRoot: true,
	Prefix: "oidc",
	Endpoints: []types.Endpoint{
		{
			Path:    "login",
			Method:  http.MethodGet,
			Handler: login,
		},
		{
			Path:    "callback",
			Method:  http.MethodGet,
			Handler: callback,
		},
		{
			Path:    "logout",
			Method:  http.MethodGet,
			Handler: logout,
		},
	},
}

func login(rc types.RouteConfig) types.EndpointHandler {
	verifier := rc.Auth.(*auth.Verifier)
	return func(w http.ResponseWriter, r *http.Request) error {
		redirectURL := r.URL.Query().Get("next")

		stateToken := auth.StateToken{
			RedirectURL: redirectURL,
			ID:          uuid.New().String(),
		}

		state, err := stateToken.String()
		if err != nil {
			return response.InternalError(err).Render(w, r)
		}

		loginHandler := func(w http.ResponseWriter) error {
			verifier.Login(w, r, state)
			return nil
		}

		return response.ManualResponse(loginHandler).Render(w, r)
	}
}

func callback(rc types.RouteConfig) types.EndpointHandler {
	verifier := rc.Auth.(*auth.Verifier)
	return func(w http.ResponseWriter, r *http.Request) error {
		state := r.URL.Query().Get("state")
		stateToken, err := auth.DecodeStateToken(state)
		if err != nil {
			return response.InternalError(err).Render(w, r)
		}

		callbackHandler := func(w http.ResponseWriter) error {
			verifier.Callback(w, r, stateToken.RedirectURL)
			return nil
		}

		return response.ManualResponse(callbackHandler).Render(w, r)
	}
}

func logout(rc types.RouteConfig) types.EndpointHandler {
	verifier := rc.Auth.(*auth.Verifier)
	return func(w http.ResponseWriter, r *http.Request) error {
		redirectURL := r.URL.Query().Get("next")

		logoutHandler := func(w http.ResponseWriter) error {
			verifier.Logout(w, r, redirectURL)
			return nil
		}

		return response.ManualResponse(logoutHandler).Render(w, r)
	}
}
