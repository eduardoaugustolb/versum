import { index, pgTable, smallint, text, unique, uuid } from "drizzle-orm/pg-core";
import { bibleChapters } from "./chapters.table.ts";

export const bibleVerses = pgTable(
  "bible_verses",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    chapterId: uuid("chapter_id")
      .notNull()
      .references(() => bibleChapters.id, {
        onDelete: "cascade",
      }),
    number: smallint("number").notNull(),
    groupStart: smallint("group_start"),
    groupEnd: smallint("group_end"),
    text: text("text").notNull(),
  },
  (table) => [
    index("bible_verses_group_start_idx").on(table.groupStart),
    index("bible_verses_group_end_idx").on(table.groupEnd),
    unique("bible_verses_chapter_id_number_unique").on(
      table.chapterId,
      table.number,
    ),
  ],
);
