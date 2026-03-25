//go:build !wasm

package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
)

// ProposalHandler menangani HTTP requests untuk Proposals.
type ProposalHandler struct {
	svc         *service.ProposalService
	compSvc     *service.CompanyService
	projSvc     *service.ProjectService
	prodSvc     *service.ProductService
	userSvc     *service.UserService
	storagePath string
	logger      *zap.Logger
}

// NewProposalHandler membuat instance ProposalHandler baru.
func NewProposalHandler(
	svc *service.ProposalService,
	compSvc *service.CompanyService,
	projSvc *service.ProjectService,
	prodSvc *service.ProductService,
	userSvc *service.UserService,
	cfg *config.Config,
	logger *zap.Logger,
) *ProposalHandler {
	return &ProposalHandler{
		svc:         svc,
		compSvc:     compSvc,
		projSvc:     projSvc,
		prodSvc:     prodSvc,
		userSvc:     userSvc,
		storagePath: cfg.StoragePath,
		logger:      logger,
	}
}

// proposalListItemDTO adalah representasi ringkas proposal untuk daftar.
type proposalListItemDTO struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"project_id"`
	CompanyID   string   `json:"company_id"`
	ProductID   string   `json:"product_id"`
	Version     int      `json:"version"`
	Status      string   `json:"status"`
	Plan        string   `json:"plan"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// proposalDetailDTO adalah representasi lengkap proposal termasuk changelog dan nama entitas terkait.
