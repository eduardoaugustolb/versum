import { and, desc, eq } from "drizzle-orm";
import { db } from "@/infrastructure/db";
import { InternalServerError } from "@/utils/app/errors";
import { consentLogs } from "../db/consent-logs.table";
import type { ConsentLog, ConsentLogsRepo, CreateConsentLogParams } from "./consent-logs.types.repository";

export class ConsentLogsRepository implements ConsentLogsRepo {
  constructor(
    private config: { dbInstance: typeof db } = { dbInstance: db }
  ) {

  }
  async createConsentLog(params: CreateConsentLogParams): Promise<ConsentLog> {
    const [res] = await this.config.dbInstance.insert(consentLogs).values(params).returning();
    if (!res) throw new InternalServerError("Failed to create consent log");
    return res;
  }
  async createConsentLogs(params: CreateConsentLogParams[]): Promise<ConsentLog[]> {
    return this.config.dbInstance.insert(consentLogs).values(params).returning();
  }
  async getConsentLogsByUserId({userId}: { userId: string }): Promise<ConsentLog[]> {
    return this.config.dbInstance.query.consentLogs.findMany({ where: (table) => eq(table.userId, userId)});
  }
  async hasConsent(params: { userId: string; purpose: string }): Promise<boolean> {
    const log = await this.config.dbInstance.query.consentLogs.findFirst({
      where: and(eq(consentLogs.userId, params.userId), eq(consentLogs.purpose, params.purpose)),
      orderBy: ({ createdAt }, { desc }) => desc(createdAt),
    });
    return !!log?.granted;
  }
}

/*
Método	Descrição
createConsentLog(params)	Insere 1 row
createConsentLogs(params[])	Insere múltiplos (batch)
getConsentLogsByUserId({ userId })	Retorna todo histórico do usuário
hasConsent({ userId, purpose })	Retorna true se o último registro para aquela finalidade for granted: true
*/
