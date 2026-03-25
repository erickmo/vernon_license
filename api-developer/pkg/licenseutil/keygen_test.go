package licenseutil_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/flashlab/vernon-license/pkg/licenseutil"
)

func TestGenerateLicenseKey_Format(t *testing.T) {
	t.Log("=== TEST: GenerateLicenseKey Format ===")
	t.Log("Goal    : License key harus berformat FL-XXXXXXXX (prefix FL- + 8 char uppercase alphanumeric)")
	t.Log("Flow    : Generate key → cek prefix FL- → cek total length 11 → cek karakter valid")

	key, err := licenseutil.GenerateLicenseKey()
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(key, "FL-") {
		t.Log("Status  : FAIL")
		t.Errorf("expected prefix FL-, got %q", key)
		return
	}
	if len(key) != 11 { // "FL-" (3) + 8 chars = 11
		t.Log("Status  : FAIL")
		t.Errorf("expected length 11, got %d (key=%q)", len(key), key)
		return
	}

	// suffix harus uppercase alphanumeric saja
	suffix := key[3:]
	re := regexp.MustCompile(`^[A-Z0-9]{8}$`)
	if !re.MatchString(suffix) {
		t.Log("Status  : FAIL")
		t.Errorf("suffix %q contains invalid characters (expected [A-Z0-9]{8})", suffix)
		return
	}

	t.Logf("Result  : %s", key)
	t.Log("Status  : PASS")
}

func TestGenerateLicenseKey_Unique(t *testing.T) {
	t.Log("=== TEST: GenerateLicenseKey Unique ===")
	t.Log("Goal    : 1000 license keys yang digenerate tidak boleh ada duplikat")
	t.Log("Flow    : Generate 1000 keys → simpan di map → cek tidak ada collision")

	const iterations = 1000
	seen := make(map[string]struct{}, iterations)

	for i := 0; i < iterations; i++ {
		key, err := licenseutil.GenerateLicenseKey()
		if err != nil {
			t.Log("Status  : FAIL")
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if _, exists := seen[key]; exists {
			t.Log("Status  : FAIL")
			t.Errorf("duplicate license key found: %s (iteration %d)", key, i)
			return
		}
		seen[key] = struct{}{}
	}

	t.Logf("Result  : %d unique keys generated, no duplicates", iterations)
	t.Log("Status  : PASS")
}

func TestGenerateOTP_Length(t *testing.T) {
	t.Log("=== TEST: GenerateOTP Length ===")
	t.Log("Goal    : OTP harus 32 karakter hex lowercase")
	t.Log("Flow    : Generate OTP → cek length 32 → cek format hex lowercase")

	key, err := licenseutil.GenerateOTP()
	if err != nil {
		t.Log("Status  : FAIL")
		t.Fatalf("unexpected error: %v", err)
	}

	if len(key) != 32 {
		t.Log("Status  : FAIL")
		t.Errorf("expected length 32, got %d (key=%q)", len(key), key)
		return
	}

	re := regexp.MustCompile(`^[0-9a-f]{32}$`)
	if !re.MatchString(key) {
		t.Log("Status  : FAIL")
		t.Errorf("key %q is not a valid 32-char hex string", key)
		return
	}

	t.Logf("Result  : %s", key)
	t.Log("Status  : PASS")
}

func TestGenerateOTP_Unique(t *testing.T) {
	t.Log("=== TEST: GenerateOTP Unique ===")
	t.Log("Goal    : 1000 OTP yang digenerate tidak boleh ada duplikat")
	t.Log("Flow    : Generate 1000 OTP → simpan di map → cek tidak ada collision")

	const iterations = 1000
	seen := make(map[string]struct{}, iterations)

	for i := 0; i < iterations; i++ {
		key, err := licenseutil.GenerateOTP()
		if err != nil {
			t.Log("Status  : FAIL")
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if _, exists := seen[key]; exists {
			t.Log("Status  : FAIL")
			t.Errorf("duplicate OTP found: %s (iteration %d)", key, i)
			return
		}
		seen[key] = struct{}{}
	}

	t.Logf("Result  : %d unique keys generated, no duplicates", iterations)
	t.Log("Status  : PASS")
}
