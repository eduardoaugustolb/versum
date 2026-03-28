import { relations } from "drizzle-orm";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";
import { likes } from "./likes.table.ts";

export const likesRelations = relations(likes, ({ one }) => ({
  user: one(users, {
    fields: [likes.userId],
    references: [users.id],
  }),
  verse: one(bibleVerses, {
    fields: [likes.verseId],
    references: [bibleVerses.id],
  }),
}));
