//go:build !wasm

// Package main adalah entrypoint untuk Vernon License API server.
// Server ini menyediakan 2 public endpoints: POST /api/v1/register dan GET /api/v1/validate,
// serta internal API untuk Vernon App (WASM) dan serving WASM app itu sendiri.
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/infrastructure/database"
	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/flashlab/vernon-license/internal/handler"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/publicapi"
	"github.com/flashlab/vernon-license/internal/service"
	ratelimit "github.com/flashlab/vernon-license/pkg/middleware"
	"github.com/flashlab/vernon-license/pkg/scheduler"
	"time"
)

func main() {
	// Load .env jika ada — ignored jika belum ada (setup wizard yang akan membuatnya).
	_ = godotenv.Load()

	// Jika .env belum ada, jalankan setup wizard dan block di sana.
	// Wizard akan buat DB, migrate, buat superuser, tulis .env, lalu restart proses ini.
	if isSetupRequired() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8081"
		}
		serveSetupWizard(port)
		return
	}

	fxApp := fx.New(
		// Config
		fx.Provide(provideConfig),

		// Logger
		fx.Provide(provideLogger),

		// Database
		fx.Provide(provideDatabase),

		// Repositories — bind concrete types ke interfaces
		fx.Provide(
			fx.Annotate(database.NewLicenseRepo, fx.As(new(domain.LicenseRepository))),
			fx.Annotate(database.NewProductRepo, fx.As(new(domain.ProductRepository))),
			fx.Annotate(database.NewAuditRepo, fx.As(new(domain.AuditLogRepository))),
			fx.Annotate(database.NewUserRepo, fx.As(new(domain.UserRepository))),
			fx.Annotate(database.NewCompanyRepo, fx.As(new(domain.CompanyRepository))),
			fx.Annotate(database.NewProjectRepo, fx.As(new(domain.ProjectRepository))),
			fx.Annotate(database.NewProposalRepo, fx.As(new(domain.ProposalRepository))),
			fx.Annotate(database.NewNotificationRepo, fx.As(new(domain.NotificationRepository))),
			fx.Annotate(database.NewOTPRepository, fx.As(new(domain.OTPRepository))),
		),

		// Services
		service.ServiceModule,

		// Public API handlers
		publicapi.Module,

		// Internal API handlers
		handler.Module,

		// HTTP server
		fx.Provide(provideRouter),
		fx.Invoke(startServer),

		// Background scheduler
		fx.Invoke(startScheduler),
	)

	fxApp.Run()
}

// provideConfig membaca konfigurasi dari environment variables.
func provideConfig() (*config.Config, error) {
	return config.Load()
}

// provideLogger membuat Uber Zap logger.
func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	var log *zap.Logger
	var err error

	if cfg.LogLevel == "debug" {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		return nil, fmt.Errorf("provideLogger: %w", err)
	}
	return log, nil
}

