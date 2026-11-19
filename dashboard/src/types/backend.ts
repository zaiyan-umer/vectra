export type BackendStatus = "healthy" | "degraded" | "down";

export interface BackendTarget {
  name: string;
  host: string;
  port: number;
  status: BackendStatus;
  avgLatency: number;
  weight: number;
}
