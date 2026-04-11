import { index, pgTable, timestamp, uuid } from "drizzle-orm/pg-core";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";

export const discoveryReadings = pgTable(
  "discovery_readings",
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
    readAt: timestamp("read_at", {
      withTimezone: true,
      precision: 3,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    index("discovery_readings_user_id_read_at_idx").on(
      table.userId,
      table.readAt.desc(),
    ),
    index("discovery_readings_verse_id_idx").on(table.verseId),
    index("discovery_readings_user_verse_unique_idx").on(
      table.userId,
      table.verseId,
    ),
  ],
);
