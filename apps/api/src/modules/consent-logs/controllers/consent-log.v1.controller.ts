import type { Context } from "hono";
import { BadRequestError } from "../../../utils/app/errors/index";
import type { Session } from "../../auth/repositories/auth.types.repository";
import { SuccessViewModel } from "@/view-models/default/success.view-model.ts";
import { ConsentLogServiceV1 } from "../services/consent-log.v1.service";

const consentLogService = new ConsentLogServiceV1();

export class ConsentLogControllerV1 {
  private readonly service: ConsentLogServiceV1;

  constructor({ service }: { service?: ConsentLogServiceV1 } = {}) {
    this.service = service ?? consentLogService;
  }

  recordConsent = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = (await c.req.json()) as { consents?: Array<{ purpose: string; granted: boolean }> };

    if (!body.consents || body.consents.length === 0) {
      throw new BadRequestError("At least one consent must be provided");
    }

    const ip = c.req.header("x-forwarded-for")?.split(",")[0]?.trim() ?? "unknown";
    const userAgent = c.req.header("user-agent") ?? "unknown";

    const logs = await this.service.recordConsents({
      userId: session.userId,
      consents: body.consents,
      ip,
      userAgent,
    });

    return c.json(SuccessViewModel.create({ consents: logs }), 201);
  };

  getConsentHistory = async (c: Context) => {
    const session = c.get("session") as Session;

    const logs = await this.service.getUserConsents({ userId: session.userId });

    return c.json(SuccessViewModel.create({ consents: logs }), 200);
  };
}
