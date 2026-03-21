// seed.ts

import { db } from "../api-deno/src/infrastructure/db/mod.ts";
import {
  bibleBooks,
  bibleChapters,
  bibleVerses,
} from "../api-deno/src/infrastructure/db/schema.ts";
import { and, eq, sql } from "drizzle-orm";

interface Versiculo {
  versiculo: number;
  texto: string;
}

interface Capitulo {
  capitulo: number;
  versiculos: Versiculo[];
}

interface Livro {
  nome: string;
  capitulos: Capitulo[];
}

interface Biblia {
  antigoTestamento: Livro[];
  novoTestamento: Livro[];
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

      const data = await response.json();
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

async function processLivro(
  livroData: Livro,
  ordem: number,
  testament: "OLD" | "NEW",
  existingBooks: any[],
  counters: {
    books: number;
    chapters: number;
    verses: number;
  },
) {
  console.log(`📖 ${livroData.nome} (${testament})`);

  const existing = existingBooks.find((b) => b.name === livroData.nome);
  let libro: any;

  if (existing) {
    libro = existing;
  } else {
    const [newLibro] = await db
      .insert(bibleBooks)
      .values({
        name: livroData.nome,
        order: ordem,
        testament: testament,
        total_chapters: livroData.capitulos.length,
      })
      .returning();

    libro = newLibro;
    counters.books++;

    await addLog(
      `📖 [${testament}] ${livroData.nome} - ${livroData.capitulos.length} caps`,
    );
  }

  for (const capituloData of livroData.capitulos) {
    const [existingChapter] = await db
      .select()
      .from(bibleChapters)
      .where(
        and(
          eq(bibleChapters.book_id, libro.id),
          eq(bibleChapters.number, capituloData.capitulo),
        ),
      );

    let capitulo = existingChapter;

    if (!existingChapter) {
      const [newCapitulo] = await db
        .insert(bibleChapters)
        .values({
          book_id: libro.id,
          number: capituloData.capitulo,
          total_verses: capituloData.versiculos.length,
        })
        .returning();

      capitulo = newCapitulo;
      counters.chapters++;
    }

    const [countResult] = await db
      .select({ count: sql<number>`count(*)` })
      .from(bibleVerses)
      .where(eq(bibleVerses.chapter_id, capitulo.id));

    const beforeCount = Number(countResult.count);

    const versiculos = capituloData.versiculos.map((vers) => ({
      chapter_id: capitulo.id,
      number: vers.versiculo,
      text: vers.texto,
    }));

    if (versiculos.length > 0) {
      await db.insert(bibleVerses).values(versiculos).onConflictDoNothing();
    }

    const [afterCountResult] = await db
      .select({ count: sql<number>`count(*)` })
      .from(bibleVerses)
      .where(eq(bibleVerses.chapter_id, capitulo.id));

    const afterCount = Number(afterCountResult.count);
    counters.verses += afterCount - beforeCount;
  }
}

async function main() {
  startTime = new Date();
  logs = [];
  endTime = null;
  hasError = false;
  messageId = null;

  console.log("🔥 SEED INICIADO");
  await addLog("🔥 SEED INICIADO");

  try {
    const url =
      "https://raw.githubusercontent.com/fidalgobr/bibliaAveMariaJSON/refs/heads/main/bibliaAveMariaRAW.json";

    const response = await fetch(url);

    if (!response.ok) throw new Error(`Falha no fetch: ${response.statusText}`);

    const bible = (await response.json()) as Biblia;

    await addLog(
      `✅ JSON OK - AT: ${bible.antigoTestamento.length} | NT: ${bible.novoTestamento.length}`,
    );

    const existingBooks = await db.select().from(bibleBooks);

    const counters = {
      books: 0,
      chapters: 0,
      verses: 0,
    };

    let ordem = 1;

    await addLog("📜 Processando Antigo Testamento...");
    for (const livro of bible.antigoTestamento) {
      await processLivro(livro, ordem++, "OLD", existingBooks, counters);
    }

    await addLog("✝️ Processando Novo Testamento...");
    for (const livro of bible.novoTestamento) {
      await processLivro(livro, ordem++, "NEW", existingBooks, counters);
    }

    endTime = new Date();

    await addLog("🎉 FINALIZADO!");
    await addLog(
      `📊 Livros: ${counters.books} | Capítulos: ${counters.chapters} | Versículos: ${counters.verses}`,
    );

    console.log("✅ FINALIZADO");
  } catch (err: any) {
    hasError = true;
    endTime = new Date();

    console.error("💀 ERRO:", err);
    await addLog(`❌ ERRO: ${err.message}`);
  }
}

await main();
