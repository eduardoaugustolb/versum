import { logger } from "@/utils/logger";
import fs, { existsSync, mkdirSync, writeFileSync } from "fs";
import { resolve } from "path";

// Lista dos livros e capítulos
const booksInfo = [
  { name: "genesis", chapters: 50, niceName: "Gênesis", slug: "gn" },
  { name: "exodo", chapters: 40, niceName: "Êxodo", slug: "ex" },
  { name: "levitico", chapters: 27, niceName: "Levítico", slug: "lv" },
  { name: "numeros", chapters: 36, niceName: "Números", slug: "nm" },
  { name: "deuteronomio", chapters: 34, niceName: "Deuteronômio", slug: "dt" },
  { name: "josue", chapters: 24, niceName: "Josué", slug: "js" },
  { name: "juizes", chapters: 21, niceName: "Juízes", slug: "jz" },
  { name: "rute", chapters: 4, niceName: "Rute", slug: "rt" },
  { name: "i-samuel", chapters: 31, niceName: "1 Samuel", slug: "1sm" },
  { name: "ii-samuel", chapters: 24, niceName: "2 Samuel", slug: "2sm" },
  { name: "i-reis", chapters: 22, niceName: "1 Reis", slug: "1rs" },
  { name: "ii-reis", chapters: 25, niceName: "2 Reis", slug: "2rs" },
  { name: "i-cronicas", chapters: 29, niceName: "1 Crônicas", slug: "1cr" },
  { name: "ii-cronicas", chapters: 36, niceName: "2 Crônicas", slug: "2cr" },
  { name: "esdras", chapters: 10, niceName: "Esdras", slug: "es" },
  { name: "neemias", chapters: 13, niceName: "Neemias", slug: "nmh" },
  { name: "tobias", chapters: 14, niceName: "Tobias", slug: "tb" },
  { name: "judite", chapters: 16, niceName: "Judite", slug: "jd" },
  { name: "ester", chapters: 10, niceName: "Ester", slug: "est" },
  { name: "jo", chapters: 42, niceName: "Jó", slug: "jo" },
  { name: "salmos", chapters: 150, niceName: "Salmos", slug: "sl" },
  { name: "i-macabeus", chapters: 16, niceName: "1 Macabeus", slug: "1mc" },
  { name: "ii-macabeus", chapters: 15, niceName: "2 Macabeus", slug: "2mc" },
  { name: "proverbios", chapters: 31, niceName: "Provérbios", slug: "pv" },
  { name: "eclesiastes", chapters: 12, niceName: "Eclesiastes", slug: "ec" },
  { name: "cantico-dos-canticos", chapters: 8, niceName: "Cântico dos Cânticos", slug: "ct" },
  { name: "sabedoria", chapters: 19, niceName: "Sabedoria", slug: "sb" },
  { name: "eclesiastico", chapters: 51, niceName: "Eclesiástico", slug: "eccl" },
  { name: "isaias", chapters: 66, niceName: "Isaías", slug: "is" },
  { name: "jeremias", chapters: 52, niceName: "Jeremias", slug: "jr" },
  { name: "lamentacoes", chapters: 5, niceName: "Lamentações", slug: "lm" },
  { name: "baruc", chapters: 6, niceName: "Baruc", slug: "br" },
  { name: "ezequiel", chapters: 48, niceName: "Ezequiel", slug: "ez" },
  { name: "daniel", chapters: 14, niceName: "Daniel", slug: "dn" },
  { name: "oseias", chapters: 14, niceName: "Oseias", slug: "os" },
  { name: "joel", chapters: 3, niceName: "Joel", slug: "jl" },
  { name: "amos", chapters: 9, niceName: "Amós", slug: "am" },
  { name: "abdias", chapters: 1, niceName: "Abdias", slug: "ab" },
  { name: "jonas", chapters: 4, niceName: "Jonas", slug: "jn" },
  { name: "miqueias", chapters: 7, niceName: "Miquéias", slug: "mq" },
  { name: "naum", chapters: 3, niceName: "Naum", slug: "nm" },
  { name: "habacuc", chapters: 3, niceName: "Habacuque", slug: "hc" },
  { name: "sofonias", chapters: 3, niceName: "Sofonias", slug: "sf" },
  { name: "ageu", chapters: 2, niceName: "Ageu", slug: "ag" },
  { name: "zacarias", chapters: 14, niceName: "Zacarias", slug: "zc" },
  { name: "malaquias", chapters: 4, niceName: "Malaquias", slug: "ml" },
  { name: "sao-mateus", chapters: 28, niceName: "São Mateus", slug: "mt" },
  { name: "sao-marcos", chapters: 16, niceName: "São Marcos", slug: "mc" },
  { name: "sao-lucas", chapters: 24, niceName: "São Lucas", slug: "lc" },
  { name: "sao-joao", chapters: 21, niceName: "São João", slug: "jo" },
  { name: "atos-dos-apostolos", chapters: 28, niceName: "Atos dos Apóstolos", slug: "at" },
  { name: "romanos", chapters: 16, niceName: "Romanos", slug: "rm" },
  { name: "i-corintios", chapters: 16, niceName: "1 Coríntios", slug: "1co" },
  { name: "ii-corintios", chapters: 13, niceName: "2 Coríntios", slug: "2co" },
  { name: "galatas", chapters: 6, niceName: "Gálatas", slug: "gl" },
  { name: "efesios", chapters: 6, niceName: "Efésios", slug: "ef" },
  { name: "filipenses", chapters: 4, niceName: "Filipenses", slug: "fp" },
  { name: "colossenses", chapters: 4, niceName: "Colossenses", slug: "cl" },
  { name: "i-tessalonicenses", chapters: 5, niceName: "1 Tessalonicenses", slug: "1ts" },
  { name: "ii-tessalonicenses", chapters: 3, niceName: "2 Tessalonicenses", slug: "2ts" },
  { name: "i-timoteo", chapters: 6, niceName: "1 Timóteo", slug: "1tm" },
  { name: "ii-timoteo", chapters: 4, niceName: "2 Timóteo", slug: "2tm" },
  { name: "tito", chapters: 3, niceName: "Tito", slug: "tt" },
  { name: "filemon", chapters: 1, niceName: "Filemom", slug: "fm" },
  { name: "hebreus", chapters: 13, niceName: "Hebreus", slug: "hb" },
  { name: "sao-tiago", chapters: 5, niceName: "São Tiago", slug: "tg" },
  { name: "i-sao-pedro", chapters: 5, niceName: "1 São Pedro", slug: "1pe" },
  { name: "ii-sao-pedro", chapters: 3, niceName: "2 São Pedro", slug: "2pe" },
  { name: "i-sao-joao", chapters: 5, niceName: "1 São João", slug: "1jo" },
  { name: "ii-sao-joao", chapters: 1, niceName: "2 São João", slug: "2jo" },
  { name: "iii-sao-joao", chapters: 1, niceName: "3 São João", slug: "3jo" },
  { name: "sao-judas", chapters: 1, niceName: "São Judas", slug: "jd" },
  { name: "apocalipse", chapters: 22, niceName: "Apocalipse", slug: "ap" }
];

