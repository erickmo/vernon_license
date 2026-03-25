package jwt_test

import (
	"testing"
	"time"

	jwtpkg "github.com/flashlab/vernon-license/pkg/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
)

const testSecret = "super-secret-key-for-testing-only"

// buildExpiredToken membuat JWT string yang sudah expired menggunakan golang-jwt langsung.
// Diperlukan karena Sign() selalu menghasilkan token dengan expiry +24h.
func buildExpiredToken(t *testing.T) string {
	t.Helper()

	now := time.Now().UTC()
	claims := &jwtpkg.Claims{
		Sub:  "user-expired",
		Name: "Expired User",
		Role: "sales",
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject:   "user-expired",
			IssuedAt:  gojwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: gojwt.NewNumericDate(now.Add(-1 * time.Hour)), // expired 1 jam lalu
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("buildExpiredToken: failed to sign: %v", err)
	}
	return signed
}

func TestSign_ValidToken(t *testing.T) {
	t.Log("=== TEST: Sign ValidToken ===")
	t.Log("Goal    : Sign() harus menghasilkan non-empty JWT string tanpa error")
	t.Log("Flow    : Panggil Sign() dengan data valid → cek hasil bukan empty string")

	token, err := jwtpkg.Sign("user-123", "Alice", "admin", testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Log("Status  : FAIL")
		t.Error("expected non-empty token string")
		return
	}
	t.Logf("Result  : token length = %d", len(token))
	t.Log("Status  : PASS")
}

func TestVerify_ValidToken(t *testing.T) {
	t.Log("=== TEST: Verify ValidToken ===")
	t.Log("Goal    : Verify() harus berhasil parse token yang baru di-sign")
	t.Log("Flow    : Sign() → Verify() dengan secret sama → tidak ada error, claims tidak nil")

	token, err := jwtpkg.Sign("user-abc", "Bob", "sales", testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("Sign() unexpected error: %v", err)
	}

	claims, err := jwtpkg.Verify(token, testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("Verify() unexpected error: %v", err)
	}
	if claims == nil {
		t.Log("Status  : FAIL")
		t.Error("expected non-nil claims")
		return
	}
	t.Logf("Result  : claims.Sub=%s, claims.Name=%s, claims.Role=%s", claims.Sub, claims.Name, claims.Role)
	t.Log("Status  : PASS")
}

func TestVerify_WrongSecret(t *testing.T) {
	t.Log("=== TEST: Verify WrongSecret ===")
	t.Log("Goal    : Verify() harus return error jika secret berbeda dari saat signing")
	t.Log("Flow    : Sign() dengan secret A → Verify() dengan secret B → expect error")

	token, err := jwtpkg.Sign("user-xyz", "Carol", "project_owner", testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("Sign() unexpected error: %v", err)
	}

	_, err = jwtpkg.Verify(token, "wrong-secret-key")
	if err == nil {
		t.Log("Status  : FAIL")
		t.Error("expected error for wrong secret, got nil")
		return
	}
	t.Logf("Result  : error = %v", err)
	t.Log("Status  : PASS")
}

func TestVerify_ExpiredToken(t *testing.T) {
	t.Log("=== TEST: Verify ExpiredToken ===")
	t.Log("Goal    : Verify() harus return error untuk token yang sudah expired")
	t.Log("Flow    : Buat token dengan ExpiresAt = now-1h → Verify() → expect error")

	expiredToken := buildExpiredToken(t)

	_, err := jwtpkg.Verify(expiredToken, testSecret)
	if err == nil {
		t.Log("Status  : FAIL")
		t.Error("expected error for expired token, got nil")
		return
	}
	t.Logf("Result  : error = %v", err)
	t.Log("Status  : PASS")
}

func TestSign_ClaimsPreserved(t *testing.T) {
	t.Log("=== TEST: Sign ClaimsPreserved ===")
	t.Log("Goal    : Claims.Sub, Name, Role harus tersimpan persis sesuai input")
	t.Log("Flow    : Sign() dengan nilai spesifik → Verify() → bandingkan setiap field")

	userID := "user-999"
	name := "Dave"
	role := "superuser"

	token, err := jwtpkg.Sign(userID, name, role, testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("Sign() unexpected error: %v", err)
	}

	claims, err := jwtpkg.Verify(token, testSecret)
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("Verify() unexpected error: %v", err)
	}

	failed := false
	if claims.Sub != userID {
		t.Errorf("Sub: expected %q, got %q", userID, claims.Sub)
		failed = true
	}
	if claims.Name != name {
		t.Errorf("Name: expected %q, got %q", name, claims.Name)
		failed = true
	}
	if claims.Role != role {
		t.Errorf("Role: expected %q, got %q", role, claims.Role)
		failed = true
	}
	if failed {
		t.Log("Status  : FAIL")
		return
	}
	t.Logf("Result  : Sub=%s, Name=%s, Role=%s", claims.Sub, claims.Name, claims.Role)
	t.Log("Status  : PASS")
}
