package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canonical/lxd-cluster-manager/internal/pkg/api/models/v1"
)

// VerifyHMAC verifies the HMAC signature of a request.
func VerifyHMAC(payload models.RemoteClusterPost, r *http.Request, secret string) (bool, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %v", err)
	}

	sig := r.Header.Get("X-CLUSTER-SIGNATURE")
	if sig == "" {
		return false, fmt.Errorf("missing signature header")
	}

	decodedSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// recompute the HMAC
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(reqBody)
	expectedMac := mac.Sum(nil)

	return hmac.Equal(decodedSig, expectedMac), nil
}
