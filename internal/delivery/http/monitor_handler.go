package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/FelipeCJSEP/pulsewatch/internal/application"
)

type MonitorHandler struct {
	useCase *application.MonitorUseCase
}

type createMonitorRequest struct {
	Name             string `json:"name"`
	TargetURL        string `json:"targetUrl"`
	IntervalSeconds  int    `json:"intervalSeconds"`
	FailureThreshold int    `json:"failureThreshold"`
}

func NewMonitorHandler(useCase *application.MonitorUseCase) *MonitorHandler {
	return &MonitorHandler{useCase: useCase}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (mh *MonitorHandler) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	var req createMonitorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}

	monitor, err := mh.useCase.CreateMonitor(r.Context(), application.CreateMonitorInput{
		Name:             req.Name,
		TargetURL:        req.TargetURL,
		IntervalSeconds:  req.IntervalSeconds,
		FailureThreshold: req.FailureThreshold,
	})

	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, monitor)
}

func (mh *MonitorHandler) ListMonitors(w http.ResponseWriter, r *http.Request) {
	monitors, err := mh.useCase.ListMonitors(r.Context())

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, monitors)
}

func (mh *MonitorHandler) ListIncidents(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("id")
	incidents, err := mh.useCase.ListIncidentsByMonitor(r.Context(), monitorID)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"erro": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, incidents)
}

func (mh *MonitorHandler) AcknowledgeIncident(w http.ResponseWriter, r *http.Request) {
	incidentID := r.PathValue("id")

	if err := mh.useCase.AcknowledgeIncident(r.Context(), incidentID); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