// Exemplo de Bible JSON
const bible: unknown = JSON.parse(await Bun.file("./src/assets/json/bible.json").text());

// Função pra checar se o versículo existe, considerando intervalos "9-24"
function isVersePresent(chapter: any, verseNumber: number): boolean {
  for (const v of chapter.verses) {
    const match = v.text.match(/^(\d+)-(\d+)/); // pega padrão 9-24
    if (match) {
      const start = parseInt(match[1]);
      const end = parseInt(match[2]);
      if (verseNumber >= start && verseNumber <= end) return true;
    }
    if (v.verse === verseNumber) return true;
  }
  return false;
}

// Função de validação
export function validateBible(bible: any, booksInfo: any[]) {
  const errors: string[] = [];

  for (const b of booksInfo) {
    let book = bible.books.find((x: any) => x.name === b.name);

    if (!book) {
      errors.push(`Livro faltando: ${b.name}`);
      continue;
    }

    // Enriquecendo o JSON
    book.niceName = b.niceName;
    book.simpleName = b.name; // já vem sem acento
    book.slug = b.slug;

    for (let c = 1; c <= b.chapters; c++) {
      const chapter = book.chapters.find((x: any) => x.chapter === c);
      if (!chapter) {
        errors.push(`Capítulo faltando: ${b.name} ${c}`);
        continue;
      }

      // Encontrar maior número de versículo para o capítulo
      const maxVerse = Math.max(...chapter.verses.map((v: any) => v.verse));
      for (let v = 1; v <= maxVerse; v++) {
        if (!isVersePresent(chapter, v)) {
          errors.push(`Versículo faltando: ${b.name} ${c}:${v}`);
        }
      }
    }
  }

  return {errors, bible};
}

// Executa validação

export const executeValidation = () => {
  const result = validateBible(bible, booksInfo);
  logger("warn", `Total de erros encontrados: ${result.errors.length}`);
  return result;
};

export const writeEnrichedBible = async (bible: unknown, timestamp: Date) => {
  const timestampStr = timestamp.toISOString().replace(/:/g, "-");
  const folderName = "output" + "/" + timestampStr;
  const filename = `bible-enriched-parsed.json`;

  await Bun.file(folderName + "/" + filename).write(JSON.stringify(bible, null, 2));

  void logger("success", `Bible enriquecida salva em ${folderName}/${filename}`);
};

export const writeMissingVerses = async (missing: string[], timestamp: Date) => {
  const timestampStr = timestamp.toISOString().replace(/:/g, "-");
  const folderName = "output" + "/" + timestampStr;
  const filename = `missing-verses.json`;
  await Bun.file(folderName + "/" + filename).write(JSON.stringify(missing, null, 2));
  void logger("success", `Versículos faltando salvos em ${folderName}/${filename}`);
};
