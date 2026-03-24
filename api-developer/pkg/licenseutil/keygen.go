// Package licenseutil menyediakan utilitas untuk generate license key dan provision API key.
package licenseutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	// licenseKeyChars adalah karakter yang digunakan untuk license key (uppercase alphanumeric).
	licenseKeyChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// licenseKeyLength adalah panjang random suffix dari license key (8 karakter).
	licenseKeyLength = 8
)

// GenerateLicenseKey menghasilkan license key dengan format FL-XXXXXXXX,
// di mana X adalah 8 karakter random uppercase alphanumeric.
func GenerateLicenseKey() (string, error) {
	suffix, err := randomString(licenseKeyChars, licenseKeyLength)
	if err != nil {
		return "", fmt.Errorf("GenerateLicenseKey: %w", err)
	}
	return "FL-" + suffix, nil
}

// GenerateProvisionAPIKey menghasilkan 32-character random hex string
// yang digunakan sebagai provision API key saat setup license.
func GenerateProvisionAPIKey() (string, error) {
	b := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("GenerateProvisionAPIKey: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// randomString menghasilkan random string sepanjang n karakter dari charset yang diberikan.
func randomString(charset string, n int) (string, error) {
	charsetLen := big.NewInt(int64(len(charset)))
	result := make([]byte, n)
	for i := range result {
		idx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("randomString: %w", err)
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result), nil
}
