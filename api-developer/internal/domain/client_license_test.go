package domain_test

import (
	"testing"
	"time"

	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/google/uuid"
)

// newLicense adalah helper untuk membuat ClientLicense minimal dengan status tertentu.
func newLicense(status string) *domain.ClientLicense {
	return &domain.ClientLicense{
		ID:         uuid.New(),
		LicenseKey: "FL-TESTKEY1",
		Status:     status,
		CreatedBy:  uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func TestClientLicense_IsValid_Active(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid Active ===")
	t.Log("Goal    : status=active + tidak ada ExpiresAt → IsValid() = true")
	t.Log("Flow    : Buat license status active tanpa ExpiresAt → panggil IsValid()")

	l := newLicense("active")
	l.ExpiresAt = nil

	result := l.IsValid()
	if !result {
		t.Log("Status  : FAIL")
		t.Error("expected IsValid()=true for active license with no expiry")
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_IsValid_ActiveNotExpired(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid ActiveNotExpired ===")
	t.Log("Goal    : status=active + ExpiresAt > now → IsValid() = true")
	t.Log("Flow    : Buat license status active, ExpiresAt = now+30d → panggil IsValid()")

	l := newLicense("active")
	future := time.Now().Add(30 * 24 * time.Hour)
	l.ExpiresAt = &future

	result := l.IsValid()
	if !result {
		t.Log("Status  : FAIL")
		t.Errorf("expected IsValid()=true for active license expiring at %v", future)
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_IsValid_ActiveExpired(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid ActiveExpired ===")
	t.Log("Goal    : status=active + ExpiresAt < now → IsValid() = false")
	t.Log("Flow    : Buat license status active, ExpiresAt = now-1d → panggil IsValid()")

	l := newLicense("active")
	past := time.Now().Add(-24 * time.Hour)
	l.ExpiresAt = &past

	result := l.IsValid()
	if result {
		t.Log("Status  : FAIL")
		t.Errorf("expected IsValid()=false for active license expired at %v", past)
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_IsValid_Trial(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid Trial ===")
	t.Log("Goal    : status=trial → IsValid() = true (approved trial)")
	t.Log("Flow    : Buat license status trial → panggil IsValid()")

	l := newLicense("trial")

	result := l.IsValid()
	if !result {
		t.Log("Status  : FAIL")
		t.Error("expected IsValid()=true for trial license")
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_IsValid_Suspended(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid Suspended ===")
	t.Log("Goal    : status=suspended → IsValid() = false")
	t.Log("Flow    : Buat license status suspended → panggil IsValid()")

	l := newLicense("suspended")

	result := l.IsValid()
	if result {
		t.Log("Status  : FAIL")
		t.Error("expected IsValid()=false for suspended license")
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_IsValid_Pending(t *testing.T) {
	t.Log("=== TEST: ClientLicense IsValid Pending ===")
	t.Log("Goal    : status=pending → IsValid() = false")
	t.Log("Flow    : Buat license status pending → panggil IsValid()")

	l := newLicense("pending")

	result := l.IsValid()
	if result {
		t.Log("Status  : FAIL")
		t.Error("expected IsValid()=false for pending license")
		return
	}
	t.Logf("Result  : IsValid() = %v", result)
	t.Log("Status  : PASS")
}

func TestClientLicense_ValidReason_Suspended(t *testing.T) {
	t.Log("=== TEST: ClientLicense ValidReason Suspended ===")
	t.Log("Goal    : status=suspended → ValidReason() = \"suspended\"")
	t.Log("Flow    : Buat license status suspended → panggil ValidReason()")

	l := newLicense("suspended")

	reason := l.ValidReason()
	if reason != "suspended" {
		t.Log("Status  : FAIL")
		t.Errorf("expected reason=%q, got %q", "suspended", reason)
		return
	}
	t.Logf("Result  : ValidReason() = %q", reason)
	t.Log("Status  : PASS")
}

func TestClientLicense_ValidReason_Expired(t *testing.T) {
	t.Log("=== TEST: ClientLicense ValidReason Expired ===")
	t.Log("Goal    : status=expired → ValidReason() = \"expired\"")
	t.Log("Flow    : Buat license status expired → panggil ValidReason()")

	l := newLicense("expired")

	reason := l.ValidReason()
	if reason != "expired" {
		t.Log("Status  : FAIL")
		t.Errorf("expected reason=%q, got %q", "expired", reason)
		return
	}
	t.Logf("Result  : ValidReason() = %q", reason)
	t.Log("Status  : PASS")
}

func TestClientLicense_ValidReason_ExpiredByDate(t *testing.T) {
	t.Log("=== TEST: ClientLicense ValidReason ExpiredByDate ===")
	t.Log("Goal    : status=active tapi ExpiresAt < now → ValidReason() = \"expired\"")
	t.Log("Flow    : Buat license status active, ExpiresAt = now-1d → panggil ValidReason()")

	l := newLicense("active")
	past := time.Now().Add(-24 * time.Hour)
	l.ExpiresAt = &past

	reason := l.ValidReason()
	if reason != "expired" {
		t.Log("Status  : FAIL")
		t.Errorf("expected reason=%q for active-but-date-expired, got %q", "expired", reason)
		return
	}
	t.Logf("Result  : ValidReason() = %q", reason)
	t.Log("Status  : PASS")
}

func TestClientLicense_ValidReason_Pending(t *testing.T) {
	t.Log("=== TEST: ClientLicense ValidReason Pending ===")
	t.Log("Goal    : status=pending → ValidReason() = \"pending_approval\"")
	t.Log("Flow    : Buat license status pending → panggil ValidReason()")

	l := newLicense("pending")

	reason := l.ValidReason()
	if reason != "pending_approval" {
		t.Log("Status  : FAIL")
		t.Errorf("expected reason=%q, got %q", "pending_approval", reason)
		return
	}
	t.Logf("Result  : ValidReason() = %q", reason)
	t.Log("Status  : PASS")
}
