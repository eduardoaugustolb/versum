import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";
import { readings } from "./readings.table.ts";
import { relations } from "drizzle-orm";

export const readingsRelations = relations(readings, ({ one }) => ({
  user: one(users, {
    fields: [readings.user_id],
    references: [users.id],
  }),
  verse: one(bibleVerses, {
    fields: [readings.verse_id],
    references: [bibleVerses.id],
  }),
}));
