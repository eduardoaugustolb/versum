import { index, pgTable, smallint, uuid, varchar } from "drizzle-orm/pg-core";
import { testamentEnum } from "./bible.enums.ts";

export const bibleBooks = pgTable(
  "bible_books",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    order: smallint("order").notNull(),
    name: varchar("name", { length: 100 }).notNull().unique(),
    slug: varchar("slug", { length: 10 }).notNull().unique(),
    niceName: varchar("nice_name", { length: 100 }).notNull(),
    testament: testamentEnum("testament").notNull(),
    totalChapters: smallint("total_chapters").notNull(),
  },
  (table) => [
    index("bible_books_testament_slug_idx").on(table.testament, table.slug),
    index("bible_books_slug_idx").on(table.slug),
    index("bible_books_nice_name_idx").on(table.niceName),
    index("bible_books_name_idx").on(table.name),
    index("bible_books_name_slug_idx").on(table.name, table.slug),
    index("bible_books_slug_nice_name_idx").on(table.slug, table.niceName),
    index("bible_books_order_idx").on(table.order),
  ],
);
