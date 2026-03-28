import { relations } from "drizzle-orm";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";
import { marks } from "./marks.table.ts";

export const marksRelations = relations(marks, ({ one }) => ({
  user: one(users, {
    fields: [marks.userId],
    references: [users.id],
  }),
  verse: one(bibleVerses, {
    fields: [marks.verseId],
    references: [bibleVerses.id],
  }),
}));
