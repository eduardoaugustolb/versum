import type { consentLogs } from "../db/consent-logs.table";

export interface ConsentLogsRepo {
  createConsentLog(params: CreateConsentLogParams): Promise<ConsentLog>;
  createConsentLogs(params: CreateConsentLogParams[]): Promise<ConsentLog[]>;
  getConsentLogsByUserId(params: { userId: string }): Promise<ConsentLog[]>;
  hasConsent(params: { userId: string; purpose: string }): Promise<boolean>;
}

export type ConsentLog = typeof consentLogs.$inferSelect;


export interface CreateConsentLogParams {
  userId: string;
  purpose: string;
  granted: boolean;
  ip: string;
  userAgent: string;
}
