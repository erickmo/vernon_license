package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	proposalpkg "github.com/flashlab/vernon-license/pkg/proposal"
)

// ProposalService menyediakan business logic untuk lifecycle proposal.
type ProposalService struct {
	repo        domain.ProposalRepository
	licenseRepo domain.LicenseRepository
	auditRepo   domain.AuditLogRepository
	notifRepo   domain.NotificationRepository
	userRepo    domain.UserRepository
	licenseSvc  *LicenseService
	storagePath string
	logger      *zap.Logger
}

// NewProposalService membuat instance ProposalService baru dengan semua dependency yang dibutuhkan.
func NewProposalService(
	repo domain.ProposalRepository,
	licenseRepo domain.LicenseRepository,
	auditRepo domain.AuditLogRepository,
	notifRepo domain.NotificationRepository,
	userRepo domain.UserRepository,
	licenseSvc *LicenseService,
	cfg *config.Config,
	logger *zap.Logger,
) *ProposalService {
	return &ProposalService{
		repo:        repo,
		licenseRepo: licenseRepo,
		auditRepo:   auditRepo,
		notifRepo:   notifRepo,
		userRepo:    userRepo,
		licenseSvc:  licenseSvc,
		storagePath: cfg.StoragePath,
		logger:      logger,
	}
}

// CreateProposalRequest adalah input untuk membuat proposal baru.
type CreateProposalRequest struct {
	ProjectID        uuid.UUID  `json:"project_id"`
	CompanyID        uuid.UUID  `json:"company_id"`
	ProductID        uuid.UUID  `json:"product_id"`
	Plan             string     `json:"plan"`
	Modules          []string   `json:"modules"`
	Apps             []string   `json:"apps"`
	MaxUsers         *int       `json:"max_users"`
	MaxTransPerMonth *int       `json:"max_trans_per_month"`
	MaxTransPerDay   *int       `json:"max_trans_per_day"`
	MaxItems         *int       `json:"max_items"`
	MaxCustomers     *int       `json:"max_customers"`
	MaxBranches      *int       `json:"max_branches"`
	MaxStorage       *int       `json:"max_storage"`
	ContractAmount   *float64   `json:"contract_amount"`
	ExpiresAt        *time.Time `json:"expires_at"`
	Notes            *string    `json:"notes"`
}

// UpdateProposalRequest adalah input untuk memperbarui proposal.
// Semua field optional — hanya field yang di-set yang akan diubah.
type UpdateProposalRequest struct {
	Plan             *string    `json:"plan"`
	Modules          []string   `json:"modules"`
	Apps             []string   `json:"apps"`
	MaxUsers         *int       `json:"max_users"`
	MaxTransPerMonth *int       `json:"max_trans_per_month"`
	MaxTransPerDay   *int       `json:"max_trans_per_day"`
	MaxItems         *int       `json:"max_items"`
	MaxCustomers     *int       `json:"max_customers"`
	MaxBranches      *int       `json:"max_branches"`
	MaxStorage       *int       `json:"max_storage"`
	ContractAmount   *float64   `json:"contract_amount"`
	ExpiresAt        *time.Time `json:"expires_at"`
	Notes            *string    `json:"notes"`
	OwnerNotes       *string    `json:"owner_notes"`
}

// List mengembalikan semua proposal (semua role dapat mengakses).
func (s *ProposalService) List(ctx context.Context) ([]*domain.Proposal, error) {
	proposals, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.List: %w", err)
	}
	return proposals, nil
}

// ListBySubmitter mengembalikan semua proposal yang dibuat oleh submitterID.
// Digunakan untuk role sales yang hanya bisa melihat proposal miliknya.
func (s *ProposalService) ListBySubmitter(ctx context.Context, submitterID uuid.UUID) ([]*domain.Proposal, error) {
	proposals, err := s.repo.FindBySubmitter(ctx, submitterID)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.ListBySubmitter: %w", err)
	}
	return proposals, nil
}

// ListByProject mengembalikan semua proposal untuk project tertentu.
func (s *ProposalService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.Proposal, error) {
	proposals, err := s.repo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.ListByProject: %w", err)
	}
	return proposals, nil
}

// GetByID mengambil satu proposal berdasarkan UUID-nya.
func (s *ProposalService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Proposal, error) {
	proposal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.GetByID: %w", err)
	}
	return proposal, nil
}

