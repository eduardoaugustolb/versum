import { CONSENT_PURPOSES } from "../schemas/v1/consent-log.v1.common.schema";
import { ConsentLogsRepository } from "../repositories/consent-logs.repository";
import type { ConsentLog, CreateConsentLogParams } from "../repositories/consent-logs.types.repository";

interface ConsentInput {
  purpose: string;
  granted: boolean;
}

export class ConsentLogServiceV1 {
  private readonly repository: ConsentLogsRepository;

  constructor({ repository }: { repository?: ConsentLogsRepository } = {}) {
    this.repository = repository ?? new ConsentLogsRepository();
  }

  async recordConsents(params: {
    userId: string;
    consents: ConsentInput[];
    ip: string;
    userAgent: string;
  }): Promise<ConsentLog[]> {
    for (const consent of params.consents) {
      if (!CONSENT_PURPOSES.includes(consent.purpose as typeof CONSENT_PURPOSES[number])) {
        throw new Error(`Invalid consent purpose: ${consent.purpose}`);
      }
    }

    const rows: CreateConsentLogParams[] = params.consents.map((c) => ({
      userId: params.userId,
      purpose: c.purpose,
      granted: c.granted,
      ip: params.ip,
      userAgent: params.userAgent,
    }));

    return this.repository.createConsentLogs(rows);
  }

  async getUserConsents({ userId }: { userId: string }): Promise<ConsentLog[]> {
    return this.repository.getConsentLogsByUserId({ userId });
  }

  async hasConsent(params: { userId: string; purpose: string }): Promise<boolean> {
    return this.repository.hasConsent(params);
  }
}
