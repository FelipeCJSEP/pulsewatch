export type MonitorStatus = 'PENDING' | 'UP' | 'DOWN' | 'PAUSED';
export type IncidentStatus = 'OPEN' | 'ACKNOWLEDGED' | 'RESOLVED';

export interface Monitor {
    id: string;
    name: string;
    targetUrl: string;
    intervalSeconds: number;
    failureThreshold: number;
    consecutiveFailures: number;
    status: MonitorStatus;
    lastCheckedAt: string | null;
    createdAt: string;
}

export interface CreateMonitorRequest {
    name: string;
    targetUrl: string;
    intervalSeconds: number;
    failureThreshold: number;
}

export interface Incident {
    id: string;
    monitorId: string;
    reason: string;
    status: IncidentStatus;
    openedAt: string;
    resolvedAt: string | null;
}