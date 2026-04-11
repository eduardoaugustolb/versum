import {
  boolean,
  index,
  pgTable,
  text,
  timestamp,
  uuid,
} from "drizzle-orm/pg-core";
import { bibleVerses } from "../../bible/db/verses.table.ts";
import { users } from "../../users/db/users.table.ts";

export const marks = pgTable(
  "marks",
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
    selectedVerseId: uuid("selected_verse_id")
      .notNull()
      .references(() => bibleVerses.id, {
        onDelete: "cascade",
      }),
    annotation: text("annotation"),
    isPublic: boolean("is_public").notNull().default(false),
    createdAt: timestamp("created_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    index("marks_user_id_created_at_idx").on(
      table.userId,
      table.createdAt.desc(),
    ),
    index("marks_verse_id_idx").on(table.verseId),
    index("marks_is_public_idx").on(table.isPublic),
  ],
);