// Create membuat draft proposal baru (sales atau PO).
// Version dihitung otomatis via NextVersion.
// Status awal adalah "draft".
// Changelog dihitung otomatis terhadap versi sebelumnya jika ada.
// Audit: action="proposal_created".
func (s *ProposalService) Create(ctx context.Context, req CreateProposalRequest, actorID uuid.UUID, actorName string) (*domain.Proposal, error) {
	// Hitung version berikutnya
	version, err := s.repo.NextVersion(ctx, req.ProjectID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.Create next version: %w", err)
	}

	now := time.Now().UTC()
	p := &domain.Proposal{
		ID:               uuid.New(),
		ProjectID:        req.ProjectID,
		CompanyID:        req.CompanyID,
		ProductID:        req.ProductID,
		Version:          version,
		Status:           "draft",
		Plan:             req.Plan,
		Modules:          req.Modules,
		Apps:             req.Apps,
		MaxUsers:         req.MaxUsers,
		MaxTransPerMonth: req.MaxTransPerMonth,
		MaxTransPerDay:   req.MaxTransPerDay,
		MaxItems:         req.MaxItems,
		MaxCustomers:     req.MaxCustomers,
		MaxBranches:      req.MaxBranches,
		MaxStorage:       req.MaxStorage,
		ContractAmount:   req.ContractAmount,
		ExpiresAt:        req.ExpiresAt,
		Notes:            req.Notes,
		SubmittedBy:      actorID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Compute changelog vs previous version jika ada
	if version > 1 {
		prev, err := s.repo.FindLatestByProjectProduct(ctx, req.ProjectID, req.ProductID)
		if err == nil {
			cl := domain.ComputeChangelog(prev, p)
			raw, jsonErr := json.Marshal(cl)
			if jsonErr == nil {
				p.Changelog = raw
			}
		}
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("ProposalService.Create persist: %w", err)
	}

	// Audit log setelah operasi sukses
	changes, _ := json.Marshal(map[string]any{
		"version": version,
		"status":  "draft",
		"plan":    req.Plan,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "proposal", p.ID, "proposal_created", actorID, actorName, changes)

	return p, nil
}

// Update memperbarui proposal.
// Role "sales": hanya bisa update jika status == "draft" → ErrProposalNotDraft.
// Role "project_owner"/"superuser": bisa update draft atau submitted.
// Changelog direcompute setelah update.
// Jika PO edit submitted → audit: "proposal_edited_by_owner" dan notifikasi ke submitter.
func (s *ProposalService) Update(ctx context.Context, id uuid.UUID, req UpdateProposalRequest, actorRole string, actorID uuid.UUID, actorName string) (*domain.Proposal, error) {
	proposal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.Update find: %w", err)
	}

	// Sales hanya bisa update draft
	if actorRole == "sales" && proposal.Status != "draft" {
		return nil, fmt.Errorf("ProposalService.Update: %w", domain.ErrProposalNotDraft)
	}

	// PO/superuser bisa update draft atau submitted; semua status lain tidak boleh
	if actorRole != "sales" && proposal.Status != "draft" && proposal.Status != "submitted" {
		return nil, fmt.Errorf("ProposalService.Update: status %q tidak bisa diedit", proposal.Status)
	}

	editedByOwner := (actorRole == "project_owner" || actorRole == "superuser") && proposal.Status == "submitted"

	// Snapshot sebelum perubahan untuk changelog
	prev := *proposal

	// Terapkan perubahan (hanya field yang di-set)
	if req.Plan != nil {
		proposal.Plan = *req.Plan
	}
	if req.Modules != nil {
		proposal.Modules = req.Modules
	}
	if req.Apps != nil {
		proposal.Apps = req.Apps
	}
	if req.MaxUsers != nil {
		proposal.MaxUsers = req.MaxUsers
	}
	if req.MaxTransPerMonth != nil {
		proposal.MaxTransPerMonth = req.MaxTransPerMonth
	}
	if req.MaxTransPerDay != nil {
		proposal.MaxTransPerDay = req.MaxTransPerDay
	}
	if req.MaxItems != nil {
		proposal.MaxItems = req.MaxItems
	}
	if req.MaxCustomers != nil {
		proposal.MaxCustomers = req.MaxCustomers
	}
	if req.MaxBranches != nil {
		proposal.MaxBranches = req.MaxBranches
	}
	if req.MaxStorage != nil {
		proposal.MaxStorage = req.MaxStorage
	}
	if req.ContractAmount != nil {
		proposal.ContractAmount = req.ContractAmount
	}
	if req.ExpiresAt != nil {
		proposal.ExpiresAt = req.ExpiresAt
	}
	if req.Notes != nil {
		proposal.Notes = req.Notes
	}
	if req.OwnerNotes != nil {
		proposal.OwnerNotes = req.OwnerNotes
	}

	// Recompute changelog vs snapshot sebelumnya
	cl := domain.ComputeChangelog(&prev, proposal)
	raw, jsonErr := json.Marshal(cl)
	if jsonErr == nil {
		proposal.Changelog = raw
	}

	proposal.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, proposal); err != nil {
		return nil, fmt.Errorf("ProposalService.Update persist: %w", err)
	}

	// Audit dan notifikasi
	auditAction := "proposal_updated"
	if editedByOwner {
		auditAction = "proposal_edited_by_owner"
		// Notifikasi ke submitter (sales)
		s.sendNotification(ctx, proposal.SubmittedBy, "proposal_edited",
			"Proposal diubah oleh PO",
			fmt.Sprintf("Proposal v%d untuk project Anda telah diubah oleh Project Owner.", proposal.Version),
			map[string]string{"proposal_id": id.String()},
		)
	}

	changes, _ := json.Marshal(map[string]any{
		"proposal_id": id.String(),
		"actor_role":  actorRole,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "proposal", id, auditAction, actorID, actorName, changes)

	return proposal, nil
}

// Submit mengubah status dari "draft" → "submitted".
// Hanya bisa jika status == "draft" → ErrProposalNotDraft.
// Notifikasi ke semua PO dan superuser.
// Audit: "proposal_submitted".
func (s *ProposalService) Submit(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	proposal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProposalService.Submit find: %w", err)
	}

	if proposal.Status != "draft" {
		return fmt.Errorf("ProposalService.Submit: %w", domain.ErrProposalNotDraft)
	}

	proposal.Status = "submitted"
	proposal.SubmittedBy = actorID
	proposal.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, proposal); err != nil {
		return fmt.Errorf("ProposalService.Submit persist: %w", err)
	}

	// Notifikasi ke semua PO dan superuser
	users, userErr := s.userRepo.FindAll(ctx)
	if userErr == nil {
		for _, u := range users {
			if (u.Role == "project_owner" || u.Role == "superuser") && u.IsActive {
				s.sendNotification(ctx, u.ID, "proposal_submitted",
					"Proposal baru menunggu review",
					fmt.Sprintf("Proposal v%d dari %s menunggu persetujuan Anda.", proposal.Version, actorName),
					map[string]string{"proposal_id": id.String()},
				)
			}
		}
	} else {
		s.logger.Warn("ProposalService.Submit: failed to load users for notification",
			zap.String("proposal_id", id.String()),
			zap.Error(userErr),
		)
	}

	// Audit log
	changes, _ := json.Marshal(map[string]string{
		"from": "draft",
		"to":   "submitted",
	})
	LogAudit(ctx, s.auditRepo, s.logger, "proposal", id, "proposal_submitted", actorID, actorName, changes)

	return nil
}

