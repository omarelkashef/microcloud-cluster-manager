package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/database"
	"github.com/canonical/lxd/shared"
)

type Config struct {
	// system configs
	Version    string
	ApiVersion string
	ServerCert *shared.CertInfo
	// db configs
	database.DBConfig
	// api configs
	ServerHost     string
	ManagementPort string
	ControlPort    string
	AllowedOrigins []string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	// oidc configs
	OIDCClientID string
	OIDCIssuer   string
	OIDCAudience string
}

// getEnvOrDefault retrieves an environment variable or returns the default value if not set.
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns the default value.
func getEnvAsInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value != "" {
		return strconv.Atoi(value)
	}
	return defaultValue, nil
}

// getServiceCert loads the TLS certificate and key from the environment.
func getServiceCert() (*shared.CertInfo, error) {
	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	tlsKeyPath := os.Getenv("TLS_KEY_PATH")
	certCAPath := os.Getenv("TLS_CA_PATH")

	if tlsCertPath == "" || tlsKeyPath == "" || certCAPath == "" {
		return nil, fmt.Errorf("TLS_CERT_PATH, TLS_KEY_PATH and TLS_CERT_PATH are required")
	}

	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	ca, err := shared.ReadCert(certCAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate: %w", err)
	}

	return shared.NewCertInfo(cert, ca, nil), nil
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig(requireCerts bool) (*Config, error) {
	// DB Config
	dbMaxIdleConns, err := getEnvAsInt("DB_MAX_IDLE", 10)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE: %w", err)
	}
	dbMaxOpenConns, err := getEnvAsInt("DB_MAX_OPEN", 2)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN: %w", err)
	}

	// Server Cert
	var serverCert *shared.CertInfo
	if requireCerts {
		serverCert, err = getServiceCert()
		if err != nil {
			return nil, err
		}
	}

	// OIDC Config
	oidcClientID := os.Getenv("OIDC_CLIENT_ID")
	oidcIssuer := os.Getenv("OIDC_ISSUER")
	oidcAudience := os.Getenv("OIDC_AUDIENCE")
	if oidcClientID == "" || oidcIssuer == "" || oidcAudience == "" {
		return nil, fmt.Errorf("OIDC_CLIENT_ID, OIDC_ISSUER, and OIDC_AUDIENCE are required")
	}

	return &Config{
		Version:        getEnvOrDefault("VERSION", "development"),
		ApiVersion:     getEnvOrDefault("API_VERSION", "1.0"),
		ServerHost:     getEnvOrDefault("SERVER_HOST", "localhost"),
		ManagementPort: getEnvOrDefault("MANAGEMENT_PORT", "9000"),
		ControlPort:    getEnvOrDefault("CONTROL_PORT", "9001"),
		AllowedOrigins: []string{"*"},
		ReadTimeout:    10,
		WriteTimeout:   10,
		IdleTimeout:    60,
		ServerCert:     serverCert,
		DBConfig: database.DBConfig{
			DBPort:         getEnvOrDefault("DB_PORT", "5432"),
			DBUser:         getEnvOrDefault("DB_USER", "admin"),
			DBPassword:     getEnvOrDefault("DB_PASSWORD", "admin"),
			DBHost:         getEnvOrDefault("DB_HOST", "localhost"),
			DBName:         getEnvOrDefault("DB_NAME", "cm"),
			DBMaxIdleConns: dbMaxIdleConns,
			DBMaxOpenConns: dbMaxOpenConns,
			DBDisableTLS:   getEnvOrDefault("DB_DISABLE_TLS", "true") == "true",
		},
		OIDCClientID: oidcClientID,
		OIDCIssuer:   oidcIssuer,
		OIDCAudience: oidcAudience,
	}, nil
}
