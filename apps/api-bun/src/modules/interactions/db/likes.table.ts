import { pgTable, timestamp, uuid, unique, index } from "drizzle-orm/pg-core";
import { users } from "../../users/db/users.table.ts";
import { bibleVerses } from "../../bible/db/verses.table.ts";

export const likes = pgTable(
  "likes",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    user_id: uuid("user_id")
      .notNull()
      .references(() => users.id, {
        onDelete: "cascade",
      }),
    verse_id: uuid("verse_id")
      .notNull()
      .references(() => bibleVerses.id, {
        onDelete: "cascade",
      }),
    created_at: timestamp("created_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    unique("likes_user_id_verse_id_unique").on(table.user_id, table.verse_id),
    index("likes_verse_id_idx").on(table.verse_id),
  ],
);