type proposalDetailDTO struct {
	ID               string           `json:"id"`
	ProjectID        string           `json:"project_id"`
	ProjectName      string           `json:"project_name"`
	CompanyID        string           `json:"company_id"`
	CompanyName      string           `json:"company_name"`
	ProductID        string           `json:"product_id"`
	ProductName      string           `json:"product_name"`
	Version          int              `json:"version"`
	Status           string           `json:"status"`
	Plan             string           `json:"plan"`
	Modules          []string         `json:"modules"`
	Apps             []string         `json:"apps"`
	MaxUsers         *int             `json:"max_users"`
	MaxTransPerMonth *int             `json:"max_trans_per_month"`
	MaxTransPerDay   *int             `json:"max_trans_per_day"`
	MaxItems         *int             `json:"max_items"`
	MaxCustomers     *int             `json:"max_customers"`
	MaxBranches      *int             `json:"max_branches"`
	MaxStorage       *int             `json:"max_storage"`
	ContractAmount   *float64         `json:"contract_amount"`
	ExpiresAt        *string          `json:"expires_at"`
	Notes            string           `json:"notes"`
	OwnerNotes       string           `json:"owner_notes"`
	RejectionReason  string           `json:"rejection_reason"`
	Changelog        *changelogDTO    `json:"changelog"`
	PDFPath          string           `json:"pdf_path"`
	SubmittedByName  string           `json:"submitted_by_name"`
	ReviewedByName   string           `json:"reviewed_by_name"`
	ReviewedAt       *string          `json:"reviewed_at"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
}

// changelogDTO adalah DTO untuk changelog proposal.
type changelogDTO struct {
	ComparedToVersion int                `json:"compared_to_version"`
	Summary           string             `json:"summary"`
	Changes           []changelogEntryDTO `json:"changes"`
	Unchanged         []string           `json:"unchanged"`
}

// changelogEntryDTO adalah satu baris diff changelog.
type changelogEntryDTO struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

// createProposalRequest adalah body JSON untuk POST /api/internal/proposals.
type createProposalRequest struct {
	ProjectID        string   `json:"project_id"`
	CompanyID        string   `json:"company_id"`
	ProductID        string   `json:"product_id"`
	Plan             string   `json:"plan"`
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	MaxUsers         *int     `json:"max_users"`
	MaxTransPerMonth *int     `json:"max_trans_per_month"`
	MaxTransPerDay   *int     `json:"max_trans_per_day"`
	MaxItems         *int     `json:"max_items"`
	MaxCustomers     *int     `json:"max_customers"`
	MaxBranches      *int     `json:"max_branches"`
	MaxStorage       *int     `json:"max_storage"`
	ContractAmount   *float64 `json:"contract_amount"`
	ExpiresAt        *string  `json:"expires_at"`
	Notes            *string  `json:"notes"`
}

// updateProposalRequest adalah body JSON untuk PUT /api/internal/proposals/{id}.
type updateProposalRequest struct {
	Plan             *string  `json:"plan"`
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	MaxUsers         *int     `json:"max_users"`
	MaxTransPerMonth *int     `json:"max_trans_per_month"`
	MaxTransPerDay   *int     `json:"max_trans_per_day"`
	MaxItems         *int     `json:"max_items"`
	MaxCustomers     *int     `json:"max_customers"`
	MaxBranches      *int     `json:"max_branches"`
	MaxStorage       *int     `json:"max_storage"`
	ContractAmount   *float64 `json:"contract_amount"`
	ExpiresAt        *string  `json:"expires_at"`
	Notes            *string  `json:"notes"`
	OwnerNotes       *string  `json:"owner_notes"`
}

// approveProposalRequest adalah body JSON untuk PUT /api/internal/proposals/{id}/approve.
type approveProposalRequest struct {
	OwnerNotes *string `json:"owner_notes"`
}

// rejectProposalRequest adalah body JSON untuk PUT /api/internal/proposals/{id}/reject.
type rejectProposalRequest struct {
	Reason string `json:"reason"`
}

// toListItemDTO mengkonversi domain.Proposal ke proposalListItemDTO.
func toProposalListItemDTO(p *domain.Proposal) proposalListItemDTO {
	return proposalListItemDTO{
		ID:        p.ID.String(),
		ProjectID: p.ProjectID.String(),
		CompanyID: p.CompanyID.String(),
		ProductID: p.ProductID.String(),
		Version:   p.Version,
		Status:    p.Status,
		Plan:      p.Plan,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}

// toDetailDTO mengkonversi domain.Proposal ke proposalDetailDTO dengan data enrich opsional.
func toProposalDetailDTO(p *domain.Proposal, companyName, projectName, productName, submittedByName, reviewedByName string) proposalDetailDTO {
	dto := proposalDetailDTO{
		ID:               p.ID.String(),
		ProjectID:        p.ProjectID.String(),
		ProjectName:      projectName,
		CompanyID:        p.CompanyID.String(),
		CompanyName:      companyName,
		ProductID:        p.ProductID.String(),
		ProductName:      productName,
		Version:          p.Version,
		Status:           p.Status,
		Plan:             p.Plan,
		Modules:          p.Modules,
		Apps:             p.Apps,
		MaxUsers:         p.MaxUsers,
		MaxTransPerMonth: p.MaxTransPerMonth,
		MaxTransPerDay:   p.MaxTransPerDay,
		MaxItems:         p.MaxItems,
		MaxCustomers:     p.MaxCustomers,
		MaxBranches:      p.MaxBranches,
		MaxStorage:       p.MaxStorage,
		ContractAmount:   p.ContractAmount,
		SubmittedByName:  submittedByName,
		ReviewedByName:   reviewedByName,
		CreatedAt:        p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        p.UpdatedAt.Format(time.RFC3339),
	}
	if p.ExpiresAt != nil {
		s := p.ExpiresAt.Format(time.RFC3339)
		dto.ExpiresAt = &s
	}
	if p.Notes != nil {
		dto.Notes = *p.Notes
	}
	if p.OwnerNotes != nil {
		dto.OwnerNotes = *p.OwnerNotes
	}
	if p.RejectionReason != nil {
		dto.RejectionReason = *p.RejectionReason
	}
	if p.PDFPath != nil {
		dto.PDFPath = *p.PDFPath
	}
	if p.ReviewedAt != nil {
		s := p.ReviewedAt.Format(time.RFC3339)
		dto.ReviewedAt = &s
	}

	// Parse changelog jika ada
	if len(p.Changelog) > 0 && string(p.Changelog) != "null" {
		var cl domain.Changelog
		if err := json.Unmarshal(p.Changelog, &cl); err == nil {
			cdto := &changelogDTO{
				ComparedToVersion: cl.ComparedToVersion,
				Summary:           cl.Summary,
				Unchanged:         cl.Unchanged,
			}
			for _, e := range cl.Changes {
				cdto.Changes = append(cdto.Changes, changelogEntryDTO{
					Field:    e.Field,
					OldValue: e.OldValue,
					NewValue: e.NewValue,
				})
			}
			dto.Changelog = cdto
		}
	}

	return dto
}

// List menangani GET /api/internal/proposals.
// Sales hanya melihat proposal miliknya; role lain melihat semua.
func (h *ProposalHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	var (
		proposals []*domain.Proposal
		err       error
	)
	if claims.Role == "sales" {
		submitterID, pErr := parseUUID(claims.Sub)
		if pErr != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user ID")
			return
		}
		proposals, err = h.svc.ListBySubmitter(r.Context(), submitterID)
	} else {
		proposals, err = h.svc.List(r.Context())
	}
	if err != nil {
		h.logger.Error("ProposalHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]proposalListItemDTO, 0, len(proposals))
	for _, p := range proposals {
		dtos = append(dtos, toProposalListItemDTO(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// ListByProject menangani GET /api/internal/projects/{projectID}/proposals.
// Mengembalikan semua proposals untuk project tertentu.
func (h *ProposalHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Project ID tidak valid")
		return
	}

	proposals, err := h.svc.ListByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("ProposalHandler.ListByProject", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]proposalListItemDTO, 0, len(proposals))
	for _, p := range proposals {
		dtos = append(dtos, toProposalListItemDTO(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// GetByID menangani GET /api/internal/proposals/{id}.
// Mengembalikan detail proposal beserta data enrich dari company, project, product, dan user.
// Sales hanya dapat mengakses proposal yang mereka buat.
func (h *ProposalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	proposal, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		h.logger.Error("ProposalHandler.GetByID", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Sales hanya boleh melihat proposal miliknya sendiri.
	if claims.Role == "sales" && proposal.SubmittedBy.String() != claims.Sub {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki akses ke proposal ini")
		return
	}

	// Enrich dengan nama entitas terkait
	companyName := proposal.CompanyID.String()
	if comp, err2 := h.compSvc.GetByID(r.Context(), proposal.CompanyID); err2 == nil {
		companyName = comp.Name
	}

	projectName := proposal.ProjectID.String()
	if proj, err2 := h.projSvc.GetByID(r.Context(), proposal.ProjectID); err2 == nil {
		projectName = proj.Name
	}

	productName := proposal.ProductID.String()
	if prod, err2 := h.prodSvc.GetByID(r.Context(), proposal.ProductID); err2 == nil {
		productName = prod.Name
	}

	submittedByName := proposal.SubmittedBy.String()
	if u, err2 := h.userSvc.GetByID(r.Context(), proposal.SubmittedBy); err2 == nil {
		submittedByName = u.Name
	}

	reviewedByName := ""
	if proposal.ReviewedBy != nil {
		if u, err2 := h.userSvc.GetByID(r.Context(), *proposal.ReviewedBy); err2 == nil {
			reviewedByName = u.Name
		}
	}

	dto := toProposalDetailDTO(proposal, companyName, projectName, productName, submittedByName, reviewedByName)
	writeJSON(w, http.StatusOK, dto)
}

// Create menangani POST /api/internal/proposals.
// Body: {project_id, company_id, product_id, plan, modules, apps, constraints..., notes}
func (h *ProposalHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req createProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	projectID, err := parseUUID(req.ProjectID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "project_id tidak valid")
		return
	}
	companyID, err := parseUUID(req.CompanyID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "company_id tidak valid")
		return
	}
	productID, err := parseUUID(req.ProductID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "product_id tidak valid")
		return
	}

	if req.Plan == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "plan wajib diisi")
		return
	}

	svcReq := service.CreateProposalRequest{
		ProjectID:        projectID,
		CompanyID:        companyID,
		ProductID:        productID,
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
		Notes:            req.Notes,
	}

	if req.ExpiresAt != nil {
		t, parseErr := time.Parse(time.RFC3339, *req.ExpiresAt)
		if parseErr != nil {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "format expires_at tidak valid (gunakan RFC3339)")
			return
		}
		svcReq.ExpiresAt = &t
	}

	proposal, err := h.svc.Create(r.Context(), svcReq, actorID, claims.Name)
	if err != nil {
		h.logger.Error("ProposalHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toProposalListItemDTO(proposal))
}

// Update menangani PUT /api/internal/proposals/{id}.
// Sales hanya bisa edit draft; PO/superuser bisa edit draft atau submitted.
func (h *ProposalHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req updateProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	svcReq := service.UpdateProposalRequest{
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
		Notes:            req.Notes,
		OwnerNotes:       req.OwnerNotes,
	}

	if req.ExpiresAt != nil {
		t, parseErr := time.Parse(time.RFC3339, *req.ExpiresAt)
		if parseErr != nil {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "format expires_at tidak valid (gunakan RFC3339)")
			return
		}
		svcReq.ExpiresAt = &t
	}

	proposal, err := h.svc.Update(r.Context(), id, svcReq, claims.Role, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProposalNotDraft) {
			writeError(w, http.StatusUnprocessableEntity, "PROPOSAL_NOT_DRAFT", "Hanya proposal draft yang bisa diedit oleh sales")
			return
		}
		h.logger.Error("ProposalHandler.Update", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toProposalListItemDTO(proposal))
}

// Submit menangani PUT /api/internal/proposals/{id}/submit.
// Mengubah status draft → submitted.
func (h *ProposalHandler) Submit(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	if err := h.svc.Submit(r.Context(), id, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProposalNotDraft) {
			writeError(w, http.StatusUnprocessableEntity, "PROPOSAL_NOT_DRAFT", "Hanya proposal dengan status draft yang bisa di-submit")
			return
		}
		h.logger.Error("ProposalHandler.Submit", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Proposal berhasil di-submit"})
}

// Approve menangani PUT /api/internal/proposals/{id}/approve.
// Hanya untuk project_owner dan superuser.
func (h *ProposalHandler) Approve(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	// Hanya PO dan superuser
	if claims.Role != "project_owner" && claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya project_owner atau superuser yang bisa menyetujui proposal")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req approveProposalRequest
	// Body opsional — abaikan decode error
	_ = json.NewDecoder(r.Body).Decode(&req)

	license, err := h.svc.Approve(r.Context(), id, req.OwnerNotes, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProposalNotSubmitted) {
			writeError(w, http.StatusUnprocessableEntity, "PROPOSAL_NOT_SUBMITTED", "Hanya proposal dengan status submitted yang bisa disetujui")
			return
		}
		h.logger.Error("ProposalHandler.Approve", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message":    "Proposal disetujui dan lisensi telah dibuat",
		"license_id": license.ID.String(),
	})
}

// Reject menangani PUT /api/internal/proposals/{id}/reject.
// Hanya untuk project_owner dan superuser.
func (h *ProposalHandler) Reject(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	// Hanya PO dan superuser
	if claims.Role != "project_owner" && claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya project_owner atau superuser yang bisa menolak proposal")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req rejectProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Reason == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "reason wajib diisi")
		return
	}

	if err := h.svc.Reject(r.Context(), id, req.Reason, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProposalNotSubmitted) {
			writeError(w, http.StatusUnprocessableEntity, "PROPOSAL_NOT_SUBMITTED", "Hanya proposal dengan status submitted yang bisa ditolak")
			return
		}
		h.logger.Error("ProposalHandler.Reject", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Proposal ditolak"})
}

// GetPDF menangani GET /api/internal/proposals/{id}/pdf.
// Melayani file PDF proposal yang sudah disetujui.
func (h *ProposalHandler) GetPDF(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	proposal, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProposalNotFound) {
			writeError(w, http.StatusNotFound, "PROPOSAL_NOT_FOUND", "Proposal tidak ditemukan")
			return
		}
		h.logger.Error("ProposalHandler.GetPDF", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	if proposal.PDFPath == nil || *proposal.PDFPath == "" {
		writeError(w, http.StatusNotFound, "PDF_NOT_FOUND", "PDF belum tersedia untuk proposal ini")
		return
	}

	// Tentukan path absolut PDF
	pdfPath := *proposal.PDFPath
	if !filepath.IsAbs(pdfPath) {
		pdfPath = filepath.Join(h.storagePath, pdfPath)
	}

	f, err := os.Open(pdfPath)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, http.StatusNotFound, "PDF_NOT_FOUND", "File PDF tidak ditemukan di server")
			return
		}
		h.logger.Error("ProposalHandler.GetPDF open file", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuka file PDF")
		return
	}
	defer f.Close()

	filename := filepath.Base(pdfPath)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, f)
}
