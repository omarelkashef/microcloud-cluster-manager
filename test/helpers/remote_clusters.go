package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/api/models/v1"
	"github.com/canonical/lxd/shared/api"
)

// FindRemoteCluster search for a remote cluster by name.
func FindRemoteCluster(env *Environment, remoteClusterName string) (*models.RemoteCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	certPublicKey, err := env.ManagementApiCert().PublicKeyX509()
	if err != nil {
		return nil, err
	}

	tlsClient, err := NewTLSHTTPClient(api.URL{}, nil, certPublicKey)
	if err != nil {
		return nil, err
	}

	output := &models.RemoteCluster{}
	path := api.NewURL().Scheme("https").Host(env.ManagementApiHost()).Path("1.0", "remote-cluster", remoteClusterName)
	err = tlsClient.Query(ctx, http.MethodGet, path, nil, output, nil)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteRemoteCluster deletes a remote cluster by name.
func DeleteRemoteCluster(env *Environment, remoteClusterName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	certPublicKey, err := env.ManagementApiCert().PublicKeyX509()
	if err != nil {
		return err
	}

	tlsClient, err := NewTLSHTTPClient(api.URL{}, nil, certPublicKey)
	if err != nil {
		return err
	}

	path := api.NewURL().Scheme("https").Host(env.ManagementApiHost()).Path("1.0", "remote-cluster", remoteClusterName)
	return tlsClient.Query(ctx, http.MethodDelete, path, nil, nil, nil)
}
