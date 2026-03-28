export interface NormalizedVerse {
  verse: number;
  text: string;
  group_start?: number;
  group_end?: number;
}

export interface NormalizedChapter {
  chapter: number;
  verses: NormalizedVerse[];
}

export interface NormalizedBook {
  name: string;
  niceName: string;
  slug: string;
  chapters: NormalizedChapter[];
}

export interface NormalizedBible {
  books: NormalizedBook[];
}

interface CompactVerseWithGroup {
  text: string;
  group_start?: number;
  group_end?: number;
}

type CompactChapter = Record<string, string | CompactVerseWithGroup>;

function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export function normalizeBookKey(value: string): string {
  return value
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]/g, "");
}

function toNiceName(value: string): string {
  return value
    .replace(/([a-z])([A-Z])/g, "$1 $2")
    .replace(/[_-]+/g, " ")
    .trim()
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function toVerse(
  input: string | CompactVerseWithGroup,
  verse: number,
): NormalizedVerse {
  if (typeof input === "string") {
    return { verse, text: input };
  }
  return {
    verse,
    text: input.text,
    group_start: input.group_start,
    group_end: input.group_end,
  };
}

function isBookNode(
  raw: unknown,
): raw is {
  name?: string;
  slug?: string;
  niceName?: string;
  chapters: CompactChapter;
} {
  if (typeof raw !== "object" || raw === null) return false;
  if (!("chapters" in raw)) return false;
  const ch = (raw as { chapters: unknown }).chapters;
  return typeof ch === "object" && ch !== null && !Array.isArray(ch);
}

/**
 * Formato esperado (`bible.json` atual):
 * `{ "genesis": { "name": "...", "slug": "...", "niceName": "...", "chapters": { "1": { "1": "..." } } } }`
 */
export function normalizeBibleJsonForSeed(
  parsedJson: unknown,
  existingBooks: Array<{ name: string; slug: string; niceName: string }>,
): NormalizedBible {
  if (typeof parsedJson !== "object" || parsedJson === null || Array.isArray(parsedJson)) {
    throw new Error("bible.json deve ser um objeto na raiz.");
  }

  const root = parsedJson as Record<string, unknown>;

  if ("books" in root) {
    throw new Error(
      'Formato com "books"[] não é suportado. Use chaves por livro e um objeto "chapters" por livro.',
    );
  }

  const books = Object.entries(root).map(([bookKey, rawBook]) => {
    if (!isBookNode(rawBook)) {
      throw new Error(
        `Livro "${bookKey}": cada entrada deve ter "chapters" como mapa (capítulo → versículos).`,
      );
    }

    const normalizedKey = normalizeBookKey(bookKey);
    const existing = existingBooks.find((book) =>
      [book.name, book.slug, book.niceName]
        .filter(Boolean)
        .map(normalizeBookKey)
        .includes(normalizedKey),
    );

    const chapters = Object.entries(rawBook.chapters)
      .map(([chapterNumber, rawChapter]) => {
        if (typeof rawChapter !== "object" || rawChapter === null || Array.isArray(rawChapter)) {
          throw new Error(
            `Livro "${bookKey}" capítulo ${chapterNumber}: esperado objeto de versículos.`,
          );
        }
        const verses = Object.entries(rawChapter as Record<string, string | CompactVerseWithGroup>)
          .map(([verseNumber, verseContent]) =>
            toVerse(verseContent, Number(verseNumber)),
          )
          .sort((a, b) => a.verse - b.verse);

        return {
          chapter: Number(chapterNumber),
          verses,
        };
      })
      .sort((a, b) => a.chapter - b.chapter);

    return {
      name: existing?.name || rawBook.name || bookKey,
      niceName: existing?.niceName || rawBook.niceName || toNiceName(bookKey),
      slug: existing?.slug || rawBook.slug || slugify(bookKey).slice(0, 10),
      chapters,
    };
  });

  return { books };
}
