import {
  boolean,
  index,
  pgTable,
  text,
  timestamp,
  uuid,
} from "drizzle-orm/pg-core";
import { users } from "../../users/db/users.table.ts";
import { bibleVerses } from "../../bible/db/verses.table.ts";

export const marks = pgTable(
  "marks",
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
    selected_text: text("selected_text").notNull(),
    annotation: text("annotation"),
    is_public: boolean("is_public").notNull().default(false),
    created_at: timestamp("created_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    index("marks_user_id_created_at_idx").on(
      table.user_id,
      table.created_at.desc(),
    ),
    index("marks_verse_id_idx").on(table.verse_id),
    index("marks_is_public_idx").on(table.is_public),
  ],
);
