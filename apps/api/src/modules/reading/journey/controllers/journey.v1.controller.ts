import type { Context } from "hono";
import { ZodError } from "zod";
import { SuccessViewModel } from "@/view-models/default/success.view-model.ts";
import type { Session } from "@/modules/auth/repositories/auth.types.repository.ts";
import { JourneyServiceV1 } from "../services/journey.v1.service.ts";

export class JourneyControllerV1 {
  private readonly service: JourneyServiceV1;

  constructor({ service }: { service?: JourneyServiceV1 } = {}) {
    this.service = service ?? new JourneyServiceV1();
  }

  getFeed = async (c: Context) => {
    const session = c.get("session") as Session;

    const rawBufferSize = c.req.query("buffer-size");
    const bufferSize =
      rawBufferSize !== undefined ? parseInt(rawBufferSize, 10) : 4;

    if (Number.isNaN(bufferSize) || bufferSize < 0 || bufferSize > 4) {
      throw new ZodError([
        {
          code: "custom",
          message: "buffer-size must be a number between 0 and 4",
          path: ["buffer-size"],
        },
      ]);
    }

    const feed = await this.service.getFeed(session.userId, bufferSize);

    return c.json(SuccessViewModel.create(feed), 200);
  };

  markCurrentAsRead = async (c: Context) => {
    const session = c.get("session") as Session;

    const result = await this.service.markCurrentAsRead(session.userId);

    return c.json(SuccessViewModel.create(result), 200);
  };

  getStatus = async (c: Context) => {
    const session = c.get("session") as Session;

    const status = await this.service.getStatus(session.userId);

    return c.json(SuccessViewModel.create(status), 200);
  };
}
