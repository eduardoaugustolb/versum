import { relations } from "drizzle-orm";
import { bibleBooks } from "./books.table.ts";
import { bibleChapters } from "./chapters.table.ts";
import { bibleVerses } from "./verses.table.ts";
import { readings } from "../../interactions/db/readings.table.ts";
import { likes } from "../../interactions/db/likes.table.ts";
import { marks } from "../../interactions/db/marks.table.ts";

export const bibleBooksRelations = relations(bibleBooks, ({ many }) => ({
  chapters: many(bibleChapters),
}));

export const bibleChaptersRelations = relations(
  bibleChapters,
  ({ one, many }) => ({
    book: one(bibleBooks, {
      fields: [bibleChapters.book_id],
      references: [bibleBooks.id],
    }),
    verses: many(bibleVerses),
  }),
);

export const bibleVersesRelations = relations(bibleVerses, ({ one, many }) => ({
  chapter: one(bibleChapters, {
    fields: [bibleVerses.chapter_id],
    references: [bibleChapters.id],
  }),
  readings: many(readings),
  likes: many(likes),
  marks: many(marks),
}));
