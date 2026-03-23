package http

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type DashboardHandler struct {
	db *sqlx.DB
}

func NewDashboardHandler(db *sqlx.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type recentItem struct {
	ID          string    `db:"id"           json:"id"`
	CompanyName string    `db:"client_name"  json:"company_name"`
	ContactName string    `db:"client_email" json:"contact_name"`
	Status      string    `db:"status"       json:"status"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var totalClients int
	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM client_licenses`).Scan(&totalClients)

	var activeCount int
	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM client_licenses WHERE status = 'active'`).Scan(&activeCount)

	var pendingCount int
	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM client_licenses WHERE status = 'trial' AND is_provisioned = false`).Scan(&pendingCount)

	var newThisMonth int
	_ = h.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM client_licenses WHERE created_at >= date_trunc('month', NOW())`).
		Scan(&newThisMonth)

	rows, err := h.db.QueryxContext(ctx,
		`SELECT id::text, client_name, client_email, status, created_at
		 FROM client_licenses ORDER BY created_at DESC LIMIT 5`)

	var recent []recentItem
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item recentItem
			if err := rows.StructScan(&item); err == nil {
				recent = append(recent, item)
			}
		}
	}
	if recent == nil {
		recent = []recentItem{}
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"total_clients":              totalClients,
		"active_companies":           activeCount,
		"paid_invoices_count":        0,
		"mrr":                        0.0,
		"mrr_growth_percent":         0.0,
		"pending_registrations":      pendingCount,
		"new_registrations_this_month": newThisMonth,
		"recent_registrations":       recent,
	})
}
