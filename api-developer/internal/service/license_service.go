package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/flashlab/vernon-license/pkg/licenseutil"
)

// LicenseService menyediakan business logic untuk manajemen license.
type LicenseService struct {
	repo        domain.LicenseRepository
	productRepo domain.ProductRepository
	auditRepo   domain.AuditLogRepository
	notifRepo   domain.NotificationRepository
	otpRepo     domain.OTPRepository
	logger      *zap.Logger
}

// NewLicenseService membuat instance LicenseService baru dengan semua dependency yang dibutuhkan.
func NewLicenseService(
	repo domain.LicenseRepository,
	productRepo domain.ProductRepository,
	auditRepo domain.AuditLogRepository,
	notifRepo domain.NotificationRepository,
	otpRepo domain.OTPRepository,
	logger *zap.Logger,
) *LicenseService {
	return &LicenseService{
		repo:        repo,
		productRepo: productRepo,
		auditRepo:   auditRepo,
		notifRepo:   notifRepo,
		otpRepo:     otpRepo,
		logger:      logger,
	}
}

// CreateLicenseRequest adalah input untuk membuat license secara langsung (direct create oleh PO).
type CreateLicenseRequest struct {
	ProjectID        uuid.UUID  `json:"project_id"`
	CompanyID        uuid.UUID  `json:"company_id"`
	ProductID        uuid.UUID  `json:"product_id"`
	Plan             string     `json:"plan"`
	Modules          []string   `json:"modules"`
	Apps             []string   `json:"apps"`
	ContractAmount   *float64   `json:"contract_amount"`
	Description      *string    `json:"description"`
	MaxUsers         *int       `json:"max_users"`
	MaxTransPerMonth *int       `json:"max_trans_per_month"`
	MaxTransPerDay   *int       `json:"max_trans_per_day"`
	MaxItems         *int       `json:"max_items"`
	MaxCustomers     *int       `json:"max_customers"`
	MaxBranches      *int       `json:"max_branches"`
	MaxStorage       *int       `json:"max_storage"`
	ExpiresAt        *time.Time `json:"expires_at"`
	CheckInterval    string     `json:"check_interval"` // "1h"|"6h"|"24h"
}

// UpdateConstraintsRequest adalah input untuk memperbarui constraint license.
type UpdateConstraintsRequest struct {
	Modules          []string   `json:"modules"`
	Apps             []string   `json:"apps"`
	MaxUsers         *int       `json:"max_users"`
	MaxTransPerMonth *int       `json:"max_trans_per_month"`
	MaxTransPerDay   *int       `json:"max_trans_per_day"`
	MaxItems         *int       `json:"max_items"`
	MaxCustomers     *int       `json:"max_customers"`
	MaxBranches      *int       `json:"max_branches"`
	MaxStorage       *int       `json:"max_storage"`
	ExpiresAt        *time.Time `json:"expires_at"`
	CheckInterval    string     `json:"check_interval"`
}

// List mengembalikan semua license yang belum dihapus (accessible oleh semua role).
func (s *LicenseService) List(ctx context.Context) ([]*domain.ClientLicense, error) {
	licenses, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("LicenseService.List: %w", err)
	}
	return licenses, nil
}

// ListByProject mengembalikan semua license untuk project tertentu.
func (s *LicenseService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.ClientLicense, error) {
	licenses, err := s.repo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("LicenseService.ListByProject: %w", err)
	}
	return licenses, nil
}

// GetByID mengambil satu license berdasarkan UUID-nya.
func (s *LicenseService) GetByID(ctx context.Context, id uuid.UUID) (*domain.ClientLicense, error) {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("LicenseService.GetByID: %w", err)
	}
	return license, nil
}

