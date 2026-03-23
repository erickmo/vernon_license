package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	createclientlicense "github.com/flashlab/flasherp-developer-api/internal/command/create_client_license"
	loginhandler "github.com/flashlab/flasherp-developer-api/internal/command/login"
	provisionlicense "github.com/flashlab/flasherp-developer-api/internal/command/provision_license"
	setupinstall "github.com/flashlab/flasherp-developer-api/internal/command/setup_install"
	updatelicenseconstraints "github.com/flashlab/flasherp-developer-api/internal/command/update_license_constraints"
	updatelicensestatus "github.com/flashlab/flasherp-developer-api/internal/command/update_license_status"
	httphandler "github.com/flashlab/flasherp-developer-api/internal/delivery/http"
	getclientlicense "github.com/flashlab/flasherp-developer-api/internal/query/get_client_license"
	getme "github.com/flashlab/flasherp-developer-api/internal/query/get_me"
	getsetupstatus "github.com/flashlab/flasherp-developer-api/internal/query/get_setup_status"
	listclientlicenses "github.com/flashlab/flasherp-developer-api/internal/query/list_client_licenses"
	"github.com/flashlab/flasherp-developer-api/infrastructure/config"
	"github.com/flashlab/flasherp-developer-api/infrastructure/database"
	"github.com/flashlab/flasherp-developer-api/pkg/commandbus"
	jwtpkg "github.com/flashlab/flasherp-developer-api/pkg/jwt"
	appmiddleware "github.com/flashlab/flasherp-developer-api/pkg/middleware"
	"github.com/flashlab/flasherp-developer-api/pkg/querybus"
)

type Handlers struct {
	fx.In

	Auth      *httphandler.AuthHandler
	License   *httphandler.LicenseHandler
	Dashboard *httphandler.DashboardHandler
	Setup     *httphandler.SetupHandler
}

type CommandHandlers struct {
	fx.In

	UpdateStatus      *updatelicensestatus.Handler
	UpdateConstraints *updatelicenseconstraints.Handler
}

type QueryHandlers struct {
	fx.In

	ListLicenses    *listclientlicenses.Handler
	GetLicense      *getclientlicense.Handler
	GetMe           *getme.Handler
	GetSetupStatus  *getsetupstatus.Handler
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

	// Query handlers
	qBus.Register("list_client_licenses.ListClientLicenses", qHandlers.ListLicenses)
	qBus.Register("get_client_license.GetClientLicense", qHandlers.GetLicense)
	qBus.Register("get_me.GetMe", qHandlers.GetMe)
	qBus.Register("get_setup_status.GetSetupStatus", qHandlers.GetSetupStatus)
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

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.RequireAuth(jwtService))

		r.Get("/api/v1/auth/me", handlers.Auth.Me)

		// Dashboard
		r.Get("/api/v1/dashboard", handlers.Dashboard.Dashboard)

		// License routes — hanya developer_sales dan superuser
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireRole("developer_sales", "superuser"))
			r.Get("/api/v1/licenses", handlers.License.List)
			r.Post("/api/v1/licenses", handlers.License.Create)
			r.Get("/api/v1/licenses/{id}", handlers.License.Get)
			r.Put("/api/v1/licenses/{id}/status", handlers.License.UpdateStatus)
			r.Put("/api/v1/licenses/{id}/constraints", handlers.License.UpdateConstraints)
			r.Post("/api/v1/licenses/{id}/provision", handlers.License.Provision)
		})
	})

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
			// Repositories — provide sebagai interface dan concrete type
			func(db *sqlx.DB) *database.UserRepository {
				return database.NewUserRepository(db)
			},
			func(db *sqlx.DB) *database.ClientLicenseRepository {
				return database.NewClientLicenseRepository(db)
			},
			// Command handlers
			func(repo *database.UserRepository, jwt *jwtpkg.Service) *loginhandler.Handler {
				return loginhandler.NewHandler(repo, jwt)
			},
			func(repo *database.ClientLicenseRepository) *createclientlicense.Handler {
				return createclientlicense.NewHandler(repo)
			},
			func(repo *database.ClientLicenseRepository) *updatelicensestatus.Handler {
				return updatelicensestatus.NewHandler(repo)
			},
			func(repo *database.ClientLicenseRepository) *updatelicenseconstraints.Handler {
				return updatelicenseconstraints.NewHandler(repo)
			},
			func(repo *database.ClientLicenseRepository) *provisionlicense.Handler {
				return provisionlicense.NewHandler(repo)
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
			// HTTP handlers
			httphandler.NewAuthHandler,
			httphandler.NewLicenseHandler,
			httphandler.NewDashboardHandler,
			httphandler.NewSetupHandler,
		),
		fx.Invoke(
			registerHandlers,
			startServer,
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
