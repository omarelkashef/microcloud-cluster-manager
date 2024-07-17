package state

import (
	"sync"

	"github.com/canonical/microcluster/microcluster"

	"github.com/canonical/lxd-site-manager/internal/oidc"
)

// SiteManagerState holds the global state of the site manager.
type SiteManagerState struct {
	// OIDCVerifier is the OpenID Connect verifier used for user authentication and validate authentication for protected API endpoints.
	OIDCVerifier *oidc.Verifier
	// MicroCluster is the MicroCluster instance.
	MicroCluster *microcluster.MicroCluster
	mu           sync.RWMutex
}

// New creates a new SiteManagerState.
func New(m *microcluster.MicroCluster) *SiteManagerState {
	return &SiteManagerState{
		MicroCluster: m,
	}
}

// SetOIDCVerifier sets the OIDCVerifier.
func (s *SiteManagerState) SetOIDCVerifier(verifier *oidc.Verifier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.OIDCVerifier = verifier
}
