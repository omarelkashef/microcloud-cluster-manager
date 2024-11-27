package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/api/models"
	"github.com/canonical/lxd-cluster-manager/test/helpers"
	"github.com/canonical/lxd/shared"
	"github.com/canonical/lxd/shared/api"
)

func testRemoteClusterSuccess(e *helpers.Environment) (testName string, testFunc func(t *testing.T)) {
	return "lxd remote cluster join and status updates under normal conditions", func(t *testing.T) {
		remoteClusterName := "remote_cluster_control_e2e"
		var condition string
		var err error
		var tokenData models.RemoteClusterTokenBody

		{
			condition = "Should be able to create token with valid data"
			tokenData, err = helpers.CreateAndReturnRemoteClusterJoinToken(e, remoteClusterName, time.Time{})
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			if tokenData.ServerName == "" {
				err = fmt.Errorf("invalid server_name")
				helpers.LogTestOutcome(t, condition, err)
			}

			if tokenData.Fingerprint == "" {
				err = fmt.Errorf("invalid fingerprint")
				helpers.LogTestOutcome(t, condition, err)
			}

			if tokenData.Address != e.ControlHost() {
				err = fmt.Errorf("invalid address")
				helpers.LogTestOutcome(t, condition, err)
			}

			if tokenData.Secret == "" {
				err = fmt.Errorf("invalid secret")
				helpers.LogTestOutcome(t, condition, err)
			}

			if tokenData.ExpiresAt == (time.Time{}) {
				err = fmt.Errorf("invalid expiry")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}

		{
			condition = "Should be able to receive a join request"
			err = sendJoinRequest(e, tokenData)
			helpers.LogTestOutcome(t, condition, err)
		}

		{
			condition = "Should be able to get remote cluster with PENDING_APPROVAL status"
			remoteCluster, err := helpers.FindRemoteCluster(e, remoteClusterName)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			if remoteCluster.Status != string(models.PENDING_APPROVAL) {
				err = fmt.Errorf("invalid remote cluster status")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}

		{
			condition = "Should have deleted the remote cluster join token after receiving join request"
			token, err := helpers.FindToken(e, remoteClusterName)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			if token != (models.RemoteClusterToken{}) {
				err = fmt.Errorf("token not deleted")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}

		{
			condition = "Should be able to approve a join request"
			err = approveJoinRequest(e, remoteClusterName)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			remoteCluster, err := helpers.FindRemoteCluster(e, remoteClusterName)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			if remoteCluster.Status != string(models.ACTIVE) {
				err = fmt.Errorf("invalid remote cluster status")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}

		{
			condition = "Should be able to receive a status update"
			response, err := sendStatusUpdate(e, tokenData)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			expected := e.ControlHost()

			if !reflect.DeepEqual(response.ClusterManagerAddress, expected) {
				err = fmt.Errorf("invalid Cluster Manager address")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}

		{
			condition = "Should be able to get remote cluster status"
			remoteCluster, err := helpers.FindRemoteCluster(e, remoteClusterName)
			if err != nil {
				helpers.LogTestOutcome(t, condition, err)
			}

			if remoteCluster.CPUTotalCount != 8 {
				err = fmt.Errorf("invalid CPU total count")
				helpers.LogTestOutcome(t, condition, err)
			}

			if remoteCluster.MemoryTotalAmount != 1024 {
				err = fmt.Errorf("invalid memory total amount")
				helpers.LogTestOutcome(t, condition, err)
			}

			if remoteCluster.DiskTotalSize != 1024 {
				err = fmt.Errorf("invalid disk total size")
				helpers.LogTestOutcome(t, condition, err)
			}

			if !reflect.DeepEqual(remoteCluster.InstanceStatuses, []models.StatusDistribution{
				{Status: "running", Count: 1},
				{Status: "stopped", Count: 2},
			}) {
				err = fmt.Errorf("invalid instance statuses")
				helpers.LogTestOutcome(t, condition, err)
			}

			if !reflect.DeepEqual(remoteCluster.MemberStatuses, []models.StatusDistribution{
				{Status: "active", Count: 1},
				{Status: "inactive", Count: 2},
			}) {
				err = fmt.Errorf("invalid member statuses")
				helpers.LogTestOutcome(t, condition, err)
			}

			helpers.LogTestOutcome(t, condition, nil)
		}
	}
}

// sendJoinRequest sets up a tls http client and sends a join request to Cluster Manager.
func sendJoinRequest(env *helpers.Environment, tokenData models.RemoteClusterTokenBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// generate dedicated lxd cert for Cluster Manager
	clientCert, err := shared.KeyPairAndCA(env.CertDir(), tokenData.ServerName, shared.CertClient, shared.CertOptions{AddHosts: false})
	if err != nil {
		return err
	}

	controlCert := env.ControlCert()
	controlCertPublicKey, err := controlCert.PublicKeyX509()
	if err != nil {
		return err
	}

	tlsClient, err := helpers.NewTLSHTTPClient(api.URL{}, clientCert, controlCertPublicKey)
	if err != nil {
		return err
	}

	input := struct {
		ClusterName             string `json:"cluster_name"`
		RemotClusterCertificate string `json:"cluster_certificate"`
	}{
		ClusterName:             tokenData.ServerName,
		RemotClusterCertificate: string(clientCert.PublicKey()),
	}

	path := api.NewURL().Scheme("https").Host(tokenData.Address).Path("1.0", "remote-cluster")
	adjustHeaders := func(req *http.Request) error {
		mac := hmac.New(sha256.New, []byte(tokenData.Secret))
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return err
		}

		mac.Write(inputBytes)
		req.Header.Set("X-CLUSTER-SIGNATURE", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		return nil
	}

	return tlsClient.Query(ctx, http.MethodPost, path, input, nil, adjustHeaders)
}

// approveJoinRequest approves a join request for a remote cluster.
func approveJoinRequest(env *helpers.Environment, remoteClusterName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	managementCert := env.ManagementCert()
	managementCertPublicKey, err := managementCert.PublicKeyX509()
	if err != nil {
		return err
	}

	tlsClient, err := helpers.NewTLSHTTPClient(api.URL{}, nil, managementCertPublicKey)
	if err != nil {
		return err
	}

	input := models.RemoteClusterPatch{
		Status: models.ACTIVE,
	}

	path := api.NewURL().Scheme("https").Host(env.ManagementHost()).Path("1.0", "remote-cluster", remoteClusterName)
	return tlsClient.Query(ctx, http.MethodPatch, path, input, nil, nil)
}

// sendStatusUpdate sends a status update to the Cluster Manager with the correct client certificate.
func sendStatusUpdate(env *helpers.Environment, tokenData models.RemoteClusterTokenBody) (*models.RemoteClusterStatusPostResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientCert, err := shared.KeyPairAndCA(env.CertDir(), tokenData.ServerName, shared.CertClient, shared.CertOptions{AddHosts: false})
	if err != nil {
		return nil, err
	}

	controlCert := env.ControlCert()
	controlCertPublicKey, err := controlCert.PublicKeyX509()
	if err != nil {
		return nil, err
	}

	tlsClient, err := helpers.NewTLSHTTPClient(api.URL{}, clientCert, controlCertPublicKey)
	if err != nil {
		return nil, err
	}

	input := models.RemoteClusterStatusPost{
		CPUTotalCount:     8,
		CPULoad1:          "0.1",
		CPULoad5:          "0.2",
		CPULoad15:         "0.3",
		MemoryTotalAmount: 1024,
		MemoryUsage:       512,
		DiskTotalSize:     1024,
		DiskUsage:         512,
		InstanceStatuses: []models.StatusDistribution{
			{Status: "running", Count: 1},
			{Status: "stopped", Count: 2},
		},
		MemberStatuses: []models.StatusDistribution{
			{Status: "active", Count: 1},
			{Status: "inactive", Count: 2},
		},
	}

	var output models.RemoteClusterStatusPostResponse
	path := api.NewURL().Scheme("https").Host(tokenData.Address).Path("1.0", "remote-cluster", "status")
	err = tlsClient.Query(ctx, http.MethodPost, path, input, &output, nil)

	return &output, err
}
