import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { DatePipe } from '@angular/common';
import { MonitorService } from './services/monitor.service';
import { Incident, Monitor } from './models/monitor.model';

@Component({
  selector: 'app-root',
  imports: [ReactiveFormsModule, DatePipe],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App implements OnInit {
  private readonly monitorService = inject(MonitorService);
  private readonly fb = inject(FormBuilder);

  // Estado reativo com Signals
  monitors = signal<Monitor[]>([]);
  selectedMonitor = signal<Monitor | null>(null);
  incidents = signal<Incident[]>([]);
  errorMessage = signal<string>('');

  // Formulário Reativo com validações
  monitorForm = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    targetUrl: ['https://', [Validators.required]],
    intervalSeconds: [10, [Validators.required, Validators.min(5)]],
    failureThreshold: [3, [Validators.required, Validators.min(1)]],
  });

  ngOnInit(): void {
    this.loadMonitors();
    // Atualiza a lista automaticamente a cada 5 segundos
    setInterval(() => this.loadMonitors(), 5000);
  }

  loadMonitors(): void {
    this.monitorService.getMonitors().subscribe({
      next: (data) => this.monitors.set(data),
      error: () => this.errorMessage.set('Erro ao carregar monitores da API Go.'),
    });
  }

  onSubmit(): void {
    if (this.monitorForm.invalid) {
      return;
    }

    this.errorMessage.set('');
    const payload = this.monitorForm.getRawValue();

    this.monitorService.createMonitor(payload).subscribe({
      next: () => {
        // Reseta o formulário para os valores padrão e recarrega a lista!
        this.monitorForm.reset({
          name: '',
          targetUrl: 'https://',
          intervalSeconds: 10,
          failureThreshold: 3,
        });
        this.loadMonitors();
      },
      error: (err) => {
        this.errorMessage.set(err.error?.erro || err.error?.error || 'Erro ao criar monitor.');
      },
    });
  }

  selectMonitor(monitor: Monitor): void {
    this.selectedMonitor.set(monitor);
    this.monitorService.getIncidents(monitor.id).subscribe({
      next: (data) => this.incidents.set(data),
    });
  }

  acknowledge(incidentId: string): void {
    this.monitorService.acknowledgeIncident(incidentId).subscribe({
      next: () => {
        const selected = this.selectedMonitor();
        if (selected) {
          this.selectMonitor(selected);
        }
      },
      error: () => this.errorMessage.set('Erro ao reconhecer incidente.'),
    });
  }
}