// DirectCreate membuat license baru secara langsung oleh PO tanpa melalui proposal.
// Menghasilkan license_key (FL-XXXXXXXX) dan OTP (32-char hex) secara otomatis.
// Status awal adalah "pending". Produk yang dimaksud harus exist dan aktif.
func (s *LicenseService) DirectCreate(ctx context.Context, req CreateLicenseRequest, actorID uuid.UUID, actorName string) (*domain.ClientLicense, error) {
	// Validasi product exist dan active
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("LicenseService.DirectCreate: %w", domain.ErrProductNotFound)
	}
	if !product.IsActive {
		return nil, fmt.Errorf("LicenseService.DirectCreate: %w", domain.ErrProductInactive)
	}

	// Generate license key
	licenseKey, err := licenseutil.GenerateLicenseKey()
	if err != nil {
		return nil, fmt.Errorf("LicenseService.DirectCreate generate license key: %w", err)
	}

	// Generate registration code
	registrationCode, err := licenseutil.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("LicenseService.DirectCreate generate registration code: %w", err)
	}

	// Normalize check interval default
	checkInterval := req.CheckInterval
	if checkInterval == "" {
		checkInterval = "6h"
	}

	now := time.Now().UTC()
	license := &domain.ClientLicense{
		ID:                     uuid.New(),
		LicenseKey:             licenseKey,
		ProjectID:              &req.ProjectID,
		CompanyID:              &req.CompanyID,
		ProductID:              req.ProductID,
		Plan:                   req.Plan,
		Status:                 "pending",
		Modules:                req.Modules,
		Apps:                   req.Apps,
		ContractAmount:         req.ContractAmount,
		Description:            req.Description,
		MaxUsers:               req.MaxUsers,
		MaxTransPerMonth:       req.MaxTransPerMonth,
		MaxTransPerDay:         req.MaxTransPerDay,
		MaxItems:               req.MaxItems,
		MaxCustomers:           req.MaxCustomers,
		MaxBranches:            req.MaxBranches,
		MaxStorage:             req.MaxStorage,
		ExpiresAt:              req.ExpiresAt,
		OTP: &registrationCode,
		CheckInterval:          checkInterval,
		IsRegistered:           false,
		CreatedBy:              &actorID,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.repo.Create(ctx, license); err != nil {
		return nil, fmt.Errorf("LicenseService.DirectCreate persist: %w", err)
	}

	// Audit log setelah operasi sukses
	meta, _ := json.Marshal(map[string]string{"creation_method": "direct"})
	s.logAuditWithMeta(ctx, "license", license.ID, "license_created", actorID, actorName, json.RawMessage("{}"), meta)

	return license, nil
}

// ListByCompany mengembalikan semua licenses untuk sebuah company.
func (s *LicenseService) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*domain.ClientLicense, error) {
	licenses, err := s.repo.FindByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("LicenseService.ListByCompany: %w", err)
	}
	return licenses, nil
}

// CreateFromProposal membuat license dari proposal yang sudah disetujui oleh PO.
// Constraints disalin dari proposal. Status awal adalah "pending".
// Dipanggil oleh ProposalService saat proposal di-approve.
func (s *LicenseService) CreateFromProposal(ctx context.Context, proposal *domain.Proposal, actorID uuid.UUID, actorName string) (*domain.ClientLicense, error) {
	// Generate license key
	licenseKey, err := licenseutil.GenerateLicenseKey()
	if err != nil {
		return nil, fmt.Errorf("LicenseService.CreateFromProposal generate license key: %w", err)
	}

	// Generate registration code
	registrationCode, err := licenseutil.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("LicenseService.CreateFromProposal generate registration code: %w", err)
	}

	now := time.Now().UTC()
	license := &domain.ClientLicense{
		ID:                     uuid.New(),
		LicenseKey:             licenseKey,
		ProjectID:              &proposal.ProjectID,
		CompanyID:              &proposal.CompanyID,
		ProductID:              proposal.ProductID,
		Plan:                   proposal.Plan,
		Status:                 "pending",
		Modules:                proposal.Modules,
		Apps:                   proposal.Apps,
		ContractAmount:         proposal.ContractAmount,
		MaxUsers:               proposal.MaxUsers,
		MaxTransPerMonth:       proposal.MaxTransPerMonth,
		MaxTransPerDay:         proposal.MaxTransPerDay,
		MaxItems:               proposal.MaxItems,
		MaxCustomers:           proposal.MaxCustomers,
		MaxBranches:            proposal.MaxBranches,
		MaxStorage:             proposal.MaxStorage,
		ExpiresAt:              proposal.ExpiresAt,
		OTP: &registrationCode,
		CheckInterval:          "6h",
		IsRegistered:           false,
		ProposalID:             &proposal.ID,
		CreatedBy:              &actorID,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.repo.Create(ctx, license); err != nil {
		return nil, fmt.Errorf("LicenseService.CreateFromProposal persist: %w", err)
	}

	// Audit log setelah operasi sukses
	meta, _ := json.Marshal(map[string]string{
		"creation_method": "from_proposal",
		"proposal_id":     proposal.ID.String(),
	})
	s.logAuditWithMeta(ctx, "license", license.ID, "license_created", actorID, actorName, json.RawMessage("{}"), meta)

	return license, nil
}

