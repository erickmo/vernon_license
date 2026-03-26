package publicapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/flashlab/vernon-license/internal/publicapi"
)

// mockLicenseRepository adalah mock untuk testing
type mockLicenseRepository struct {
	createFn                  func(ctx context.Context, l *domain.ClientLicense) error
	findByOTPFn               func(ctx context.Context, code, slug string) (*domain.ClientLicense, error)
	findByCompanyAndProductFn func(ctx context.Context, companyID, productID uuid.UUID) (*domain.ClientLicense, error)
	updateRegistrationFn      func(ctx context.Context, id uuid.UUID, instanceURL, instanceName string) error
	updateLastPullAtFn        func(ctx context.Context, id uuid.UUID) error
}

func (m *mockLicenseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ClientLicense, error) {
	return nil, nil
}

func (m *mockLicenseRepository) FindByKey(ctx context.Context, key string) (*domain.ClientLicense, error) {
	return nil, nil
}

func (m *mockLicenseRepository) FindByOTP(ctx context.Context, code, slug string) (*domain.ClientLicense, error) {
	if m.findByOTPFn != nil {
		return m.findByOTPFn(ctx, code, slug)
	}
	return nil, domain.ErrLicenseNotFound
}

func (m *mockLicenseRepository) FindByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.ClientLicense, error) {
	return nil, nil
}

func (m *mockLicenseRepository) FindAll(ctx context.Context) ([]*domain.ClientLicense, error) {
	return nil, nil
}

func (m *mockLicenseRepository) FindByCompanyAndProduct(ctx context.Context, companyID, productID uuid.UUID) (*domain.ClientLicense, error) {
	if m.findByCompanyAndProductFn != nil {
		return m.findByCompanyAndProductFn(ctx, companyID, productID)
	}
	return nil, domain.ErrLicenseNotFound
}

func (m *mockLicenseRepository) FindExpiring(ctx context.Context, withinDays int) ([]*domain.ClientLicense, error) {
	return nil, nil
}

func (m *mockLicenseRepository) Create(ctx context.Context, l *domain.ClientLicense) error {
	if m.createFn != nil {
		return m.createFn(ctx, l)
	}
	return nil
}

func (m *mockLicenseRepository) Update(ctx context.Context, l *domain.ClientLicense) error {
	return nil
}

func (m *mockLicenseRepository) UpdateRegistration(ctx context.Context, id uuid.UUID, instanceURL, instanceName string) error {
	if m.updateRegistrationFn != nil {
		return m.updateRegistrationFn(ctx, id, instanceURL, instanceName)
	}
	return nil
}

func (m *mockLicenseRepository) UpdateLastPullAt(ctx context.Context, id uuid.UUID) error {
	if m.updateLastPullAtFn != nil {
		return m.updateLastPullAtFn(ctx, id)
	}
	return nil
}

func (m *mockLicenseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return nil
}

func (m *mockLicenseRepository) UpdateSuperuser(ctx context.Context, id uuid.UUID, username string) error {
	return nil
}

func (m *mockLicenseRepository) FindByCompany(ctx context.Context, companyID uuid.UUID) ([]*domain.ClientLicense, error) {
	return nil, nil
}

// MockProductRepository
type mockProductRepository struct {
	findByIDFn   func(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	findBySlugFn func(ctx context.Context, slug string) (*domain.Product, error)
}

func (m *mockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrProductNotFound
}

func (m *mockProductRepository) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(ctx, slug)
	}
	return nil, domain.ErrProductNotFound
}

func (m *mockProductRepository) FindAll(ctx context.Context, includeInactive bool) ([]*domain.Product, error) {
	return nil, nil
}

func (m *mockProductRepository) Create(ctx context.Context, p *domain.Product) error {
	return nil
}

func (m *mockProductRepository) Update(ctx context.Context, p *domain.Product) error {
	return nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// MockAuditLogRepository
type mockAuditLogRepository struct {
	createFn func(ctx context.Context, log *domain.AuditLog) error
}

func (m *mockAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	if m.createFn != nil {
		return m.createFn(ctx, log)
	}
	return nil
}

func (m *mockAuditLogRepository) FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.AuditLog, error) {
	return nil, nil
}

// MockCompanyRepository
type mockCompanyRepository struct {
	findByNameFn func(ctx context.Context, name string) (*domain.Company, error)
	createFn     func(ctx context.Context, c *domain.Company) error
}

func (m *mockCompanyRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return nil, domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) FindAll(ctx context.Context) ([]*domain.Company, error) {
	return nil, nil
}

func (m *mockCompanyRepository) FindByName(ctx context.Context, name string) (*domain.Company, error) {
	if m.findByNameFn != nil {
		return m.findByNameFn(ctx, name)
	}
	return nil, domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) Create(ctx context.Context, c *domain.Company) error {
	if m.createFn != nil {
		return m.createFn(ctx, c)
	}
	return nil
}

func (m *mockCompanyRepository) Update(ctx context.Context, c *domain.Company) error {
	return nil
}

