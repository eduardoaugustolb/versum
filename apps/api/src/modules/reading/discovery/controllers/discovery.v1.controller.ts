import type { Context } from "hono";
import { SuccessViewModel } from "@/view-models/default/success.view-model.ts";
import type { Session } from "@/modules/auth/repositories/auth.types.repository.ts";
import type { VerseWithContext } from "../services/discovery.v1.service.ts";
import { DiscoveryServiceV1 } from "../services/discovery.v1.service.ts";

export class DiscoveryControllerV1 {
  private readonly service: DiscoveryServiceV1;

  constructor({ service }: { service?: DiscoveryServiceV1 } = {}) {
    this.service = service ?? new DiscoveryServiceV1();
  }

  getNextVerses = async (c: Context) => {
    const chapterId = c.req.query("chapterId");

    if (!chapterId) {
      return c.json(SuccessViewModel.create<VerseWithContext[]>(), 200);
    }

    const verses = await this.service.getNextVerses(chapterId);
    return c.json(SuccessViewModel.create(verses), 200);
  };

  markVersesAsRead = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = (await c.req.json()) as { verseIds: string[] };

    await this.service.markVersesAsRead({
      userId: session.userId,
      verseIds: body.verseIds,
    });

    return c.json(SuccessViewModel.create({ success: true }), 200);
  };

  getStats = async (c: Context) => {
    const session = c.get("session") as Session;
    const stats = await this.service.getStats(session.userId);
    return c.json(SuccessViewModel.create(stats), 200);
  };
}
