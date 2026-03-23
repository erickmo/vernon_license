package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	createclientlicense "github.com/flashlab/flasherp-developer-api/internal/command/create_client_license"
	provisionlicense "github.com/flashlab/flasherp-developer-api/internal/command/provision_license"
	updatelicenseconstraints "github.com/flashlab/flasherp-developer-api/internal/command/update_license_constraints"
	updatelicensestatus "github.com/flashlab/flasherp-developer-api/internal/command/update_license_status"
	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
	getclientlicense "github.com/flashlab/flasherp-developer-api/internal/query/get_client_license"
	listclientlicenses "github.com/flashlab/flasherp-developer-api/internal/query/list_client_licenses"
	"github.com/flashlab/flasherp-developer-api/pkg/commandbus"
	"github.com/flashlab/flasherp-developer-api/pkg/middleware"
	"github.com/flashlab/flasherp-developer-api/pkg/querybus"
)

type LicenseHandler struct {
	cmdBus                     *commandbus.CommandBus
	queryBus                   *querybus.QueryBus
	createHandler              *createclientlicense.Handler
	provisionHandler           *provisionlicense.Handler
}

func NewLicenseHandler(
	cmdBus *commandbus.CommandBus,
	queryBus *querybus.QueryBus,
	createHandler *createclientlicense.Handler,
	provisionHandler *provisionlicense.Handler,
) *LicenseHandler {
	return &LicenseHandler{
		cmdBus:           cmdBus,
		queryBus:         queryBus,
		createHandler:    createHandler,
		provisionHandler: provisionHandler,
	}
}

func (h *LicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page := 1
	pageSize := 20
	if v := q.Get("page"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := q.Get("page_size"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			pageSize = n
		}
	}

	result, err := querybus.Dispatch[*listclientlicenses.ListClientLicensesResult](
		r.Context(), h.queryBus,
		listclientlicenses.ListClientLicenses{
			Status:   q.Get("status"),
			Product:  q.Get("product"),
			Search:   q.Get("search"),
			Page:     page,
			PageSize: pageSize,
		},
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "gagal ambil daftar lisensi")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"items":     result.Items,
		"total":     result.Total,
		"page":      result.Page,
		"page_size": result.PageSize,
	})
}

