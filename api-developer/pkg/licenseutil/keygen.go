package licenseutil

import (
	"crypto/rand"
	"fmt"
)

// GenerateLicenseKey menghasilkan license key aman format FL-XXXXXXXX-XXXXXXXX-XXXXXXXX-XXXXXXXX.
func GenerateLicenseKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("FL-%02X%02X%02X%02X-%02X%02X%02X%02X-%02X%02X%02X%02X-%02X%02X%02X%02X",
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7],
		b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]), nil
}

// GenerateSecurePassword menghasilkan password acak 16 karakter.
func GenerateSecurePassword() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%&*"
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	result := make([]byte, 16)
	for i, v := range b {
		result[i] = charset[int(v)%len(charset)]
	}
	return string(result), nil
}
