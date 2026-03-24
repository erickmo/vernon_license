package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	createclientlicense "github.com/flashlab/flasherp-developer-api/internal/command/create_client_license"
	createproduct "github.com/flashlab/flasherp-developer-api/internal/command/create_product"
	deleteproduct "github.com/flashlab/flasherp-developer-api/internal/command/delete_product"
	loginhandler "github.com/flashlab/flasherp-developer-api/internal/command/login"
	markallread "github.com/flashlab/flasherp-developer-api/internal/command/mark_all_notifications_read"
	markread "github.com/flashlab/flasherp-developer-api/internal/command/mark_notification_read"
	provisionlicense "github.com/flashlab/flasherp-developer-api/internal/command/provision_license"
	registerdevice "github.com/flashlab/flasherp-developer-api/internal/command/register_device"
	setupinstall "github.com/flashlab/flasherp-developer-api/internal/command/setup_install"
	unregisterdevice "github.com/flashlab/flasherp-developer-api/internal/command/unregister_device"
	updatelicenseconstraints "github.com/flashlab/flasherp-developer-api/internal/command/update_license_constraints"
	updatelicensestatus "github.com/flashlab/flasherp-developer-api/internal/command/update_license_status"
	updateproduct "github.com/flashlab/flasherp-developer-api/internal/command/update_product"
	httphandler "github.com/flashlab/flasherp-developer-api/internal/delivery/http"
	getclientlicense "github.com/flashlab/flasherp-developer-api/internal/query/get_client_license"
	getdashboard "github.com/flashlab/flasherp-developer-api/internal/query/get_dashboard"
	getme "github.com/flashlab/flasherp-developer-api/internal/query/get_me"
	getproduct "github.com/flashlab/flasherp-developer-api/internal/query/get_product"
	getsetupstatus "github.com/flashlab/flasherp-developer-api/internal/query/get_setup_status"
	getunreadcount "github.com/flashlab/flasherp-developer-api/internal/query/get_unread_count"
	listauditlogs "github.com/flashlab/flasherp-developer-api/internal/query/list_audit_logs"
	listclientlicenses "github.com/flashlab/flasherp-developer-api/internal/query/list_client_licenses"
	listnotifications "github.com/flashlab/flasherp-developer-api/internal/query/list_notifications"
	listproducts "github.com/flashlab/flasherp-developer-api/internal/query/list_products"
	"github.com/flashlab/flasherp-developer-api/infrastructure/config"
	"github.com/flashlab/flasherp-developer-api/infrastructure/database"
	"github.com/flashlab/flasherp-developer-api/internal/service"
	"github.com/flashlab/flasherp-developer-api/pkg/commandbus"
	jwtpkg "github.com/flashlab/flasherp-developer-api/pkg/jwt"
	appmiddleware "github.com/flashlab/flasherp-developer-api/pkg/middleware"
	"github.com/flashlab/flasherp-developer-api/pkg/querybus"
)

//go:embed web
var webFiles embed.FS

type Handlers struct {
	fx.In

	Auth         *httphandler.AuthHandler
	License      *httphandler.LicenseHandler
	Dashboard    *httphandler.DashboardHandler
	Setup        *httphandler.SetupHandler
	Product      *httphandler.ProductHandler
	Audit        *httphandler.AuditHandler
	Notification *httphandler.NotificationHandler
}

type CommandHandlers struct {
	fx.In

	UpdateStatus      *updatelicensestatus.Handler
	UpdateConstraints *updatelicenseconstraints.Handler
	CreateProduct     *createproduct.Handler
	UpdateProduct     *updateproduct.Handler
	DeleteProduct     *deleteproduct.Handler
	MarkRead          *markread.Handler
	MarkAllRead       *markallread.Handler
	RegisterDevice    *registerdevice.Handler
	UnregisterDevice  *unregisterdevice.Handler
}

