import { index, pgTable, timestamp, unique, uuid } from "drizzle-orm/pg-core";
import { bibleChapters } from "../../bible/db/chapters.table.ts";
import { users } from "../../users/db/users.table.ts";

export const journeyReadings = pgTable(
  "journey_readings",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    userId: uuid("user_id")
      .notNull()
      .references(() => users.id, {
        onDelete: "cascade",
      }),
    chapterId: uuid("chapter_id")
      .notNull()
      .references(() => bibleChapters.id, {
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
    unique("journey_readings_user_chapter_unique_idx").on(
      table.userId,
      table.chapterId,
    ),
    index("journey_readings_user_id_read_at_idx").on(
      table.userId,
      table.readAt.desc(),
    ),
    index("journey_readings_chapter_id_idx").on(table.chapterId),
  ],
);
