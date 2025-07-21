package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"item-pdp-service/internal/application/http/handlers"
	"item-pdp-service/internal/application/http/routes"
	"item-pdp-service/internal/application/usecase"
	"item-pdp-service/internal/infrastructure/config"
	"item-pdp-service/internal/infrastructure/database"
	"item-pdp-service/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Provide dependencies
		fx.Provide(
			func() (*config.Config, error) {
				return config.Load("./configs")
			},
			setupLogger,
			database.NewConnection,
			persistence.NewPostgresItemRepository,
			// Mock services for dependency injection (part of intentional flaws)
			func() usecase.InventoryService {
				return &mockInventoryService{}
			},
			func() usecase.CategoryService {
				return &mockCategoryService{}
			},
			func() usecase.PricingService {
				return &mockPricingService{}
			},
			usecase.NewItemUseCase,
			handlers.NewItemHandler,
			setupGinEngine,
			setupServer,
		),
		// Invoke the server
		fx.Invoke(runServer),
	).Run()
}

// setupLogger configures the logger
func setupLogger(cfg *config.Config) zerolog.Logger {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	if cfg.Log.Format == "pretty" || cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return log.Logger.With().
		Str("service", cfg.App.Name).
		Str("version", cfg.App.Version).
		Str("environment", cfg.App.Environment).
		Logger()
}

// setupGinEngine configures the Gin engine
func setupGinEngine(cfg *config.Config, itemHandler *handlers.ItemHandler) *gin.Engine {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	// Setup middlewares
	routes.SetupMiddlewares(router)

	// Setup routes
	routes.SetupRoutes(router, itemHandler)

	return router
}

// setupServer creates the HTTP server
func setupServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

// runServer starts the HTTP server with graceful shutdown
func runServer(lc fx.Lifecycle, cfg *config.Config, server *http.Server, db *database.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().
				Str("address", cfg.GetServerAddress()).
				Msg("Starting HTTP server")

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("Failed to start server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down server")

			// Create shutdown context with timeout
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Shutdown server gracefully
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Server forced to shutdown")
				return err
			}

			// Close database connection
			if err := db.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close database connection")
				return err
			}

			log.Info().Msg("Server stopped gracefully")
			return nil
		},
	})

	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info().Msg("Received shutdown signal")
		// The fx framework will handle the shutdown via OnStop hooks
	}()
}

// Mock service implementations for dependency injection (part of intentional flaws)
type mockInventoryService struct{}

func (s *mockInventoryService) ReserveInventory(ctx context.Context, itemID string, quantity int) error {
	return nil
}

func (s *mockInventoryService) ReleaseInventory(ctx context.Context, itemID string, quantity int) error {
	return nil
}

type mockCategoryService struct{}

func (s *mockCategoryService) ValidateCategory(ctx context.Context, category string) error {
	return nil
}

func (s *mockCategoryService) GetCategoryDiscounts(ctx context.Context, category string) (float64, error) {
	return 0.0, nil
}

type mockPricingService struct{}

func (s *mockPricingService) CalculatePrice(ctx context.Context, basePrice float64, category string) (float64, error) {
	return basePrice, nil
}

func (s *mockPricingService) ApplyDiscounts(ctx context.Context, price float64, itemID string) (float64, error) {
	return price, nil
}