type QueryHandlers struct {
	fx.In

	ListLicenses      *listclientlicenses.Handler
	GetLicense        *getclientlicense.Handler
	GetMe             *getme.Handler
	GetSetupStatus    *getsetupstatus.Handler
	ListProducts      *listproducts.Handler
	GetProduct        *getproduct.Handler
	ListAuditLogs     *listauditlogs.Handler
	ListNotifications *listnotifications.Handler
	GetUnreadCount    *getunreadcount.Handler
	GetDashboard      *getdashboard.Handler
}

func registerHandlers(
	cmdBus *commandbus.CommandBus,
	qBus *querybus.QueryBus,
	cmdHandlers CommandHandlers,
	qHandlers QueryHandlers,
) {
	// Command handlers
	cmdBus.Register("update_license_status.UpdateLicenseStatus", cmdHandlers.UpdateStatus)
	cmdBus.Register("update_license_constraints.UpdateLicenseConstraints", cmdHandlers.UpdateConstraints)
	cmdBus.Register("create_product.CreateProduct", cmdHandlers.CreateProduct)
	cmdBus.Register("update_product.UpdateProduct", cmdHandlers.UpdateProduct)
	cmdBus.Register("delete_product.DeleteProduct", cmdHandlers.DeleteProduct)
	cmdBus.Register("mark_notification_read.MarkNotificationRead", cmdHandlers.MarkRead)
	cmdBus.Register("mark_all_notifications_read.MarkAllNotificationsRead", cmdHandlers.MarkAllRead)
	cmdBus.Register("register_device.RegisterDevice", cmdHandlers.RegisterDevice)
	cmdBus.Register("unregister_device.UnregisterDevice", cmdHandlers.UnregisterDevice)

	// Query handlers
	qBus.Register("list_client_licenses.ListClientLicenses", qHandlers.ListLicenses)
	qBus.Register("get_client_license.GetClientLicense", qHandlers.GetLicense)
	qBus.Register("get_me.GetMe", qHandlers.GetMe)
	qBus.Register("get_setup_status.GetSetupStatus", qHandlers.GetSetupStatus)
	qBus.Register("list_products.ListProducts", qHandlers.ListProducts)
	qBus.Register("get_product.GetProduct", qHandlers.GetProduct)
	qBus.Register("list_audit_logs.ListAuditLogs", qHandlers.ListAuditLogs)
	qBus.Register("list_notifications.ListNotifications", qHandlers.ListNotifications)
	qBus.Register("get_unread_count.GetUnreadCount", qHandlers.GetUnreadCount)
	qBus.Register("get_dashboard.GetDashboard", qHandlers.GetDashboard)
}

