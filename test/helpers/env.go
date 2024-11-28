package helpers

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/canonical/lxd/shared"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	// Succeed is a checkmark symbol.
	Succeed = "\u2713"
	// Failed is a cross symbol.
	Failed = "\u2717"
)

// Environment represents the environment for the tests.
type Environment struct {
	rootDir             string
	testDir             string
	certDir             string
	processIDs          []int
	kClient             *kubernetes.Clientset
	managementCert      *shared.CertInfo
	controlCert         *shared.CertInfo
	managementHost      string
	controlHost         string
	remoteClusters      []string
	remoteClusterTokens []string
}

// NewEnv creates a new environment.
func NewEnv() *Environment {
	return &Environment{
		rootDir:        getProjectRoot(),
		testDir:        "",
		certDir:        "",
		managementHost: "localhost:9000",
		controlHost:    "localhost:9001",
	}
}

// RootDir returns the root directory of the project.
func (e *Environment) RootDir() string {
	return e.rootDir
}

// TestDir returns the test directory.
func (e *Environment) DataDir() string {
	return e.rootDir + "/test/e2e/data"
}

// CertDir returns the certificate directory.
func (e *Environment) CertDir() string {
	return e.rootDir + "/test/e2e/data/certs"
}

// AddProcessID adds a process ID to the environment.
func (e *Environment) AddProcessID(pid int) {
	e.processIDs = append(e.processIDs, pid)
}

// ProcessIDs returns the pids for all processes started during a test run.
func (e *Environment) ProcessIDs() []int {
	return e.processIDs
}

// AddRemoteCluster mark a remote cluster for removal during test cleanup.
func (e *Environment) RemoveRemoteCluster(name string) {
	e.remoteClusters = append(e.remoteClusters, name)
}

// AddRemoteClusterToken mark a remote cluster join token for removal during test cleanup.
func (e *Environment) RemoveRemoteClusterToken(token string) {
	e.remoteClusterTokens = append(e.remoteClusterTokens, token)
}

// Init initializes the environment.
func (e *Environment) Init() error {
	// cleanup data dir in case if test didn't exit properly last time
	if err := os.RemoveAll(e.DataDir()); err != nil {
		return err
	}

	err := e.setKubeClient()
	if err != nil {
		return err
	}

	err = e.setCertificates()
	if err != nil {
		return err
	}

	err = e.setTestMode("true")
	if err != nil {
		return err
	}

	return nil
}

// Cleanup cleans up the environment.
func (e *Environment) Cleanup() error {
	pids := e.ProcessIDs()
	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", strconv.Itoa(pid))
		if err := killCmd.Run(); err != nil {
			continue
		}
	}

	if err := os.RemoveAll(e.DataDir()); err != nil {
		return err
	}

	// remove all new clusters and tokens created during the test
	for _, cluster := range e.remoteClusters {
		err := DeleteRemoteCluster(e, cluster)
		if err != nil {
			continue
		}
	}

	for _, token := range e.remoteClusterTokens {
		err := DeleteRemoteClusterJoinToken(e, token)
		if err != nil {
			continue
		}
	}

	// Reset the test mode
	err := e.setTestMode("false")
	if err != nil {
		return err
	}

	return nil
}

// ManagementCert returns the management certificate.
func (e *Environment) ManagementCert() *shared.CertInfo {
	return e.managementCert
}

// ControlCert returns the control certificate.
func (e *Environment) ControlCert() *shared.CertInfo {
	return e.controlCert
}

// ManagementHost returns the management host.
func (e *Environment) ManagementHost() string {
	return e.managementHost
}

// ControlHost returns the control host.
func (e *Environment) ControlHost() string {
	return e.controlHost
}

func (e *Environment) setTestMode(val string) error {
	deploymentName := "management-depl"
	containerName := "management"

	// Get the Deployment
	deployment, err := e.kClient.AppsV1().Deployments("default").Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	// Find the container in the deployment's spec
	var container *corev1.Container
	for i := range deployment.Spec.Template.Spec.Containers {
		if deployment.Spec.Template.Spec.Containers[i].Name == containerName {
			container = &deployment.Spec.Template.Spec.Containers[i]
			break
		}
	}

	if container == nil {
		return fmt.Errorf("container %s not found in deployment %s", containerName, deploymentName)
	}

	// Update the environment variable in the container
	updated := false
	for i := range container.Env {
		if container.Env[i].Name == "TEST_MODE" {
			container.Env[i].Value = val
			updated = true
			break
		}
	}

	if !updated {
		// Add the environment variable if it doesn't exist
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "TEST_MODE",
			Value: val,
		})
	}

	// Update the deployment with the new environment variable
	_, err = e.kClient.AppsV1().Deployments("default").Update(context.TODO(), deployment, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %v", err)
	}

	// Wait for the deployment to be ready
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutChan := time.After(60 * time.Second)

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timed out waiting for deployment to be ready")
		case <-ticker.C:
			// Get the latest status of the deployment
			deployment, err := e.kClient.AppsV1().Deployments("default").Get(context.TODO(), deploymentName, v1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get deployment: %v", err)
			}

			// Check if the deployment has finished updating
			if deployment.Status.AvailableReplicas == *deployment.Spec.Replicas {
				// All replicas are available, the deployment is complete
				return nil
			}
		}
	}
}

func (e *Environment) setCertificates() error {
	// Helper function to retrieve and validate certificate data from a secret
	getCertificateData := func(secretName string) (cert, key, ca []byte, err error) {
		secret, err := e.kClient.CoreV1().Secrets("default").Get(context.TODO(), secretName, v1.GetOptions{})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not get secret %q: %w", secretName, err)
		}

		cert, certExist := secret.Data["tls.crt"]
		key, keyExist := secret.Data["tls.key"]
		ca, caExist := secret.Data["ca.crt"]

		if !certExist || !keyExist || !caExist {
			return nil, nil, nil, fmt.Errorf("secret %q does not contain all required keys", secretName)
		}

		return cert, key, ca, nil
	}

	// Helper function to write a file and handle errors
	getCertInfo := func(certPem []byte, keyPem []byte, caPem []byte) (*shared.CertInfo, error) {
		cert, err := tls.X509KeyPair(certPem, keyPem)
		if err != nil {
			return nil, fmt.Errorf("could not load key pair: %w", err)
		}

		ca, err := shared.ParseCert(caPem)
		if err != nil {
			return nil, fmt.Errorf("could not parse CA certificate: %w", err)
		}

		return shared.NewCertInfo(cert, ca, nil), nil
	}

	// Fetch management and control certificate data
	managementCert, managementKey, managementCA, err := getCertificateData("management-cert-secret")
	if err != nil {
		return err
	}

	certInfo, err := getCertInfo(managementCert, managementKey, managementCA)
	if err != nil {
		return err
	}
	e.managementCert = certInfo

	controlCert, controlKey, controlCA, err := getCertificateData("control-cert-secret")
	if err != nil {
		return err
	}
	certInfo, err = getCertInfo(controlCert, controlKey, controlCA)
	if err != nil {
		return err
	}
	e.controlCert = certInfo

	return nil
}

func (e *Environment) setKubeClient() error {
	var kubeconfig string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return fmt.Errorf("could not find home directory for kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("could not build kubeconfig: %v", err)
	}

	kClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("could not create kubernetes client: %v", err)
	}

	e.kClient = kClient
	return nil
}

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}

		// Move to the parent directory, Dir excludes the last file path element
		parent := filepath.Dir(wd)
		if parent == wd {
			// Reached the root of the filesystem
			break
		}

		wd = parent
	}

	panic("could not find go.mod")
}
