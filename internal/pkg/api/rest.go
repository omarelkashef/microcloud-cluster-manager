package api

import (
	"net/http"
	"path"

	"github.com/canonical/lxd/lxd/response"
	"github.com/gorilla/mux"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/logger"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
)

func registerRoutes(mux *mux.Router, routes []types.RouteGroup, rc types.RouteConfig) {
	// Register route groups
	for _, rg := range routes {
		registerRouteGroup(mux, rg, rc)
	}
}

func registerRouteGroup(mux *mux.Router, rg types.RouteGroup, rc types.RouteConfig) {
	routeGroupPath := path.Join("/", rg.Prefix)
	if !rg.IsRoot {
		routeGroupPath = path.Join("/", rc.Version, rg.Prefix)
	}

	// apply middlewares at route group level
	sr := mux.PathPrefix(routeGroupPath).Subrouter()
	if len(rg.Middlewares) > 0 {
		for _, m := range rg.Middlewares {
			sr.Use(m(rc))
		}
	}

	for _, e := range rg.Endpoints {
		registerEndpoint(sr, routeGroupPath, e, rc)
	}
}

func registerEndpoint(mux *mux.Router, prefix string, e types.Endpoint, rc types.RouteConfig) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := e.Handler(rc)(w, r)
		if err != nil {
			logger.Log.Errorw("internal error", "ERROR", err)
			renderErr := response.InternalError(err).Render(w, r)
			if renderErr != nil {
				logger.Log.Errorw("failed to write error response", "path", path.Join(prefix, e.Path), "ERROR", renderErr.Error())
			}
		}
	})

	// in case if the endpoint is the root of the route group
	ep := ""
	if e.Path != "" {
		ep = path.Join("/", e.Path)
	}

	mux.Handle(ep, handler).Methods(e.Method)
}
