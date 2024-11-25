package api

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/canonical/lxd/lxd/response"
	"github.com/gorilla/mux"

	"github.com/canonical/lxd-cluster-manager/config"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/database"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/logger"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
)

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type ApiConfig struct {
	Shutdown  chan os.Signal
	DB        *database.DB
	Auth      types.Authenticator
	EnvConfig *config.Config
}

// Api is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type Api struct {
	mux       *mux.Router
	shutdown  chan os.Signal
	db        *database.DB
	auth      types.Authenticator
	envConfig *config.Config
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApi(cfg ApiConfig) *Api {
	mux := mux.NewRouter()
	mux.StrictSlash(false)
	mux.SkipClean(true)
	mux.UseEncodedPath()

	return &Api{
		mux:       mux,
		shutdown:  cfg.Shutdown,
		db:        cfg.DB,
		auth:      cfg.Auth,
		envConfig: cfg.EnvConfig,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *Api) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// UseGlobalMiddleWares adds global middlewares to the router.
func (a *Api) UseGlobalMiddleWares(mw ...mux.MiddlewareFunc) {
	a.mux.Use(mw...)
}

// RegisterRoutes adds the routes to the router.
func (a *Api) RegisterRoutes(routes []types.RouteGroup) {
	rc := types.RouteConfig{
		Auth: a.auth,
		DB:   a.db,
		Env:  a.envConfig,
	}

	registerRoutes(a.mux, routes, rc)

	a.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := response.SyncResponse(true, []string{"/1.0"}).Render(w, r)
		if err != nil {
			logger.Log.Errorw("Failed to write HTTP response", "url", r.URL, "err", err.Error())
		}
	})

	a.mux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Infow("Sending top level 404", "url", r.URL)
		w.Header().Set("Content-Type", "application/json")
		err := response.NotFound(nil).Render(w, r)
		if err != nil {
			logger.Log.Error("Failed to write HTTP response", "url", r.URL, "err", err.Error())
		}
	})
}

// Mux returns the router.
func (a *Api) Mux() *mux.Router {
	return a.mux
}

// ServeHTTP implements the http.Handler interface.
func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
