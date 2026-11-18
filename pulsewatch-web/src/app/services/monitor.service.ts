import { inject, Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { CreateMonitorRequest, Incident, Monitor } from "../models/monitor.model";
import { HttpClient } from "@angular/common/http";

@Injectable({
  providedIn: 'root',
})
export class MonitorService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api';

  // Listar todos os monitores (GET /api/monitors)
  getMonitors(): Observable<Monitor[]> {
    return this.http.get<Monitor[]>(`${this.apiUrl}/monitors`);
  }

  // Criar um novo monitor (POST /api/monitors)
  createMonitor(payload: CreateMonitorRequest): Observable<Monitor> {
    return this.http.post<Monitor>(`${this.apiUrl}/monitors`, payload)
  }

  // Listar incidentes de um monitor (GET /api/monitors/{id}/incidents)
  getIncidents(monitorID: string): Observable<Incident[]> {
    return this.http.get<Incident[]>(`${this.apiUrl}/monitors/${monitorID}/incidents`)
  }

  acknowledgeIncident(incidentId: string): Observable<void> {
    return this.http.patch<void>(`${this.apiUrl}/incidents/${incidentId}/acknowledge`, {});
    }
}