// Approve menyetujui proposal (PO/superuser).
// Hanya bisa jika status == "submitted" → ErrProposalNotSubmitted.
// Menghasilkan PDF, menyimpannya, lalu membuat license via licenseSvc.CreateFromProposal.
// Notifikasi ke submitter (sales) dan audit "proposal_approved".
func (s *ProposalService) Approve(ctx context.Context, id uuid.UUID, ownerNotes *string, actorID uuid.UUID, actorName string) (*domain.ClientLicense, error) {
	proposal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.Approve find: %w", err)
	}

	if proposal.Status != "submitted" {
		return nil, fmt.Errorf("ProposalService.Approve: %w", domain.ErrProposalNotSubmitted)
	}

	if ownerNotes != nil {
		proposal.OwnerNotes = ownerNotes
	}

	now := time.Now().UTC()
	proposal.Status = "approved"
	proposal.ReviewedBy = &actorID
	proposal.ReviewedAt = &now

	// Generate PDF
	pdfData := proposalpkg.PDFData{
		Proposal:     proposal,
		CompanyName:  proposal.CompanyID.String(), // akan di-enrich di handler jika perlu
		ProjectName:  proposal.ProjectID.String(),
		ProductName:  proposal.ProductID.String(),
		ReviewerName: actorName,
		VendorName:   "Vernon License",
		VendorAddress: "",
		VendorPhone:  "",
		VendorEmail:  "",
	}

	pdfContent, pdfErr := proposalpkg.GeneratePDF(pdfData)
	if pdfErr != nil {
		s.logger.Warn("ProposalService.Approve: failed to generate PDF",
			zap.String("proposal_id", id.String()),
			zap.Error(pdfErr),
		)
	} else {
		pdfPath, saveErr := proposalpkg.SavePDF(s.storagePath, id.String(), pdfContent)
		if saveErr != nil {
			s.logger.Warn("ProposalService.Approve: failed to save PDF",
				zap.String("proposal_id", id.String()),
				zap.Error(saveErr),
			)
		} else {
			proposal.PDFPath = &pdfPath
			proposal.PDFGeneratedAt = &now
		}
	}

	if err := s.repo.Update(ctx, proposal); err != nil {
		return nil, fmt.Errorf("ProposalService.Approve persist: %w", err)
	}

	// Buat license dari proposal yang sudah disetujui
	license, err := s.licenseSvc.CreateFromProposal(ctx, proposal, actorID, actorName)
	if err != nil {
		return nil, fmt.Errorf("ProposalService.Approve create license: %w", err)
	}

	// Notifikasi ke submitter (sales)
	s.sendNotification(ctx, proposal.SubmittedBy, "proposal_approved",
		"Proposal disetujui",
		fmt.Sprintf("Proposal v%d Anda telah disetujui oleh %s. Lisensi telah dibuat.", proposal.Version, actorName),
		map[string]string{
			"proposal_id": id.String(),
			"license_id":  license.ID.String(),
		},
	)

	// Audit log
	changes, _ := json.Marshal(map[string]any{
		"from":       "submitted",
		"to":         "approved",
		"license_id": license.ID.String(),
	})
	LogAudit(ctx, s.auditRepo, s.logger, "proposal", id, "proposal_approved", actorID, actorName, changes)

	return license, nil
}