func startServer(lc fx.Lifecycle, cfg *config.Config, jwtService *jwtpkg.Service, handlers Handlers) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(appmiddleware.CORS())
	r.Use(appmiddleware.Logger(logger))
	r.Use(appmiddleware.Recoverer())

	// Public routes
	r.Get("/api/v1/setup/status", handlers.Setup.Status)
	r.Post("/api/v1/setup/install", handlers.Setup.Install)
	r.Post("/api/v1/auth/login", handlers.Auth.Login)
	r.Get("/api/v1/client/license", handlers.License.GetClientLicense)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.RequireAuth(jwtService))

		r.Get("/api/v1/auth/me", handlers.Auth.Me)

		// Dashboard
		r.Get("/api/v1/dashboard", handlers.Dashboard.Dashboard)

		// License routes — sales dan superuser
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("sales", "developer_sales", "project_owner", "superuser"))
			r.Get("/api/v1/licenses", handlers.License.List)
			r.Post("/api/v1/licenses", handlers.License.Create)
			r.Get("/api/v1/licenses/{id}", handlers.License.Get)
			r.Put("/api/v1/licenses/{id}/status", handlers.License.UpdateStatus)
			r.Put("/api/v1/licenses/{id}/constraints", handlers.License.UpdateConstraints)
			r.Post("/api/v1/licenses/{id}/provision", handlers.License.Provision)
			r.Get("/api/v1/licenses/{id}/audit", handlers.Audit.GetLicenseAudit)
		})

		// Product routes — read: semua role, write: project_owner/superuser
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("sales", "developer_sales", "project_owner", "superuser"))
			r.Get("/api/v1/products", handlers.Product.List)
			r.Get("/api/v1/products/{id}", handlers.Product.Get)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("project_owner", "superuser"))
			r.Post("/api/v1/products", handlers.Product.Create)
			r.Put("/api/v1/products/{id}", handlers.Product.Update)
			r.Delete("/api/v1/products/{id}", handlers.Product.Delete)
		})

		// Global audit — superuser only
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("superuser"))
			r.Get("/api/v1/audit", handlers.Audit.GetAllAudit)
		})

		// Notification routes — semua authenticated user
		r.Get("/api/v1/notifications", handlers.Notification.ListNotifications)
		r.Put("/api/v1/notifications/{id}/read", handlers.Notification.MarkRead)
		r.Put("/api/v1/notifications/read-all", handlers.Notification.MarkAllRead)
		r.Get("/api/v1/notifications/unread-count", handlers.Notification.GetUnreadCount)
		r.Post("/api/v1/devices", handlers.Notification.RegisterDevice)
		r.Delete("/api/v1/devices/{token}", handlers.Notification.UnregisterDevice)
	})

	// Serve PWA static files — must be after API routes
	webFS, err := fs.Sub(webFiles, "web")
	if err == nil {
		fileServer := http.FileServer(http.FS(webFS))
		// Chi requires explicit GET registration for FileServer to work correctly
		r.Get("/", fileServer.ServeHTTP)
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
			http.StripPrefix(pathPrefix, fileServer).ServeHTTP(w, r)
		})
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info().Str("port", cfg.HTTPPort).Msg("server dimulai")
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal().Err(err).Msg("server error")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("server dihentikan")
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}

