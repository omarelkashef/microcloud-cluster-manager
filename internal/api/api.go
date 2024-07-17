package api

import (
	"net/http"

	"github.com/canonical/lxd/lxd/response"
	microState "github.com/canonical/microcluster/state"

	"github.com/canonical/lxd-site-manager/internal/api/types"
	"github.com/canonical/lxd-site-manager/internal/state"
)

func authHandler(siteManagerState *state.SiteManagerState) types.AccessHandler {
	return func(clusterState microState.State, r *http.Request) (bool, response.Response) {
		// always allow unix socket requests
		if r.RemoteAddr == "@" {
			return true, nil
		}

		if r.TLS == nil {
			return false, response.Forbidden(nil)
		}

		if siteManagerState.OIDCVerifier == nil {
			return false, response.Forbidden(nil)
		}

		_, resp, err := siteManagerState.OIDCVerifier.Auth(r.Context(), r)
		if err != nil {
			return false, resp
		}

		return true, resp
	}
}
