package http

import (
	"context"
	"net/http"
	"os"
	"study-golang-backend/internal/delivery/http/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// ProvideJWTSecret loads the JWT secret from environment variables
func ProvideJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		slog.Error("JWT_SECRET not found in environment")
		os.Exit(1)
	}
	return []byte(secret)
}

// NewGinEngine initializes Gin with middlewares
func NewGinEngine() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware(), gin.Recovery())
	return r
}

// RunServer registers routes and starts the Gin server within Fx lifecycle
func RunServer(
	lc fx.Lifecycle,
	r *gin.Engine,
	userHandler *UserHandler,
	productHandler *ProductHandler,
	cartHandler *CartHandler,
) {
	// Register user public routes
	userHandler.RegisterRouter(r.Group(""))

	// Register product protected routes
	jwtSecret := ProvideJWTSecret()
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	productHandler.RegisterRouter(protected)

	// Register cart protected routes
	cartProtected := r.Group("")
	cartProtected.Use(middleware.AuthMiddleware(jwtSecret))
	cartHandler.RegisterRouter(cartProtected)

	// Define HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Append Fx start/stop lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Starting HTTP server on :8080")
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("HTTP server failed", slog.String("error", err.Error()))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Stopping HTTP server gracefully...")
			return srv.Shutdown(ctx)
		},
	})
}

// Module exports the HTTP module providers and invokes the server runner
var Module = fx.Module("http",
	fx.Provide(
		ProvideJWTSecret,
		NewGinEngine,
		NewUserHandler,
		NewProductHandler,
		NewCartHandler,
	),
	fx.Invoke(RunServer),
)