func (h *LicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	var req struct {
		ClientName       string     `json:"client_name"`
		ClientEmail      string     `json:"client_email"`
		Product          string     `json:"product"`
		Plan             string     `json:"plan"`
		MaxUsers         *int       `json:"max_users"`
		MaxTransPerMonth *int       `json:"max_trans_per_month"`
		MaxTransPerDay   *int       `json:"max_trans_per_day"`
		MaxItems         *int       `json:"max_items"`
		MaxCustomers     *int       `json:"max_customers"`
		MaxBranches      *int       `json:"max_branches"`
		ExpiresAt        *time.Time `json:"expires_at"`
		FlashERPURL      *string    `json:"flasherp_url"`
		ProvisionAPIKey  *string    `json:"provision_api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	if req.ClientName == "" || req.ClientEmail == "" {
		respondError(w, http.StatusBadRequest, "client_name dan client_email wajib diisi")
		return
	}

	createdBy, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "user ID tidak valid")
		return
	}

	plan := req.Plan
	if plan == "" {
		plan = clientlicense.PlanSaaS
	}

	result, err := h.createHandler.HandleCreate(r.Context(), createclientlicense.CreateClientLicense{
		ClientName:       req.ClientName,
		ClientEmail:      req.ClientEmail,
		Product:          req.Product,
		Plan:             plan,
		MaxUsers:         req.MaxUsers,
		MaxTransPerMonth: req.MaxTransPerMonth,
		MaxTransPerDay:   req.MaxTransPerDay,
		MaxItems:         req.MaxItems,
		MaxCustomers:     req.MaxCustomers,
		MaxBranches:      req.MaxBranches,
		ExpiresAt:        req.ExpiresAt,
		FlashERPURL:      req.FlashERPURL,
		ProvisionAPIKey:  req.ProvisionAPIKey,
		CreatedBy:        createdBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "gagal buat lisensi: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"license_key":      result.LicenseKey,
		"client_email":     result.ClientEmail,
		"initial_password": result.InitialPassword,
		"warning":          "Simpan initial_password sekarang. Password ini hanya ditampilkan sekali dan tidak disimpan di sistem.",
	})
}

func (h *LicenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := querybus.Dispatch[*clientlicense.ClientLicense](
		r.Context(), h.queryBus,
		getclientlicense.GetClientLicense{ID: id},
	)
	if err != nil {
		if errors.Is(err, clientlicense.ErrNotFound) {
			respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, "terjadi kesalahan sistem")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *LicenseHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "status wajib diisi")
		return
	}

	// Resolve license key dari id (bisa UUID atau license key)
	licenseKey, err := h.resolveLicenseKey(r, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
		return
	}

	if err := h.cmdBus.Dispatch(r.Context(), updatelicensestatus.UpdateLicenseStatus{
		LicenseKey: licenseKey,
		Status:     req.Status,
	}); err != nil {
		if errors.Is(err, clientlicense.ErrNotFound) {
			respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, "gagal update status: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "status berhasil diperbarui"})
}

func (h *LicenseHandler) UpdateConstraints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		MaxUsers         *int       `json:"max_users"`
		MaxTransPerMonth *int       `json:"max_trans_per_month"`
		MaxTransPerDay   *int       `json:"max_trans_per_day"`
		MaxItems         *int       `json:"max_items"`
		MaxCustomers     *int       `json:"max_customers"`
		MaxBranches      *int       `json:"max_branches"`
		ExpiresAt        *time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	licenseKey, err := h.resolveLicenseKey(r, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
		return
	}

	if err := h.cmdBus.Dispatch(r.Context(), updatelicenseconstraints.UpdateLicenseConstraints{
		LicenseKey:       licenseKey,
		MaxUsers:         req.MaxUsers,
		MaxTransPerMonth: req.MaxTransPerMonth,
		MaxTransPerDay:   req.MaxTransPerDay,
		MaxItems:         req.MaxItems,
		MaxCustomers:     req.MaxCustomers,
		MaxBranches:      req.MaxBranches,
		ExpiresAt:        req.ExpiresAt,
	}); err != nil {
		if errors.Is(err, clientlicense.ErrNotFound) {
			respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, "gagal update constraints: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "constraints berhasil diperbarui"})
}

func (h *LicenseHandler) Provision(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	licenseKey, err := h.resolveLicenseKey(r, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
		return
	}

	if err := h.provisionHandler.Handle(r.Context(), provisionlicense.ProvisionLicense{
		LicenseKey: licenseKey,
	}); err != nil {
		if errors.Is(err, clientlicense.ErrNotFound) {
			respondError(w, http.StatusNotFound, "lisensi tidak ditemukan")
			return
		}
		if errors.Is(err, clientlicense.ErrMissingFlashERPURL) {
			respondError(w, http.StatusBadRequest, "URL FlashERP belum dikonfigurasi")
			return
		}
		respondError(w, http.StatusBadGateway, "gagal provisioning: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "provisioning berhasil dilakukan"})
}

// resolveLicenseKey mengambil license_key dari id (UUID atau license key langsung)
func (h *LicenseHandler) resolveLicenseKey(r *http.Request, id string) (string, error) {
	result, err := querybus.Dispatch[*clientlicense.ClientLicense](
		r.Context(), h.queryBus,
		getclientlicense.GetClientLicense{ID: id},
	)
	if err != nil {
		return "", err
	}
	return result.LicenseKey, nil
}

func parseInt(s string) (int, error) {
	var n int
	_, err := parseIntVal(s, &n)
	return n, err
}

func parseIntVal(s string, n *int) (int, error) {
	val := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("bukan angka")
		}
		val = val*10 + int(c-'0')
	}
	*n = val
	return val, nil
}