// Reject menolak proposal (PO/superuser).
// Hanya bisa jika status == "submitted" → ErrProposalNotSubmitted.
// Set rejection_reason, reviewed_by, reviewed_at, status = "rejected".
// Notifikasi ke submitter dan audit "proposal_rejected".
func (s *ProposalService) Reject(ctx context.Context, id uuid.UUID, reason string, actorID uuid.UUID, actorName string) error {
	proposal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProposalService.Reject find: %w", err)
	}

	if proposal.Status != "submitted" {
		return fmt.Errorf("ProposalService.Reject: %w", domain.ErrProposalNotSubmitted)
	}

	now := time.Now().UTC()
	proposal.Status = "rejected"
	proposal.RejectionReason = &reason
	proposal.ReviewedBy = &actorID
	proposal.ReviewedAt = &now
	proposal.UpdatedAt = now

	if err := s.repo.Update(ctx, proposal); err != nil {
		return fmt.Errorf("ProposalService.Reject persist: %w", err)
	}

	// Notifikasi ke submitter (sales)
	s.sendNotification(ctx, proposal.SubmittedBy, "proposal_rejected",
		"Proposal ditolak",
		fmt.Sprintf("Proposal v%d Anda ditolak oleh %s. Alasan: %s", proposal.Version, actorName, reason),
		map[string]string{
			"proposal_id":      id.String(),
			"rejection_reason": reason,
		},
	)

	// Audit log
	changes, _ := json.Marshal(map[string]string{
		"from":             "submitted",
		"to":               "rejected",
		"rejection_reason": reason,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "proposal", id, "proposal_rejected", actorID, actorName, changes)

	return nil
}

// sendNotification mengirim notifikasi secara fire-and-forget.
// Error hanya dicatat ke logger dan tidak dipropagasi.
func (s *ProposalService) sendNotification(ctx context.Context, userID uuid.UUID, notifType, title, body string, data map[string]string) {
	raw, _ := json.Marshal(data)
	n := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Data:      raw,
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.notifRepo.Create(ctx, n); err != nil {
		s.logger.Warn("ProposalService: failed to send notification",
			zap.String("user_id", userID.String()),
			zap.String("type", notifType),
			zap.Error(err),
		)
	}
}