// Activate mengubah status license menjadi "active".
// Transisi valid: pending → active, trial → active, expired → active (renew).
func (s *LicenseService) Activate(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.Activate find: %w", err)
	}

	if err := validateTransition(license.Status, "active"); err != nil {
		return fmt.Errorf("LicenseService.Activate: %w", err)
	}

	prevStatus := license.Status

	if err := s.repo.UpdateStatus(ctx, id, "active"); err != nil {
		return fmt.Errorf("LicenseService.Activate update: %w", err)
	}

	// Audit log setelah operasi sukses
	changes, _ := json.Marshal(map[string]string{"from": prevStatus, "to": "active"})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "status_changed", actorID, actorName, changes)

	return nil
}

// Suspend mengubah status license menjadi "suspended".
// Transisi valid: active → suspended.
func (s *LicenseService) Suspend(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.Suspend find: %w", err)
	}

	if err := validateTransition(license.Status, "suspended"); err != nil {
		return fmt.Errorf("LicenseService.Suspend: %w", err)
	}

	prevStatus := license.Status

	if err := s.repo.UpdateStatus(ctx, id, "suspended"); err != nil {
		return fmt.Errorf("LicenseService.Suspend update: %w", err)
	}

	// Audit log setelah operasi sukses
	changes, _ := json.Marshal(map[string]string{"from": prevStatus, "to": "suspended"})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "status_changed", actorID, actorName, changes)

	return nil
}

// Renew mengubah status license dari "expired" menjadi "active" dan memperbarui expires_at.
// Jika newExpiresAt nil, license tidak memiliki tanggal kadaluarsa baru (perpetual).
func (s *LicenseService) Renew(ctx context.Context, id uuid.UUID, newExpiresAt *time.Time, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.Renew find: %w", err)
	}

	if err := validateTransition(license.Status, "active"); err != nil {
		return fmt.Errorf("LicenseService.Renew: %w", err)
	}

	license.Status = "active"
	license.ExpiresAt = newExpiresAt
	license.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, license); err != nil {
		return fmt.Errorf("LicenseService.Renew update: %w", err)
	}

	// Audit log setelah operasi sukses
	changes, _ := json.Marshal(map[string]any{
		"from":           "expired",
		"to":             "active",
		"new_expires_at": newExpiresAt,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "license_renewed", actorID, actorName, changes)

	return nil
}

// SetTrial mengubah status license dari "pending" menjadi "trial".
// Hanya bisa dilakukan oleh PO.
func (s *LicenseService) SetTrial(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.SetTrial find: %w", err)
	}

	if err := validateTransition(license.Status, "trial"); err != nil {
		return fmt.Errorf("LicenseService.SetTrial: %w", err)
	}

	prevStatus := license.Status

	if err := s.repo.UpdateStatus(ctx, id, "trial"); err != nil {
		return fmt.Errorf("LicenseService.SetTrial update: %w", err)
	}

	// Audit log setelah operasi sukses
	changes, _ := json.Marshal(map[string]string{"from": prevStatus, "to": "trial"})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "status_changed", actorID, actorName, changes)

	return nil
}

