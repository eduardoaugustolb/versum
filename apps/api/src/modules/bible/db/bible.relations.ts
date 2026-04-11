import { relations } from "drizzle-orm";
import { likes } from "../../interactions/db/likes.table.ts";
import { marks } from "../../interactions/db/marks.table.ts";
import { journeyReadings } from "../../interactions/db/readings.table.ts";
import { bibleBooks } from "./books.table.ts";
import { bibleChapters } from "./chapters.table.ts";
import { bibleVerses } from "./verses.table.ts";

export const bibleBooksRelations = relations(bibleBooks, ({ many }) => ({
  chapters: many(bibleChapters),
}));

export const bibleChaptersRelations = relations(
  bibleChapters,
  ({ one, many }) => ({
    book: one(bibleBooks, {
      fields: [bibleChapters.bookId],
      references: [bibleBooks.id],
    }),
    verses: many(bibleVerses),
  }),
);

export const bibleVersesRelations = relations(bibleVerses, ({ one, many }) => ({
  chapter: one(bibleChapters, {
    fields: [bibleVerses.chapterId],
    references: [bibleChapters.id],
  }),
  readings: many(journeyReadings),
  likes: many(likes),
  marks: many(marks),
}));
