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
    text: text("text").notNull(),
  },
  (table) => [
    unique("bible_verses_chapter_id_number_unique").on(
      table.chapterId,
      table.number,
    ),
  ],
);
