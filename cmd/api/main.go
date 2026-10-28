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
	pengaturanRepo := repository.NewPengaturanRepository(database)
	proyekRepo := repository.NewProyekRepository(database)
	pirRepo := repository.NewPIRRepository(database)
	surveyRepo := repository.NewSurveyRepository(database)
	laporanRepo := repository.NewLaporanRepository(database)
	tindakLanjutRepo := repository.NewTindakLanjutRepository(database)

	// 4. Init services
	authService := service.NewAuthService(cfg, userRepo)
	pengaturanService := service.NewPengaturanService(pengaturanRepo)
	proyekService := service.NewProyekService(proyekRepo)
	pirService := service.NewPIRService(pirRepo)
	surveyService := service.NewSurveyService(surveyRepo, pirRepo)
	laporanService := service.NewLaporanService(laporanRepo, proyekRepo)
	tindakLanjutService := service.NewTindakLanjutService(tindakLanjutRepo)

	// 5. Init handlers
	authHandler := handler.NewAuthHandler(authService)
	pengaturanHandler := handler.NewPengaturanHandler(pengaturanService)
	proyekHandler := handler.NewProyekHandler(proyekService)
	dashboardHandler := handler.NewDashboardHandler(proyekService)
	pirHandler := handler.NewPIRHandler(pirService, "./uploads")
	surveyHandler := handler.NewSurveyHandler(surveyService)
	laporanHandler := handler.NewLaporanHandler(laporanService, "./uploads")
	tindakLanjutHandler := handler.NewTindakLanjutHandler(tindakLanjutService, "./uploads")

	// 6. Setup router
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	// 7. CORS (must be before any routes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static file serving for uploads
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Fallback JSON for 404 & 405
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Endpoint tidak ditemukan: "+req.Method+" "+req.URL.Path, nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		response.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode tidak diizinkan: "+req.Method+" "+req.URL.Path, nil)
	})

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		response.Success(w, map[string]string{
			"message": "PIR Backend API is running",
		}, nil)
	})

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

		// Public Survey routes for respondents (no auth needed)
		r.Get("/public/survey/{id}", surveyHandler.GetPublic)
		r.Post("/public/survey/{id}/submit", surveyHandler.PublicSubmit)

		// Auth routes (protected)
		r.Group(func(r chi.Router) {

			//auth
			r.Use(middleware.JWTAuth(cfg.JWTSecret))
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/logout", authHandler.Logout)

			//pengaturan
			r.Get("/pengaturan", pengaturanHandler.GetAll)
			r.Put("/pengaturan/{unit}", pengaturanHandler.Update)

			// Proyek
			r.Get("/proyek", proyekHandler.GetAll)
			r.Post("/proyek", proyekHandler.Create)
			r.Get("/proyek/{id}", proyekHandler.GetByID)

			// Dashboard
			r.Get("/dashboard/stats", dashboardHandler.Stats)
			r.Get("/dashboard/statistik-pir", dashboardHandler.StatistikPIR)

			// PIR Form Step 1 & Dokumentasi
			r.Post("/pir", pirHandler.Create)
			r.Get("/pir", pirHandler.GetByProyekOrList)
			r.Get("/pir/{id}", pirHandler.GetByID)
			r.Put("/pir/{id}", pirHandler.Update)
			r.Put("/pir/{id}/approve", pirHandler.Approve)
			r.Get("/pir/{id}/dokumentasi", pirHandler.GetDokumentasi)
			r.Post("/pir/{id}/dokumentasi", pirHandler.UploadDokumentasi)
			r.Delete("/pir/dokumentasi/{dokId}", pirHandler.DeleteDokumentasi)

			// Survey Step 2 (PIC)
			r.Post("/pir/{id}/survey/questions", surveyHandler.SaveQuestions)
			r.Post("/pir/{id}/survey", surveyHandler.Save)
			r.Get("/pir/{id}/survey", surveyHandler.Get)

			// Laporan Step 3
			r.Post("/laporan", laporanHandler.Create)
			r.Get("/laporan", laporanHandler.GetAll)
			r.Get("/laporan/{id}", laporanHandler.GetByID)
			r.Get("/laporan/{id}/pdf", laporanHandler.ExportPDF)
			r.Get("/laporan/proyek/{proyekId}", laporanHandler.GetByProyekID)
			r.Post("/laporan/{id}/upload-ttd", laporanHandler.UploadSignedPDF)

			// Tindak Lanjut Step 4
			r.Get("/tindak-lanjut", tindakLanjutHandler.GetAll)
			r.Get("/tindak-lanjut/{id}", tindakLanjutHandler.GetByID)
			r.Post("/tindak-lanjut", tindakLanjutHandler.Create)
			r.Put("/tindak-lanjut/{id}", tindakLanjutHandler.Update)
			r.Post("/tindak-lanjut/{id}/bukti", tindakLanjutHandler.UploadBukti)
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
