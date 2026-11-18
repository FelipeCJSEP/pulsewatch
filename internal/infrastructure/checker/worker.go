package checker

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FelipeCJSEP/pulsewatch/internal/application"
	"github.com/FelipeCJSEP/pulsewatch/internal/domain"
)

type Worker struct {
	useCase *application.MonitorUseCase
	client  *http.Client
}

func NewWorker(useCase *application.MonitorUseCase) *Worker {
	return &Worker{
		useCase: useCase,

		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (w *Worker) checkOne(ctx context.Context, m *domain.Monitor) {
	resp, err := w.client.Get(m.TargetURL)

	if err != nil {
		// falhou por erro de rede, DNS ou Timeout
		_ = w.useCase.ProcessCheckResult(ctx, m.ID, false, err.Error())
		slog.Error("Erro ao verificar monitor", "error", err)
		return
	}

	defer resp.Body.Close()

	// Se o sattus HTTP for 2xx ou 3xx, é sucesso
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		_ = w.useCase.ProcessCheckResult(ctx, m.ID, true, "")
	} else {
		reason := fmt.Sprintf("HTTP %d", resp.StatusCode)
		_ = w.useCase.ProcessCheckResult(ctx, m.ID, false, reason)
	}

	slog.Info("Monitor verificado", "name", m.Name, "url", m.TargetURL)
}

func (w *Worker) checkAll(ctx context.Context) {
	monitors, err := w.useCase.ListMonitors(ctx)

	if err != nil {
		slog.Error("Erro ao listar monitores no worker", "error", err)
		return
	}

	for _, m := range monitors {
		if m.Status != domain.StatusPaused {
			go w.checkOne(ctx, m)
		}
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	slog.Info("Worker de Health Check iniciado (intervalo: 10s)")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker encerrado")
			return
		case <-ticker.C:
			w.checkAll(ctx)
		}
	}
}
