//go:build !wasm

package handler

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
)

// DashboardHandler menangani HTTP requests untuk Dashboard stats.
type DashboardHandler struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewDashboardHandler membuat instance DashboardHandler baru.
func NewDashboardHandler(db *sqlx.DB, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{db: db, logger: logger}
}

// ExpiringLicense merepresentasikan lisensi yang akan segera expired.
type ExpiringLicense struct {
	ID         string `json:"id" db:"id"`
	LicenseKey string `json:"license_key" db:"license_key"`
	Company    string `json:"company" db:"company"`
	ExpiresAt  string `json:"expires_at"`
	DaysLeft   int    `json:"days_left" db:"days_left"`
}

// expiringLicenseRow adalah struct untuk scan dari database.
type expiringLicenseRow struct {
	ID         string    `db:"id"`
	LicenseKey string    `db:"license_key"`
	Company    string    `db:"company"`
	ExpiresAt  time.Time `db:"expires_at"`
	DaysLeft   int       `db:"days_left"`
}

// ActivityItem merepresentasikan satu item aktivitas terbaru.
type ActivityItem struct {
	EntityType string `json:"entity_type" db:"entity_type"`
	Action     string `json:"action" db:"action"`
	ActorName  string `json:"actor_name" db:"actor_name"`
	CreatedAt  string `json:"created_at"`
}

// activityItemRow adalah struct untuk scan dari database.
type activityItemRow struct {
	EntityType string    `db:"entity_type"`
	Action     string    `db:"action"`
	ActorName  string    `db:"actor_name"`
	CreatedAt  time.Time `db:"created_at"`
}

// DashboardStats adalah agregasi statistik untuk halaman dashboard.
type DashboardStats struct {
	TotalLicenses     int               `json:"total_licenses"`
	ActiveLicenses    int               `json:"active_licenses"`
	PendingLicenses   int               `json:"pending_licenses"`
	SuspendedLicenses int               `json:"suspended_licenses"`
	ExpiredLicenses   int               `json:"expired_licenses"`
	TotalCompanies    int               `json:"total_companies"`
	TotalProposals    int               `json:"total_proposals"`
	PendingProposals  int               `json:"pending_proposals"`
	TotalRevenue      float64           `json:"total_revenue"`
	ExpiringLicenses  []ExpiringLicense `json:"expiring_licenses"`
	RecentActivity    []ActivityItem    `json:"recent_activity"`
}

// GetStats menangani GET /api/internal/dashboard.
// Mengembalikan statistik agregasi untuk Vernon App dashboard.
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	_, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	ctx := r.Context()
	stats := DashboardStats{
		ExpiringLicenses: []ExpiringLicense{},
		RecentActivity:   []ActivityItem{},
	}

	// License counts per status
	const licenseCountQ = `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'active') AS active,
			COUNT(*) FILTER (WHERE status = 'pending') AS pending,
			COUNT(*) FILTER (WHERE status = 'suspended') AS suspended,
			COUNT(*) FILTER (WHERE status = 'expired') AS expired
		FROM client_licenses
		WHERE deleted_at IS NULL`

	// Semua query bersifat non-fatal — log error tapi tetap lanjut dengan nilai 0.
	if err := h.db.QueryRowContext(ctx, licenseCountQ).Scan(
		&stats.TotalLicenses,
		&stats.ActiveLicenses,
		&stats.PendingLicenses,
		&stats.SuspendedLicenses,
		&stats.ExpiredLicenses,
	); err != nil {
		h.logger.Warn("DashboardHandler.GetStats: license counts", zap.Error(err))
	}

	// Total companies
	const companyCountQ = `SELECT COUNT(*) FROM companies WHERE deleted_at IS NULL`
	if err := h.db.QueryRowContext(ctx, companyCountQ).Scan(&stats.TotalCompanies); err != nil {
		h.logger.Warn("DashboardHandler.GetStats: company count", zap.Error(err))
	}

	// Proposal counts
	const proposalCountQ = `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'submitted') AS pending
		FROM proposals
		WHERE deleted_at IS NULL`
	if err := h.db.QueryRowContext(ctx, proposalCountQ).Scan(&stats.TotalProposals, &stats.PendingProposals); err != nil {
		h.logger.Warn("DashboardHandler.GetStats: proposal counts", zap.Error(err))
	}

	// Total revenue dari active licenses dengan contract_amount
	const revenueQ = `
		SELECT COALESCE(SUM(contract_amount), 0)
		FROM client_licenses
		WHERE status = 'active' AND contract_amount IS NOT NULL AND deleted_at IS NULL`
	if err := h.db.QueryRowContext(ctx, revenueQ).Scan(&stats.TotalRevenue); err != nil {
		h.logger.Warn("DashboardHandler.GetStats: revenue", zap.Error(err))
	}

	// Expiring licenses dalam 30 hari — non-fatal, skip jika query gagal.
	const expiringQ = `
		SELECT
			cl.id::text,
			cl.license_key,
			c.name AS company,
			cl.expires_at,
			EXTRACT(DAY FROM cl.expires_at - NOW())::int AS days_left
		FROM client_licenses cl
		JOIN companies c ON c.id = cl.company_id
		WHERE cl.status = 'active'
			AND cl.expires_at IS NOT NULL
			AND cl.expires_at > NOW()
			AND cl.expires_at <= NOW() + INTERVAL '30 days'
			AND cl.deleted_at IS NULL
		ORDER BY cl.expires_at ASC
		LIMIT 10`

	if rows, qErr := h.db.QueryxContext(ctx, expiringQ); qErr != nil {
		h.logger.Warn("DashboardHandler.GetStats: expiring query skipped", zap.Error(qErr))
	} else {
		defer rows.Close()
		for rows.Next() {
			var row expiringLicenseRow
			if err := rows.StructScan(&row); err != nil {
				h.logger.Error("DashboardHandler.GetStats: scan expiring row", zap.Error(err))
				continue
			}
			stats.ExpiringLicenses = append(stats.ExpiringLicenses, ExpiringLicense{
				ID:         row.ID,
				LicenseKey: row.LicenseKey,
				Company:    row.Company,
				ExpiresAt:  row.ExpiresAt.Format("2006-01-02"),
				DaysLeft:   row.DaysLeft,
			})
		}
	}

	// Recent activity — 10 audit logs terbaru, non-fatal.
	const activityQ = `
		SELECT entity_type, action, actor_name, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT 10`

	if actRows, qErr := h.db.QueryxContext(ctx, activityQ); qErr != nil {
		h.logger.Warn("DashboardHandler.GetStats: activity query skipped", zap.Error(qErr))
	} else {
		defer actRows.Close()
		for actRows.Next() {
			var row activityItemRow
			if err := actRows.StructScan(&row); err != nil {
				h.logger.Error("DashboardHandler.GetStats: scan activity row", zap.Error(err))
				continue
			}
			stats.RecentActivity = append(stats.RecentActivity, ActivityItem{
				EntityType: row.EntityType,
				Action:     row.Action,
				ActorName:  row.ActorName,
				CreatedAt:  row.CreatedAt.Format("2006-01-02T15:04:05Z"),
			})
		}
	}

	writeJSON(w, http.StatusOK, stats)
}