// provideDatabase membuka koneksi ke PostgreSQL.
func provideDatabase(cfg *config.Config, lc fx.Lifecycle) (*sqlx.DB, error) {
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("provideDatabase: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

// provideRouter membuat Chi router dengan semua routes terdaftar:
// public API, internal API, dan Vernon App WASM handler.
func provideRouter(
	registerHandler *publicapi.RegisterHandler,
	validateHandler *publicapi.ValidateHandler,
	authHandler *handler.AuthHandler,
	setupHandler *handler.SetupHandler,
	companyHandler *handler.CompanyHandler,
	projectHandler *handler.ProjectHandler,
	licenseHandler *handler.LicenseHandler,
	proposalHandler *handler.ProposalHandler,
	productHandler *handler.ProductHandler,
	userHandler *handler.UserHandler,
	notifHandler *handler.NotificationHandler,
	dashboardHandler *handler.DashboardHandler,
	cfg *config.Config,
	log *zap.Logger,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public API routes dengan rate limiting 60 req/min per IP
	r.Group(func(r chi.Router) {
		r.Use(ratelimit.NewRateLimiter(60))
		r.Post("/api/v1/register", registerHandler.Handle)
		r.Get("/api/v1/validate", validateHandler.Handle)
	})

	// Internal API — no auth required
	r.Post("/api/internal/setup/install", setupHandler.Install)
	r.Get("/api/internal/setup/status", setupHandler.GetStatus)
	r.Post("/api/internal/auth/login", authHandler.Login)

	// Internal API — JWT auth required
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/api/internal/auth/me", authHandler.GetMe)

		// Companies
		r.Get("/api/internal/companies", companyHandler.List)
		r.Post("/api/internal/companies", companyHandler.Create)
		r.Get("/api/internal/companies/{id}", companyHandler.GetByID)
		r.Put("/api/internal/companies/{id}", companyHandler.Update)
		r.Delete("/api/internal/companies/{id}", companyHandler.Delete)

		// Projects (nested under company + standalone)
		r.Get("/api/internal/companies/{companyID}/projects", projectHandler.ListByCompany)
		r.Get("/api/internal/companies/{companyID}/licenses", licenseHandler.ListByCompany)
		r.Post("/api/internal/companies/{companyID}/projects", projectHandler.Create)
		r.Get("/api/internal/projects/{id}", projectHandler.GetByID)
		r.Put("/api/internal/projects/{id}", projectHandler.Update)
		r.Delete("/api/internal/projects/{id}", projectHandler.Delete)

		// Licenses
		r.Get("/api/internal/licenses", licenseHandler.List)
		r.Post("/api/internal/licenses", licenseHandler.Create)
		r.Get("/api/internal/licenses/{id}", licenseHandler.GetByID)
		r.Get("/api/internal/licenses/{id}/otp", licenseHandler.GetOTP)
		r.Put("/api/internal/licenses/{id}/activate", licenseHandler.Activate)
		r.Put("/api/internal/licenses/{id}/suspend", licenseHandler.Suspend)
		r.Put("/api/internal/licenses/{id}/renew", licenseHandler.Renew)
		r.Put("/api/internal/licenses/{id}/constraints", licenseHandler.UpdateConstraints)
		r.Put("/api/internal/licenses/{id}/status", licenseHandler.SetStatus)
		// Audit logs hanya bisa diakses oleh superuser
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("superuser"))
			r.Get("/api/internal/licenses/{id}/audit", licenseHandler.GetAuditLogs)
		})
		r.Get("/api/internal/projects/{projectID}/licenses", licenseHandler.ListByProject)

		// Proposals
		r.Get("/api/internal/proposals", proposalHandler.List)
		r.Post("/api/internal/proposals", proposalHandler.Create)
		r.Get("/api/internal/proposals/{id}", proposalHandler.GetByID)
		r.Put("/api/internal/proposals/{id}", proposalHandler.Update)
		r.Put("/api/internal/proposals/{id}/submit", proposalHandler.Submit)
		r.Put("/api/internal/proposals/{id}/approve", proposalHandler.Approve)
		r.Put("/api/internal/proposals/{id}/reject", proposalHandler.Reject)
		r.Get("/api/internal/proposals/{id}/pdf", proposalHandler.GetPDF)
		r.Get("/api/internal/projects/{projectID}/proposals", proposalHandler.ListByProject)

		// Dashboard
		r.Get("/api/internal/dashboard", dashboardHandler.GetStats)
		r.Get("/api/internal/dashboard/otp", dashboardHandler.GetOTP)

		// Products
		r.Get("/api/internal/products", productHandler.List)
		r.Post("/api/internal/products", productHandler.Create)
		r.Get("/api/internal/products/{id}", productHandler.GetByID)
		r.Put("/api/internal/products/{id}", productHandler.Update)
		r.Delete("/api/internal/products/{id}", productHandler.Delete)

		// Users
		r.Get("/api/internal/users", userHandler.List)
		r.Post("/api/internal/users", userHandler.Create)
		r.Put("/api/internal/users/{id}/active", userHandler.SetActive)

		// Notifications
		r.Get("/api/internal/notifications", notifHandler.List)
		r.Get("/api/internal/notifications/unread-count", notifHandler.CountUnread)
		r.Put("/api/internal/notifications/{id}/read", notifHandler.MarkRead)
		r.Put("/api/internal/notifications/read-all", notifHandler.MarkAllRead)
	})

	// Serve static web assets (CSS, icons, manifest)
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	r.Handle("/manifest.json", http.FileServer(http.Dir("web")))

	// Vernon App WASM — serve untuk semua non-API routes
	wasmHandler := &app.Handler{
		Name:            "Vernon License",
		Description:     "License Management System",
		Author:          "FlashLab",
		ThemeColor:      "#4D2975",
		BackgroundColor: "#0F0A1A",
		LoadingLabel:    "Loading Vernon...",
		Styles:          []string{"/web/app.css"},
	}
	// go-app's ServeHTTP returns 404 jika path tidak cocok dengan route yang terdaftar.
	// Untuk path frontend (bukan go-app internal assets), kita rewrite URL ke "/" agar
	// go-app selalu menyajikan WASM shell. Client-side routing (go-app) kemudian membaca
	// window.location.pathname untuk navigasi yang benar.
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		goAppInternal := map[string]bool{
			"/wasm_exec.js": true, "/app.js": true, "/app.js.gz": true,
			"/sw.js": true, "/manifest.webmanifest": true, "/robots.txt": true,
			"/app-worker.js": true,
		}
		if !goAppInternal[req.URL.Path] {
			req.URL.Path = "/"
		}
		wasmHandler.ServeHTTP(w, req)
	}))

	log.Info("Router configured",
		zap.String("public_api", "/api/v1/{register,validate}"),
		zap.String("internal_api", "/api/internal/*"),
		zap.String("app", "/* (WASM)"),
	)
	return r
}

// startServer mendaftarkan lifecycle hook untuk start/stop HTTP server.
func startServer(lc fx.Lifecycle, r *chi.Mux, cfg *config.Config, log *zap.Logger) {
	server := &http.Server{
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", ":"+cfg.Port)
			if err != nil {
				return fmt.Errorf("startServer: listen: %w", err)
			}
			log.Info("Starting Vernon License API", zap.String("port", cfg.Port))
			fmt.Printf("\n  Vernon License berjalan → http://localhost:%s\n\n", cfg.Port)
			go func() {
				if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
					log.Error("server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down server")
			return server.Shutdown(ctx)
		},
	})
}

// startScheduler mendaftarkan lifecycle hook untuk start/stop background scheduler.
func startScheduler(lc fx.Lifecycle, otpService *service.OTPService, log *zap.Logger) {
	sched := scheduler.New(log)

	// Cleanup expired OTP records setiap jam
	sched.Schedule("cleanup-expired-otp", 1*time.Hour, func(ctx context.Context) error {
		return otpService.CleanupExpired(ctx)
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go sched.Start(ctx)
			log.Info("Background scheduler started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sched.Stop()
			log.Info("Background scheduler stopped")
			return nil
		},
	})
}
