package main

import (
	"fmt"
	"os"

	"github.com/canonical/lxd-cluster-manager/cmd/admin"
	cluster_connector "github.com/canonical/lxd-cluster-manager/cmd/cluster-connector"
	management_api "github.com/canonical/lxd-cluster-manager/cmd/management-api"
	"github.com/canonical/lxd-cluster-manager/internal/pkg/logger"
)

var SERVICES = []string{
	"management-api",
	"cluster-connector",
	"admin",
}

func main() {
	// Perform the startup and shutdown sequence.
	err := run()
	if err != nil {
		logger.Log.Errorw("startup", "ERROR", err)
		logger.Log.Sync()
		os.Exit(1)
	}

	logger.Cleanup()
}

func run() error {
	// Get the service name from environment
	service := os.Getenv("SERVICE")
	if service == "" {
		return fmt.Errorf("service name is required, it should be one of: %v", SERVICES)
	}

	// Initialize the logger for the service
	logger.SetService(service)

	if service == "management-api" {
		return management_api.Run()
	}

	if service == "cluster-connector" {
		return cluster_connector.Run()
	}

	if service == "admin" {
		return admin.Run()
	}

	return fmt.Errorf("service name is invalid, it should be one of: %v", SERVICES)
}
