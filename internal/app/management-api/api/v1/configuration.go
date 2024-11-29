package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/canonical/lxd-cluster-manager/config"
	"github.com/canonical/lxd-cluster-manager/internal/app/management-api/core/auth"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/api/models/v1"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/database"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/types"
	"github.com/canonical/lxd/lxd/response"
)

// Configuration defines the API for the configuration endpoint.
var Configuration = types.RouteGroup{
	Prefix: "configuration",
	Middlewares: []types.RouteMiddleware{
		auth.AuthMiddleware,
	},
	Endpoints: []types.Endpoint{
		{
			Method:  http.MethodGet,
			Handler: configurationGet,
		},
	},
}

func configurationGet(rc types.RouteConfig) types.EndpointHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		configs := mapEnvToConfig(*rc.Env)
		return response.SyncResponse(true, configs).Render(w, r)
	}
}

// generateConnString generates a PostgreSQL connection string from a DBConfig struct.
func generateConnString(config database.DBConfig) string {
	sslMode := "require"
	if config.DBDisableTLS {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		sslMode,
	)
}

// mapEnvToConfig maps application environment values to config data exposed via the API.
func mapEnvToConfig(cfg config.Config) models.Configuration {
	return models.Configuration{
		ApiVersion: models.ConfigData{
			Value:       cfg.ApiVersion,
			Title:       "API Version",
			Description: "The version of the API being used.",
		},
		ClusterConnectorAddress: models.ConfigData{
			Value:       cfg.ClusterConnectorAddress,
			Title:       "Cluster Connector Address",
			Description: "The host address for the cluster-connector API.",
		},
		OIDCClientID: models.ConfigData{
			Value:       cfg.OIDCClientID,
			Title:       "OIDC Client ID",
			Description: "The OpenID Conenct client ID for the application.",
		},
		OIDCIssuer: models.ConfigData{
			Value:       cfg.OIDCIssuer,
			Title:       "OIDC Issuer",
			Description: "OpenID Connect Discovery URL for the provider.",
		},
		OIDCAudience: models.ConfigData{
			Value:       cfg.OIDCAudience,
			Title:       "OIDC Audience",
			Description: "Expected audience value for the application.",
		},
		DBConnectionString: models.ConfigData{
			Value:       generateConnString(cfg.DBConfig),
			Title:       "Database Connection",
			Description: "The connection string for the Postgres database.",
		},
		DBMaxIdleConns: models.ConfigData{
			Value:       strconv.Itoa(cfg.DBMaxIdleConns),
			Title:       "Database Max Idle Connections",
			Description: "The maximum number of idle connections in the database pool.",
		},
		DBMaxOpenConns: models.ConfigData{
			Value:       strconv.Itoa(cfg.DBMaxOpenConns),
			Title:       "Database Max Open Connections",
			Description: "The maximum number of open connections in the database pool.",
		},
	}
}
