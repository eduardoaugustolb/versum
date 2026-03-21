import { index, pgTable, timestamp, uuid } from "drizzle-orm/pg-core";
import { users } from "../../users/db/users.table.ts";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { readingModeEnum } from "./interactions.enums.ts";

export const readings = pgTable(
  "readings",
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
    mode: readingModeEnum("mode").notNull(),
    read_at: timestamp("read_at", {
      withTimezone: true,
      precision: 3,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    index("readings_user_id_read_at_idx").on(
      table.user_id,
      table.read_at.desc(),
    ),
    index("readings_verse_id_idx").on(table.verse_id),
  ],
);
