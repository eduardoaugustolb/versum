import { index, pgTable, smallint, uuid, varchar } from "drizzle-orm/pg-core";
import { testamentEnum } from "./bible.enums.ts"; // Corrected import path

export const bibleBooks = pgTable(
  "bible_books",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    name: varchar("name", { length: 100 }).notNull().unique(),
    order: smallint("order").notNull().unique(),
    testament: testamentEnum("testament").notNull(),
    total_chapters: smallint("total_chapters").notNull(),
  },
  (table) => [
    index("bible_books_testament_order_idx").on(table.testament, table.order),
  ],
);
