/** Nomes reais das tabelas no PostgreSQL (whitelist para TRUNCATE). */
export const TRUNCATABLE_TABLES = [
  { name: "bible_books", label: "Bíblia — livros" },
  { name: "bible_chapters", label: "Bíblia — capítulos" },
  { name: "bible_verses", label: "Bíblia — versículos" },
  { name: "readings", label: "Interações — leituras" },
  { name: "likes", label: "Interações — curtidas" },
  { name: "marks", label: "Interações — marcações" },
  { name: "users", label: "Usuários" },
  { name: "sessions", label: "Sessões (auth)" },
  { name: "magic_links", label: "Magic links (auth)" },
] as const;

export type TruncatableTableName = (typeof TRUNCATABLE_TABLES)[number]["name"];

const NAMES = new Set<string>(
  TRUNCATABLE_TABLES.map((t) => t.name),
);

export function assertTruncatableTableNames(names: string[]): asserts names is TruncatableTableName[] {
  for (const n of names) {
    if (!NAMES.has(n)) {
      throw new Error(`Tabela não permitida para TRUNCATE: ${n}`);
    }
  }
}

/** Presets: ordem irrelevante — o PostgreSQL resolve com CASCADE. */
export const TRUNCATE_PRESETS: Record<
  string,
  { title: string; description: string; tables: readonly TruncatableTableName[] }
> = {
  all: {
    title: "Todas as tabelas da aplicação",
    description:
      "Limpa bíblia, usuários, sessões, magic links e interações (leituras, likes, marks).",
    tables: [
      "readings",
      "likes",
      "marks",
      "bible_verses",
      "bible_chapters",
      "bible_books",
      "sessions",
      "magic_links",
      "users",
    ],
  },
  bible: {
    title: "Bíblia + interações em versículos",
    description:
      "Livros, capítulos, versículos e registros que referenciam versículos (leituras, likes, marks).",
    tables: [
      "readings",
      "likes",
      "marks",
      "bible_verses",
      "bible_chapters",
      "bible_books",
    ],
  },
  interactions: {
    title: "Só interações (leituras, likes, marks)",
    description: "Não altera bíblia nem usuários.",
    tables: ["readings", "likes", "marks"],
  },
  auth: {
    title: "Sessões e magic links",
    description: "Autenticação; não remove usuários.",
    tables: ["sessions", "magic_links"],
  },
  users_chain: {
    title: "Usuários e dados ligados",
    description:
      "Usuários, sessões e interações; não remove bíblia nem magic links pendentes.",
    tables: ["readings", "likes", "marks", "sessions", "users"],
  },
};
