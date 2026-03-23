package list_client_licenses

import (
	"context"
	"fmt"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type ListClientLicenses struct {
	Status   string
	Product  string
	Search   string
	Page     int
	PageSize int
}

type ListClientLicensesResult struct {
	Items    []*clientlicense.ClientLicense
	Total    int
	Page     int
	PageSize int
}

type Handler struct {
	licenseRepo clientlicense.ReadRepository
}

func NewHandler(licenseRepo clientlicense.ReadRepository) *Handler {
	return &Handler{licenseRepo: licenseRepo}
}

func (h *Handler) Handle(ctx context.Context, q any) (any, error) {
	query, ok := q.(ListClientLicenses)
	if !ok {
		return nil, fmt.Errorf("query tidak valid")
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	filter := clientlicense.ListFilter{
		Status:   query.Status,
		Product:  query.Product,
		Search:   query.Search,
		Page:     page,
		PageSize: pageSize,
	}

	items, total, err := h.licenseRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil daftar lisensi: %w", err)
	}

	return &ListClientLicensesResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
