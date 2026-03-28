import { db } from "../../../../infrastructure/db/index.ts";
import { bibleBooks } from "../../../../modules/bible/db/books.table.ts";
import { bibleChapters } from "../../../../modules/bible/db/chapters.table.ts";
import { bibleVerses } from "../../../../modules/bible/db/verses.table.ts";
import { and, eq, count } from "drizzle-orm";
import { logger } from "@/utils/logger";
import {
  normalizeBibleJsonForSeed,
  type NormalizedBook,
} from "../bible-json-normalize.ts";

function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

const WEBHOOK_URL =
  "https://discord.com/api/webhooks/1463937488962850837/kc9xkCsbzsXASyoxiMQDkWh3aQLsvejHYIf9CK6eABrg6QlGOdAjHHpg0MT5LO38zf6R";

let messageId: string | null = null;
let logs: string[] = [];
let startTime: Date = new Date();
let endTime: Date | null = null;
let hasError = false;

function formatDate(date: Date): string {
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${day}/${month}/${year} - ${hours}h${minutes}`;
}

async function updateDiscordMessage() {
  const logsText = logs.join("\n");
  const truncatedLogs =
    logsText.length > 3800 ? "...\n" + logsText.slice(-3800) : logsText;

  const embed = {
    title: `Logs do seed ${startTime.toISOString()}`,
    description: `**Logs:**\n\`\`\`\n${truncatedLogs}\n\`\`\``,
    fields: [
      {
        name: "Começou em:",
        value: formatDate(startTime),
        inline: true,
      },
      {
        name: "Terminou em:",
        value: endTime ? formatDate(endTime) : "Em andamento",
        inline: true,
      },
      {
        name: "Houve erros?:",
        value: hasError ? "Sim" : "Não",
        inline: true,
      },
    ],
    color: hasError ? 0xe74c3c : endTime ? 0x2ecc71 : 0xffaa00,
    timestamp: new Date().toISOString(),
  };

  try {
    if (!messageId) {
      const response = await fetch(WEBHOOK_URL + "?wait=true", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ embeds: [embed] }),
      });

      const data = (await response.json()) as { id: string };
      messageId = data.id;
    } else {
      await fetch(`${WEBHOOK_URL}/messages/${messageId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ embeds: [embed] }),
      });
    }
  } catch (err) {
    console.error("Erro ao atualizar Discord:", err);
  }
}

async function addLog(message: string) {
  logs.push(message);
  await updateDiscordMessage();
}

async function processBook(
  bookData: NormalizedBook,
  testament: "OLD" | "NEW",
  existingBooks: any[],
  options: {
    insertBooks: boolean;
    insertChapters: boolean;
    insertVerses: boolean;
  },
  counters: {
    books: number;
    chapters: number;
    verses: number;
  },
) {
  const slug = bookData.slug || slugify(bookData.name);
  console.log(`📖 ${bookData.name} (${slug}) - ${testament}`);

  const existing = existingBooks.find(
    (b) => b.name === bookData.name || b.slug === slug,
  );

  if (options.insertBooks) {
    if (existing) {
      await addLog(`♻️ [${testament}] ${bookData.name} - já existe`);
    } else {
      const [newBook] = await db
        .insert(bibleBooks)
        .values({
          name: bookData.name,
          slug: slug,
          niceName: bookData.niceName || bookData.name,
          testament: testament,
          total_chapters: bookData.chapters.length,
        })
        .returning();

      existingBooks.push(newBook);
      counters.books++;

      await addLog(
        `📖 [${testament}] ${bookData.name} (${slug}) - ${bookData.chapters.length} caps`,
      );
    }
  }

  const book =
    existing || existingBooks.find((b) => b.name === bookData.name);

  if (!book) {
    await addLog(
      `❌ Livro não encontrado após tentativa de inserção: ${bookData.name}`,
    );
    return;
  }

  for (const chapterData of bookData.chapters) {
    const [existingChapter] = await db
      .select()
      .from(bibleChapters)
      .where(
        and(
          eq(bibleChapters.book_id, book.id),
          eq(bibleChapters.number, chapterData.chapter),
        ),
      );

    let chapter = existingChapter;

    if (options.insertChapters && !existingChapter) {
      const [newChapter] = await db
        .insert(bibleChapters)
        .values({
          book_id: book.id,
          number: chapterData.chapter,
          total_verses: chapterData.verses.length,
        })
        .returning();

      chapter = newChapter;
      counters.chapters++;
    }

    if (!chapter) continue;

    if (options.insertVerses) {
      const [countResult] = await db
        .select({ total: count() })
        .from(bibleVerses)
        .where(eq(bibleVerses.chapter_id, chapter.id));

      const beforeCount = Number(countResult?.total ?? 0);

      const verses = chapterData.verses.map((vers) => ({
        chapter_id: chapter!.id,
        number: vers.verse,
        groupStart: vers.group_start,
        groupEnd: vers.group_end,
        text: vers.text,
      }));

      if (verses.length > 0) {
        await db
          .insert(bibleVerses)
          .values(verses)
          .onConflictDoNothing();
      }

      const [afterCountResult] = await db
        .select({ total: count() })
        .from(bibleVerses)
        .where(eq(bibleVerses.chapter_id, chapter!.id));

      const afterCount = Number(afterCountResult?.total ?? 0);
      counters.verses += afterCount - beforeCount;
    }
  }
}

export interface SeedBibleOptions {
  insertBooks: boolean;
  insertChapters: boolean;
  insertVerses: boolean;
}

export async function seedBibleFromJson(
  jsonPath: string,
  options: SeedBibleOptions,
) {
  startTime = new Date();
  logs = [];
  endTime = null;
  hasError = false;
  messageId = null;

  console.log("🔥 SEED INICIADO");
  await addLog("🔥 SEED INICIADO");

  try {
    const bibleJson = await Bun.file(jsonPath).text();
    const existingBooks = await db.select().from(bibleBooks);
    const bible = normalizeBibleJsonForSeed(JSON.parse(bibleJson), existingBooks);

    await addLog(`✅ JSON OK - Total de livros: ${bible.books.length}`);

    const counters = {
      books: 0,
      chapters: 0,
      verses: 0,
    };

    await addLog("📜 Processando livros...");

    const oldTestamentCount = 46;

    for (const book of bible.books) {
      const index = bible.books.indexOf(book);
      const testament: "OLD" | "NEW" =
        index < oldTestamentCount ? "OLD" : "NEW";
      await processBook(book, testament, existingBooks, options, counters);
    }

    endTime = new Date();

    await addLog("🎉 FINALIZADO!");
    await addLog(
      `📊 Livros: ${counters.books} | Capítulos: ${counters.chapters} | Versículos: ${counters.verses}`,
    );

    logger("success", "Seed da bíblia concluído!");
    console.log("✅ FINALIZADO");
  } catch (err: any) {
    hasError = true;
    endTime = new Date();

    console.error("💀 ERRO:", err);
    await addLog(`❌ ERRO: ${err.message}`);
    logger("error", `Erro no seed: ${err.message}`);
  }
}
