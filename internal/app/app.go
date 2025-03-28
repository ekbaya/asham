package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ekbaya/asham/internal/wire"
	"github.com/ekbaya/asham/pkg/config"
	"github.com/ekbaya/asham/pkg/db"
	"github.com/ekbaya/asham/pkg/db/migrations"
	"gorm.io/gorm"
)

type App struct {
	addr   string
	db     *gorm.DB
	server *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
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
