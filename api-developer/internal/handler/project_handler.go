//go:build !wasm

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
)

// ProjectHandler menangani HTTP requests untuk Projects.
type ProjectHandler struct {
	svc    *service.ProjectService
	logger *zap.Logger
}

// NewProjectHandler membuat instance ProjectHandler baru.
func NewProjectHandler(svc *service.ProjectService, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{svc: svc, logger: logger}
}

// projectDTO adalah representasi Project yang dikembalikan ke client.
type projectDTO struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// createProjectRequest adalah body JSON untuk POST /api/internal/companies/{companyID}/projects.
type createProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// updateProjectRequest adalah body JSON untuk PUT /api/internal/projects/{id}.
type updateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// toProjectDTO mengkonversi domain.Project ke projectDTO.
func toProjectDTO(p *domain.Project) projectDTO {
	return projectDTO{
		ID:          p.ID.String(),
		CompanyID:   p.CompanyID.String(),
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
	}
}

// ListByCompany menangani GET /api/internal/companies/{companyID}/projects.
// Mengembalikan semua projects milik sebuah company.
func (h *ProjectHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyID, err := parseUUID(chi.URLParam(r, "companyID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Company ID tidak valid")
		return
	}

	projects, err := h.svc.ListByCompany(r.Context(), companyID)
	if err != nil {
		h.logger.Error("ProjectHandler.ListByCompany", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]projectDTO, 0, len(projects))
	for _, p := range projects {
		dtos = append(dtos, toProjectDTO(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// GetByID menangani GET /api/internal/projects/{id}.
// Mengembalikan satu project berdasarkan UUID.
func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	project, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project tidak ditemukan")
			return
		}
		h.logger.Error("ProjectHandler.GetByID", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toProjectDTO(project))
}

// Create menangani POST /api/internal/companies/{companyID}/projects.
// Body: {name, description}
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	companyID, err := parseUUID(chi.URLParam(r, "companyID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Company ID tidak valid")
		return
	}

	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name wajib diisi")
		return
	}

	project, err := h.svc.Create(r.Context(), service.CreateProjectRequest{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("ProjectHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toProjectDTO(project))
}

// Update menangani PUT /api/internal/projects/{id}.
// Body: {name, description, status}
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req updateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name wajib diisi")
		return
	}

	project, err := h.svc.Update(r.Context(), id, service.UpdateProjectRequest{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("ProjectHandler.Update", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toProjectDTO(project))
}

// Delete menangani DELETE /api/internal/projects/{id}.
// Melakukan soft-delete pada project.
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.Delete(r.Context(), id, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrProjectNotFound) {
			writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project tidak ditemukan")
			return
		}
		h.logger.Error("ProjectHandler.Delete", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Project berhasil dihapus"})
}
