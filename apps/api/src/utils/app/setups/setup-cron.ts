import { logger } from "@versum/logger";
import { env } from "../../env/parser.ts";
import { PurgeService } from "../../../modules/auth/services/purge.service.ts";

export class SetupCron {
  constructor() {
    if (env.CRON_ENABLED !== "true") {
      logger("info", "[CRON] Desabilitado via CRON_ENABLED env var");
      return;
    }
    this.setupDailyPurge();
  }

  private setupDailyPurge() {
    Bun.cron("0 3 * * *", async () => {
      const service = new PurgeService();
      const result = await service.runDailyPurge();
      logger(
        { level: "info", icon: "🧹" },
        `[PURGE] Magic links: ${result.magicLinks} deletados, Sessões: ${result.sessions} deletadas`,
      );
    });

    logger("info", "[CRON] Purge job registrado — 03:00 diário");
  }
}
