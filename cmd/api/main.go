package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/FelipeCJSEP/pulsewatch/internal/application"
	httpdelivery "github.com/FelipeCJSEP/pulsewatch/internal/delivery/http"
	"github.com/FelipeCJSEP/pulsewatch/internal/infrastructure/checker"
	"github.com/FelipeCJSEP/pulsewatch/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	dbURL := "postgres://pulsewatch:pulsewatch123@localhost:5432/pulsewatch_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, dbURL)

	if err != nil {
		slog.Error("Erro ao conectar com o banco", "erro", err)
		os.Exit(1)
	}

	defer pool.Close()

	// injeção de dependência
	monitorRepo := postgres.NewMonitorRepository(pool)
	incidentRepo := postgres.NewIncidentRepository(pool)
	monitorUseCase := application.NewMonitorUseCase(monitorRepo, incidentRepo)

	//inicia worker de Health Check em segundo plano
	worker := checker.NewWorker(monitorUseCase)
	go worker.Start(ctx)

	// handlers e rotas HTPP
	handler := httpdelivery.NewMonitorHandler(monitorUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/monitors", handler.CreateMonitor)
	mux.HandleFunc("GET /api/monitors", handler.ListMonitors)
	mux.HandleFunc("GET /api/monitors/{id}/incidents", handler.ListIncidents)
	mux.HandleFunc("PATCH /api/incidents/{id}/acknowledge", handler.AcknowledgeIncident)

	// sobe o servidor HTTP na porta 8080 com suporte a CORS
	slog.Info("Servidor PulseWatch API rodando em http://localhost:8080")

	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		slog.Error("Erro no servidor HTTP", "error", err)
	}
}
