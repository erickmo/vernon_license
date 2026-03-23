package provision_license

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type ProvisionLicense struct {
	LicenseKey string
}

type provisionPayload struct {
	LicenseKey         string     `json:"license_key"`
	ClientName         string     `json:"client_name"`
	Plan               string     `json:"plan"`
	Status             string     `json:"status"`
	MaxUsers           *int       `json:"max_users"`
	MaxTransPerMonth   *int       `json:"max_trans_per_month"`
	MaxTransPerDay     *int       `json:"max_trans_per_day"`
	MaxItems           *int       `json:"max_items"`
	MaxCustomers       *int       `json:"max_customers"`
	MaxBranches        *int       `json:"max_branches"`
	ExpiresAt          *time.Time `json:"expires_at"`
	AdminEmail         string     `json:"admin_email"`
	AdminPasswordHash  string     `json:"admin_password_hash"`
}

type Handler struct {
	licenseRepo clientlicense.WriteRepository
	httpClient  *http.Client
}

func NewHandler(licenseRepo clientlicense.WriteRepository) *Handler {
	return &Handler{
		licenseRepo: licenseRepo,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	c, ok := cmd.(ProvisionLicense)
	if !ok {
		return fmt.Errorf("perintah tidak valid")
	}

	license, err := h.licenseRepo.GetByKey(ctx, c.LicenseKey)
	if err != nil {
		return err
	}

	if license.FlashERPURL == nil || *license.FlashERPURL == "" {
		return clientlicense.ErrMissingFlashERPURL
	}

	payload := provisionPayload{
		LicenseKey:       license.LicenseKey,
		ClientName:       license.ClientName,
		Plan:             license.Plan,
		Status:           license.Status,
		MaxUsers:         license.MaxUsers,
		MaxTransPerMonth: license.MaxTransPerMonth,
		MaxTransPerDay:   license.MaxTransPerDay,
		MaxItems:         license.MaxItems,
		MaxCustomers:     license.MaxCustomers,
		MaxBranches:      license.MaxBranches,
		ExpiresAt:        license.ExpiresAt,
		AdminEmail:       license.ClientEmail,
		AdminPasswordHash: "",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal marshal payload: %w", err)
	}

	url := *license.FlashERPURL + "/internal/provision"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gagal buat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if license.ProvisionAPIKey != nil {
		req.Header.Set("X-Provision-Key", *license.ProvisionAPIKey)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal kirim request ke FlashERP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FlashERP menolak provisioning (status %d): %s", resp.StatusCode, string(respBody))
	}

	license.IsProvisioned = true
	license.UpdatedAt = time.Now().UTC()
	return h.licenseRepo.Update(ctx, license)
}