func main() {
	// Pre-flight: jika .env tidak ada atau DB tidak bisa dikoneksi, jalankan setup wizard
	if needsSetup() {
		runSetupServer("8081")
		return
	}

	// Handle sinyal OS untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := fx.New(
		fx.Provide(
			config.Load,
			// DB
			func(cfg *config.Config) (*sqlx.DB, error) {
				return database.NewPostgres(cfg.DatabaseURL)
			},
			// JWT
			func(cfg *config.Config) *jwtpkg.Service {
				return jwtpkg.NewService(cfg.JWTSecret, cfg.JWTExpHours)
			},
			// Buses
			commandbus.New,
			querybus.New,
			// Repositories
			func(db *sqlx.DB) *database.UserRepository {
				return database.NewUserRepository(db)
			},
			func(db *sqlx.DB) *database.ClientLicenseRepository {
				return database.NewClientLicenseRepository(db)
			},
			func(db *sqlx.DB) *database.ProductRepository {
				return database.NewProductRepository(db)
			},
			func(db *sqlx.DB) *database.AuditRepository {
				return database.NewAuditRepository(db)
			},
			func(db *sqlx.DB) *database.NotificationRepository {
				return database.NewNotificationRepository(db)
			},
			// Command handlers — license
			func(repo *database.UserRepository, jwt *jwtpkg.Service) *loginhandler.Handler {
				return loginhandler.NewHandler(repo, jwt)
			},
			func(repo *database.ClientLicenseRepository, cfg *config.Config) *createclientlicense.Handler {
				return createclientlicense.NewHandler(repo, cfg)
			},
			func(repo *database.ClientLicenseRepository, auditRepo *database.AuditRepository) *updatelicensestatus.Handler {
				return updatelicensestatus.NewHandler(repo, auditRepo)
			},
			func(repo *database.ClientLicenseRepository, auditRepo *database.AuditRepository) *updatelicenseconstraints.Handler {
				return updatelicenseconstraints.NewHandler(repo, auditRepo)
			},
			func(repo *database.ClientLicenseRepository, auditRepo *database.AuditRepository) *provisionlicense.Handler {
				return provisionlicense.NewHandler(repo, auditRepo)
			},
			// Command handlers — product
			func(repo *database.ProductRepository, auditRepo *database.AuditRepository) *createproduct.Handler {
				return createproduct.NewHandler(repo, auditRepo)
			},
			func(readRepo *database.ProductRepository, writeRepo *database.ProductRepository, auditRepo *database.AuditRepository) *updateproduct.Handler {
				return updateproduct.NewHandler(readRepo, writeRepo, auditRepo)
			},
			func(repo *database.ProductRepository) *deleteproduct.Handler {
				return deleteproduct.NewHandler(repo)
			},
			// Command handlers — notification
			func(repo *database.NotificationRepository) *markread.Handler {
				return markread.NewHandler(repo)
			},
			func(repo *database.NotificationRepository) *markallread.Handler {
				return markallread.NewHandler(repo)
			},
			func(repo *database.NotificationRepository) *registerdevice.Handler {
				return registerdevice.NewHandler(repo)
			},
			func(repo *database.NotificationRepository) *unregisterdevice.Handler {
				return unregisterdevice.NewHandler(repo)
			},
			// Query handlers
			func(repo *database.ClientLicenseRepository) *listclientlicenses.Handler {
				return listclientlicenses.NewHandler(repo)
			},
			func(repo *database.ClientLicenseRepository) *getclientlicense.Handler {
				return getclientlicense.NewHandler(repo)
			},
			func(repo *database.UserRepository) *getme.Handler {
				return getme.NewHandler(repo)
			},
			func(repo *database.UserRepository) *getsetupstatus.Handler {
				return getsetupstatus.NewHandler(repo)
			},
			func(writeRepo *database.UserRepository, readRepo *database.UserRepository) *setupinstall.Handler {
				return setupinstall.NewHandler(writeRepo, readRepo)
			},
			// Query handlers — product
			func(repo *database.ProductRepository) *listproducts.Handler {
				return listproducts.NewHandler(repo)
			},
			func(repo *database.ProductRepository) *getproduct.Handler {
				return getproduct.NewHandler(repo)
			},
			// Query handlers — audit
			func(repo *database.AuditRepository) *listauditlogs.Handler {
				return listauditlogs.NewHandler(repo)
			},
			// Query handlers — notification
			func(repo *database.NotificationRepository) *listnotifications.Handler {
				return listnotifications.NewHandler(repo)
			},
			func(repo *database.NotificationRepository) *getunreadcount.Handler {
				return getunreadcount.NewHandler(repo)
			},
			// Query handlers — dashboard
			func(db *sqlx.DB) *getdashboard.Handler {
				return getdashboard.NewHandler(db)
			},
			// Services
			func(db *sqlx.DB, repo *database.NotificationRepository) *service.NotificationService {
				return service.NewNotificationService(db, repo)
			},
			// HTTP handlers
			httphandler.NewAuthHandler,
			httphandler.NewLicenseHandler,
			httphandler.NewDashboardHandler,
			httphandler.NewSetupHandler,
			httphandler.NewProductHandler,
			httphandler.NewAuditHandler,
			httphandler.NewNotificationHandler,
		),
		fx.Invoke(
			registerHandlers,
			startServer,
			func(lc fx.Lifecycle, svc *service.NotificationService) {
				ctx, cancel := context.WithCancel(context.Background())
				lc.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						svc.StartExpiryChecker(ctx)
						return nil
					},
					OnStop: func(_ context.Context) error {
						cancel()
						return nil
					},
				})
			},
		),
	)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := app.Stop(ctx); err != nil {
			log.Fatal().Err(err).Msg("gagal stop app")
		}
	}()

	app.Run()
}