func (m *mockCompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// MockOTPRepository
type mockOTPRepository struct {
	isActiveFn func(ctx context.Context, code string) error
}

func (m *mockOTPRepository) IsActive(ctx context.Context, code string) error {
	if m.isActiveFn != nil {
		return m.isActiveFn(ctx, code)
	}
	return nil
}

func (m *mockOTPRepository) GetActive(ctx context.Context) (string, error) {
	return "mock-otp", nil
}

func TestRegisterHandler_ValidRequest(t *testing.T) {
	t.Log("=== TEST: RegisterHandler ValidRequest ===")
	t.Log("Goal    : Registrasi dengan valid OTP & product harus sukses, license dibuat dengan status pending")
	t.Log("Flow    : POST /api/v1/register dengan valid data → expect 201 dengan license_key & status=pending")

	mockLic := &mockLicenseRepository{
		createFn: func(ctx context.Context, l *domain.ClientLicense) error {
			return nil
		},
	}

	mockProd := &mockProductRepository{
		findBySlugFn: func(ctx context.Context, slug string) (*domain.Product, error) {
			if slug == "flasherp" {
				return &domain.Product{
					ID:       uuid.New(),
					Name:     "FlashERP",
					Slug:     "flasherp",
					IsActive: true,
				}, nil
			}
			return nil, domain.ErrProductNotFound
		},
	}

	mockAudit := &mockAuditLogRepository{}

	mockOTP := &mockOTPRepository{
		isActiveFn: func(ctx context.Context, code string) error {
			if code == "test_code_123" {
				return nil
			}
			return errors.New("not latest active OTP")
		},
	}

	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		LicenseCheckInterval: "6h",
	}

	mockComp := &mockCompanyRepository{}
	handler := publicapi.NewRegisterHandler(mockLic, mockProd, mockComp, mockAudit, cfg, logger, mockOTP)

	reqBody := map[string]string{
		"otp":           "test_code_123",
		"client_name":  "Test Company",
		"instance_url":  "http://test.example.com",
		"app_name":      "flasherp",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Log("Status  : FAIL")
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["license_key"] == "" || resp["license_key"] == nil {
		t.Log("Status  : FAIL")
		t.Errorf("expected non-empty license_key, got %v", resp["license_key"])
		return
	}

	if resp["status"] != "pending" {
		t.Log("Status  : FAIL")
		t.Errorf("expected status=pending, got %v", resp["status"])
		return
	}

	t.Logf("Result  : Status=%d, LicenseKey=%s, LicenseStatus=%s", w.Code, resp["license_key"], resp["status"])
	t.Log("Status  : PASS")
}

func TestRegisterHandler_InvalidOTP(t *testing.T) {
	t.Log("=== TEST: RegisterHandler InvalidOTP ===")
	t.Log("Goal    : Registration dengan OTP yang bukan terbaru harus return 403")
	t.Log("Flow    : POST dengan OTP tidak valid → expect 403 INVALID_CLIENT_CODE")

	mockLic := &mockLicenseRepository{}
	mockProd := &mockProductRepository{
		findBySlugFn: func(ctx context.Context, slug string) (*domain.Product, error) {
			if slug == "flasherp" {
				return &domain.Product{
					ID:       uuid.New(),
					Name:     "FlashERP",
					Slug:     "flasherp",
					IsActive: true,
				}, nil
			}
			return nil, domain.ErrProductNotFound
		},
	}
	mockAudit := &mockAuditLogRepository{}

	mockOTP := &mockOTPRepository{
		isActiveFn: func(ctx context.Context, code string) error {
			return errors.New("not latest active OTP")
		},
	}

	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}

	mockComp := &mockCompanyRepository{}
	handler := publicapi.NewRegisterHandler(mockLic, mockProd, mockComp, mockAudit, cfg, logger, mockOTP)

	reqBody := map[string]string{
		"otp":           "invalid_code",
		"client_name":  "Test",
		"instance_url":  "http://test.com",
		"app_name":      "flasherp",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	if w.Code != http.StatusForbidden {
		t.Log("Status  : FAIL")
		t.Errorf("expected 403, got %d", w.Code)
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	errCode := resp["error"].(map[string]interface{})["code"]

	if errCode != "INVALID_CLIENT_CODE" {
		t.Log("Status  : FAIL")
		t.Errorf("expected error code INVALID_CLIENT_CODE, got %v", errCode)
		return
	}

	t.Logf("Result  : Status=%d, Error=%s", w.Code, errCode)
	t.Log("Status  : PASS")
}

func TestRegisterHandler_ProductNotFound(t *testing.T) {
	t.Log("=== TEST: RegisterHandler ProductNotFound ===")
	t.Log("Goal    : Registration dengan product yang tidak ada harus return 403")
	t.Log("Flow    : POST dengan product_slug yang tidak ada → expect 403 PRODUCT_NOT_FOUND")

	mockLic := &mockLicenseRepository{}
	mockProd := &mockProductRepository{
		findBySlugFn: func(ctx context.Context, slug string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
	}
	mockAudit := &mockAuditLogRepository{}
	mockOTP := &mockOTPRepository{}

	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{}

	mockComp := &mockCompanyRepository{}
	handler := publicapi.NewRegisterHandler(mockLic, mockProd, mockComp, mockAudit, cfg, logger, mockOTP)

	reqBody := map[string]string{
		"otp":           "test_code_123",
		"client_name":  "Test",
		"instance_url":  "http://test.com",
		"app_name":      "nonexistent",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	if w.Code != http.StatusForbidden {
		t.Log("Status  : FAIL")
		t.Errorf("expected 403, got %d", w.Code)
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	errCode := resp["error"].(map[string]interface{})["code"]

	if errCode != "PRODUCT_NOT_FOUND" {
		t.Log("Status  : FAIL")
		t.Errorf("expected error code PRODUCT_NOT_FOUND, got %v", errCode)
		return
	}

	t.Logf("Result  : Status=%d, Error=%s", w.Code, errCode)
	t.Log("Status  : PASS")
}
