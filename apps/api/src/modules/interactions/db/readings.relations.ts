import { relations } from "drizzle-orm";
import { bibleChapters } from "../../bible/db/chapters.table.ts";
import { users } from "../../users/db/users.table.ts";
import { journeyReadings } from "./readings.table.ts";

export const journeyReadingsRelations = relations(
  journeyReadings,
  ({ one }) => ({
    user: one(users, {
      fields: [journeyReadings.userId],
      references: [users.id],
    }),
    chapter: one(bibleChapters, {
      fields: [journeyReadings.chapterId],
      references: [bibleChapters.id],
    }),
  }),
);
