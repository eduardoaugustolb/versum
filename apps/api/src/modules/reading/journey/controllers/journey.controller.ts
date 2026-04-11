import type { Context } from "hono";
import { ZodError } from "zod";
import { SuccessViewModel } from "../../../../view-models/default/success.view-model.ts";
import type { Session } from "../../../auth/repositories/auth.types.repository.ts";
import { JourneyService } from "../services/journey.service.ts";

export class JourneyController {
  private readonly service: JourneyService;

  constructor({ service }: { service?: JourneyService } = {}) {
    this.service = service ?? new JourneyService();
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
