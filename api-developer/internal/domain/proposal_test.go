package domain_test

import (
	"testing"
	"time"

	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/google/uuid"
)

// baseProposal membuat Proposal minimal untuk digunakan sebagai basis test.
func baseProposal(version int) *domain.Proposal {
	maxUsers := 10
	amount := 5000.0
	notes := "initial notes"
	return &domain.Proposal{
		ID:             uuid.New(),
		ProjectID:      uuid.New(),
		CompanyID:      uuid.New(),
		ProductID:      uuid.New(),
		Version:        version,
		Status:         "draft",
		Modules:        []string{"accounting", "inventory"},
		Apps:           []string{"web"},
		Plan:           "standard",
		MaxUsers:       &maxUsers,
		ContractAmount: &amount,
		Notes:          &notes,
		SubmittedBy:    uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func TestComputeChangelog_NoChanges(t *testing.T) {
	t.Log("=== TEST: ComputeChangelog NoChanges ===")
	t.Log("Goal    : Dua proposal identik → Changes kosong, Unchanged tidak kosong")
	t.Log("Flow    : Buat prev dan curr identik → ComputeChangelog() → cek len(Changes)==0 dan len(Unchanged)>0")

	prev := baseProposal(1)
	curr := baseProposal(2)

	// Salin semua field comparable agar identik
	curr.Plan = prev.Plan
	curr.Modules = prev.Modules
	curr.Apps = prev.Apps
	curr.ContractAmount = prev.ContractAmount
	curr.ExpiresAt = prev.ExpiresAt
	curr.MaxUsers = prev.MaxUsers
	curr.MaxTransPerMonth = prev.MaxTransPerMonth
	curr.MaxTransPerDay = prev.MaxTransPerDay
	curr.MaxItems = prev.MaxItems
	curr.MaxCustomers = prev.MaxCustomers
	curr.MaxBranches = prev.MaxBranches
	curr.MaxStorage = prev.MaxStorage
	curr.Notes = prev.Notes

	cl := domain.ComputeChangelog(prev, curr)

	if len(cl.Changes) != 0 {
		t.Log("Status  : FAIL")
		t.Errorf("expected 0 changes, got %d: %+v", len(cl.Changes), cl.Changes)
		return
	}
	if len(cl.Unchanged) == 0 {
		t.Log("Status  : FAIL")
		t.Error("expected non-empty Unchanged list")
		return
	}
	t.Logf("Result  : Changes=%d, Unchanged=%d", len(cl.Changes), len(cl.Unchanged))
	t.Log("Status  : PASS")
}

func TestComputeChangelog_PlanChanged(t *testing.T) {
	t.Log("=== TEST: ComputeChangelog PlanChanged ===")
	t.Log("Goal    : Jika Plan berubah → ada di Changes dengan field=\"plan\"")
	t.Log("Flow    : Buat prev.Plan=\"standard\", curr.Plan=\"enterprise\" → ComputeChangelog() → cek Changes")

	prev := baseProposal(1)
	curr := baseProposal(2)
	curr.Plan = "enterprise" // ubah plan

	cl := domain.ComputeChangelog(prev, curr)

	found := false
	for _, entry := range cl.Changes {
		if entry.Field == "plan" {
			found = true
			if entry.OldValue != "standard" {
				t.Log("Status  : FAIL")
				t.Errorf("plan OldValue: expected %q, got %v", "standard", entry.OldValue)
				return
			}
			if entry.NewValue != "enterprise" {
				t.Log("Status  : FAIL")
				t.Errorf("plan NewValue: expected %q, got %v", "enterprise", entry.NewValue)
				return
			}
		}
	}
	if !found {
		t.Log("Status  : FAIL")
		t.Errorf("field \"plan\" tidak ditemukan di Changes: %+v", cl.Changes)
		return
	}
	t.Logf("Result  : found plan change in Changes (len=%d)", len(cl.Changes))
	t.Log("Status  : PASS")
}

func TestComputeChangelog_MultipleChanges(t *testing.T) {
	t.Log("=== TEST: ComputeChangelog MultipleChanges ===")
	t.Log("Goal    : Beberapa field berubah → semua harus ada di Changes")
	t.Log("Flow    : Ubah plan, max_users, dan contract_amount → cek semua ada di Changes")

	prev := baseProposal(1)
	curr := baseProposal(2)

	// Ubah 3 field
	curr.Plan = "premium"
	newMaxUsers := 50
	curr.MaxUsers = &newMaxUsers
	newAmount := 12000.0
	curr.ContractAmount = &newAmount

	cl := domain.ComputeChangelog(prev, curr)

	wantChangedFields := map[string]bool{
		"plan":            false,
		"max_users":       false,
		"contract_amount": false,
	}

	for _, entry := range cl.Changes {
		if _, ok := wantChangedFields[entry.Field]; ok {
			wantChangedFields[entry.Field] = true
		}
	}

	allFound := true
	for field, found := range wantChangedFields {
		if !found {
			t.Errorf("expected field %q in Changes but not found", field)
			allFound = false
		}
	}
	if !allFound {
		t.Log("Status  : FAIL")
		return
	}
	t.Logf("Result  : %d changes found, all 3 expected fields present", len(cl.Changes))
	t.Log("Status  : PASS")
}

func TestComputeChangelog_VersionTracked(t *testing.T) {
	t.Log("=== TEST: ComputeChangelog VersionTracked ===")
	t.Log("Goal    : Changelog.ComparedToVersion harus sama dengan prev.Version")
	t.Log("Flow    : Buat prev dengan Version=3, curr dengan Version=4 → cek ComparedToVersion=3")

	prev := baseProposal(3)
	curr := baseProposal(4)
	curr.Plan = "changed" // minimal satu perubahan

	cl := domain.ComputeChangelog(prev, curr)

	if cl.ComparedToVersion != prev.Version {
		t.Log("Status  : FAIL")
		t.Errorf("ComparedToVersion: expected %d (prev.Version), got %d", prev.Version, cl.ComparedToVersion)
		return
	}
	t.Logf("Result  : ComparedToVersion = %d", cl.ComparedToVersion)
	t.Log("Status  : PASS")
}
