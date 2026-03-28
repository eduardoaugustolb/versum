import { integer, pgTable, unique, uuid } from "drizzle-orm/pg-core";
import { bibleBooks } from "./books.table.ts";

export const bibleChapters = pgTable(
  "bible_chapters",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    bookId: uuid("book_id")
      .notNull()
      .references(() => bibleBooks.id, { onDelete: "cascade" }),
    number: integer("number").notNull(),
    totalVerses: integer("total_verses").notNull(),
  },
  (table) => [
    unique("bible_chapters_book_id_number_unique").on(
      table.bookId,
      table.number,
    ),
  ],
);
