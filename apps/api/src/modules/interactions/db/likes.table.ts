import { index, pgTable, timestamp, unique, uuid } from "drizzle-orm/pg-core";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";

export const likes = pgTable(
  "likes",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    userId: uuid("user_id")
      .notNull()
      .references(() => users.id, {
        onDelete: "cascade",
      }),
    verseId: uuid("verse_id")
      .notNull()
      .references(() => bibleVerses.id, {
        onDelete: "cascade",
      }),
    createdAt: timestamp("created_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    unique("likes_user_id_verse_id_unique").on(table.userId, table.verseId),
    index("likes_verse_id_idx").on(table.verseId),
  ],
);