// UpdateConstraints memperbarui constraint license (max_users, modules, apps, dll.).
// Tidak mengubah status license.
func (s *LicenseService) UpdateConstraints(ctx context.Context, id uuid.UUID, req UpdateConstraintsRequest, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.UpdateConstraints find: %w", err)
	}

	// Catat nilai sebelumnya untuk audit
	before, _ := json.Marshal(map[string]any{
		"modules":             license.Modules,
		"apps":                license.Apps,
		"max_users":           license.MaxUsers,
		"max_trans_per_month": license.MaxTransPerMonth,
		"max_trans_per_day":   license.MaxTransPerDay,
		"max_items":           license.MaxItems,
		"max_customers":       license.MaxCustomers,
		"max_branches":        license.MaxBranches,
		"max_storage":         license.MaxStorage,
		"expires_at":          license.ExpiresAt,
		"check_interval":      license.CheckInterval,
	})

	// Terapkan perubahan
	license.Modules = req.Modules
	license.Apps = req.Apps
	license.MaxUsers = req.MaxUsers
	license.MaxTransPerMonth = req.MaxTransPerMonth
	license.MaxTransPerDay = req.MaxTransPerDay
	license.MaxItems = req.MaxItems
	license.MaxCustomers = req.MaxCustomers
	license.MaxBranches = req.MaxBranches
	license.MaxStorage = req.MaxStorage
	license.ExpiresAt = req.ExpiresAt
	if req.CheckInterval != "" {
		license.CheckInterval = req.CheckInterval
	}
	license.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, license); err != nil {
		return fmt.Errorf("LicenseService.UpdateConstraints update: %w", err)
	}

	// Audit log setelah operasi sukses
	after, _ := json.Marshal(map[string]any{
		"modules":             license.Modules,
		"apps":                license.Apps,
		"max_users":           license.MaxUsers,
		"max_trans_per_month": license.MaxTransPerMonth,
		"max_trans_per_day":   license.MaxTransPerDay,
		"max_items":           license.MaxItems,
		"max_customers":       license.MaxCustomers,
		"max_branches":        license.MaxBranches,
		"max_storage":         license.MaxStorage,
		"expires_at":          license.ExpiresAt,
		"check_interval":      license.CheckInterval,
	})
	changes, _ := json.Marshal(map[string]json.RawMessage{
		"before": before,
		"after":  after,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "license_updated", actorID, actorName, changes)

	return nil
}

// UpdateStatus mengubah status license langsung tanpa validasi transisi (superuser only).
func (s *LicenseService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("LicenseService.UpdateStatus: %w", err)
	}
	return nil
}

// ActivateWithSuperuser mengubah status license ke "active" dan membuat superuser di client app.
func (s *LicenseService) ActivateWithSuperuser(ctx context.Context, id uuid.UUID, username, password string, actorID uuid.UUID, actorName string) error {
	return s.changeStatusWithSuperuser(ctx, id, "active", username, password, actorID, actorName)
}

// SetTrialWithSuperuser mengubah status license ke "trial" dan membuat superuser di client app.
func (s *LicenseService) SetTrialWithSuperuser(ctx context.Context, id uuid.UUID, username, password string, actorID uuid.UUID, actorName string) error {
	return s.changeStatusWithSuperuser(ctx, id, "trial", username, password, actorID, actorName)
}

// ResetSuperuser membuat/mereset superuser di client app tanpa mengubah status license.
func (s *LicenseService) ResetSuperuser(ctx context.Context, id uuid.UUID, username, password string, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.ResetSuperuser find: %w", err)
	}

	returnedUsername, err := s.callCreateSuperuser(ctx, license, username, password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateSuperuser(ctx, id, returnedUsername); err != nil {
		return fmt.Errorf("LicenseService.ResetSuperuser update: %w", err)
	}

	changes, _ := json.Marshal(map[string]string{"username": returnedUsername})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "superuser_reset", actorID, actorName, changes)

	return nil
}

