import type { Context } from "hono";
import { SuccessViewModel } from "../../../../view-models/default/success.view-model.ts";
import type { Session } from "../../../auth/repositories/auth.types.repository.ts";
import type { VerseWithContext } from "../services/discovery.service.ts";
import { DiscoveryService } from "../services/discovery.service.ts";

export class DiscoveryController {
  private readonly service: DiscoveryService;

  constructor({ service }: { service?: DiscoveryService } = {}) {
    this.service = service ?? new DiscoveryService();
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
