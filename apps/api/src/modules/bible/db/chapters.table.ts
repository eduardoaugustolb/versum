import { integer, pgTable, unique, uuid } from "drizzle-orm/pg-core";
import { bibleBooks } from "./books.table.ts";

export const bibleChapters = pgTable(
  "bible_chapters",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    book_id: uuid("book_id")
      .notNull()
      .references(() => bibleBooks.id, { onDelete: "cascade" }),
    number: integer("number").notNull(),
    total_verses: integer("total_verses").notNull(),
  },
  (table) => [
    unique("bible_chapters_book_id_number_unique").on(
      table.book_id,
      table.number,
    ),
  ],
);