// changeStatusWithSuperuser adalah helper untuk activate/trial dengan superuser creation.
func (s *LicenseService) changeStatusWithSuperuser(ctx context.Context, id uuid.UUID, targetStatus, username, password string, actorID uuid.UUID, actorName string) error {
	license, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("LicenseService.changeStatusWithSuperuser find: %w", err)
	}

	if err := validateTransition(license.Status, targetStatus); err != nil {
		return fmt.Errorf("LicenseService.changeStatusWithSuperuser: %w", err)
	}

	prevStatus := license.Status

	// Call client app untuk create superuser
	returnedUsername, err := s.callCreateSuperuser(ctx, license, username, password)
	if err != nil {
		return err
	}

	// Update status
	if err := s.repo.UpdateStatus(ctx, id, targetStatus); err != nil {
		return fmt.Errorf("LicenseService.changeStatusWithSuperuser update status: %w", err)
	}

	// Update superuser username
	if err := s.repo.UpdateSuperuser(ctx, id, returnedUsername); err != nil {
		s.logger.Error("changeStatusWithSuperuser: UpdateSuperuser failed", zap.Error(err))
	}

	changes, _ := json.Marshal(map[string]string{
		"from":     prevStatus,
		"to":       targetStatus,
		"username": returnedUsername,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "license", id, "status_changed_with_superuser", actorID, actorName, changes)

	return nil
}

// callCreateSuperuser memanggil client app untuk membuat superuser.
// Mengirim POST ke {instance_url}/api/v1/create-superuser dengan body {otp, license_key, username, password}.
// Mengembalikan username dari response client app.
func (s *LicenseService) callCreateSuperuser(ctx context.Context, license *domain.ClientLicense, username, password string) (string, error) {
	if license.InstanceURL == nil || *license.InstanceURL == "" {
		return "", fmt.Errorf("LicenseService.callCreateSuperuser: %w", domain.ErrLicenseNoInstanceURL)
	}

	// Ambil OTP aktif
	activeOTP, err := s.otpRepo.GetActive(ctx)
	if err != nil {
		return "", fmt.Errorf("LicenseService.callCreateSuperuser: %w", domain.ErrNoActiveOTP)
	}

	// Construct URL
	instanceURL := strings.TrimRight(*license.InstanceURL, "/")
	callbackURL := instanceURL + "/api/v1/create-superuser"

	// Build request body
	reqBody, _ := json.Marshal(map[string]string{
		"otp":         activeOTP,
		"license_key": license.LicenseKey,
		"username":    username,
		"password":    password,
	})

	// Call client app
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, callbackURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("LicenseService.callCreateSuperuser new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("LicenseService.callCreateSuperuser call failed: %w: %w", domain.ErrSuperuserCreationFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errBody map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return "", fmt.Errorf("LicenseService.callCreateSuperuser: client returned %d: %v: %w", resp.StatusCode, errBody, domain.ErrSuperuserCreationFailed)
	}

	// Parse response — expect {"username": "..."}
	var result struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("LicenseService.callCreateSuperuser decode response: %w", err)
	}

	if result.Username == "" {
		result.Username = username
	}

	return result.Username, nil
}


// validateTransition memvalidasi apakah perubahan status dari `from` ke `to` diizinkan.
// Mengembalikan domain.ErrLicenseInvalidTransition jika transisi tidak valid.
func validateTransition(from, to string) error {
	valid := map[string][]string{
		"pending":   {"active", "trial"},
		"trial":     {"active"},
		"active":    {"suspended", "expired"},
		"suspended": {"active"},
		"expired":   {"active"},
	}

	allowed, ok := valid[from]
	if !ok {
		return fmt.Errorf("status %q: %w", from, domain.ErrLicenseInvalidTransition)
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("transition %q → %q: %w", from, to, domain.ErrLicenseInvalidTransition)
}

// logAuditWithMeta mencatat audit log dengan tambahan metadata.
// Jika gagal, error hanya dicatat ke logger dan tidak dipropagasi.
func (s *LicenseService) logAuditWithMeta(
	ctx context.Context,
	entityType string,
	entityID uuid.UUID,
	action string,
	actorID uuid.UUID,
	actorName string,
	changes json.RawMessage,
	metadata json.RawMessage,
) {
	entry := &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		ActorID:    actorID,
		ActorName:  actorName,
		Changes:    changes,
		Metadata:   metadata,
	}
	if err := s.auditRepo.Create(ctx, entry); err != nil {
		s.logger.Error("audit log failed",
			zap.String("entity_type", entityType),
			zap.String("entity_id", entityID.String()),
			zap.String("action", action),
			zap.Error(err),
		)
	}
}
