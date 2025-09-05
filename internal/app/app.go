package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ekbaya/asham/internal/wire"
	"github.com/ekbaya/asham/pkg/config"
	"github.com/ekbaya/asham/pkg/db"
	"github.com/ekbaya/asham/pkg/db/migrations"
	"github.com/ekbaya/asham/pkg/db/redis"
	"gorm.io/gorm"
)

type App struct {
	addr   string
	db     *gorm.DB
	server *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize redis client with environment variables
	redisAddr := fmt.Sprintf("%s:%s", 
		os.Getenv("REDIS_HOST"), 
		os.Getenv("REDIS_PORT"))
	if os.Getenv("REDIS_HOST") == "" {
		redisAddr = "localhost:6379" // fallback for development
	}
	redis.InitRedis(redisAddr, "", 0)

	// Initialize database connection
	db, err := db.NewPostgresConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed database
	migrations.SeedDatabase(db)

	// Initialize services with db connection
	services, err := wire.InitializeServices(db)
	if err != nil {
		panic(err)
	}

	// Initialize router with dependencies
	router, err := InitRoutes(services)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize routes: %w", err)
	}

	if cfg.SEED_PERMISSIONS {
		// Auto-register permission mappings
		err = registerPermissionsAndResources(router, services.RbacService, services.PermissionResourceService)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-register permission mappings: %w", err)
		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	return &App{
		addr:   fmt.Sprintf(":%s", cfg.Server.Port),
		db:     db,
		server: server,
	}, nil
}

// Run starts the HTTP server
func (a *App) Run() error {
	fmt.Printf("Starting server on %s...\n", a.addr)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down server...")
	return a.server.Shutdown(ctx)
}
