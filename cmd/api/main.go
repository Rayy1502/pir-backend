package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Rayy1502/pir-backend/internal/config"
	"github.com/Rayy1502/pir-backend/internal/db"
	"github.com/Rayy1502/pir-backend/internal/handler"
	"github.com/Rayy1502/pir-backend/internal/middleware"
	"github.com/Rayy1502/pir-backend/internal/repository"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Connect ke SQL Server
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()
	log.Println("✅ Connected to SQL Server")

	// 3. Init repositories
	userRepo := repository.NewUserRepository(database)

	// 4. Init services
	authService := service.NewAuthService(cfg, userRepo)

	// 5. Init handlers
	authHandler := handler.NewAuthHandler(authService)

	// 6. Setup router  ← INI YANG HILANG DI FILE KAMU
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	// 7. CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 8. Health check
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		if err := database.PingContext(req.Context()); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "DB_ERROR", "Database unreachable", nil)
			return
		}
		response.Success(w, map[string]any{
			"status": "ok",
			"app":    "PIR",
			"time":   time.Now().Format(time.RFC3339),
		}, nil)
	})

	// 9. API v1
	r.Route("/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Post("/auth/login", authHandler.Login)

		// Auth routes (protected)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWTSecret))
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/logout", authHandler.Logout)
		})
	})

	// 10. Server setup
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 11. Start server
	go func() {
		log.Printf("🚀 PIR Backend running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 12. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("✅ Server exited gracefully")
}
