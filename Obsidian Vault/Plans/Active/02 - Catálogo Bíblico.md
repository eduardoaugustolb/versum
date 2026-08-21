---
title: "02 - Catálogo Bíblico"
section: Plans
subsection: Active
type: implementation-plan
status: completed
date: 2026-08-04
tags: [versum, plans, api, go, catalog, postgres]
up: "[[Plans/Active/_Index|Planos Ativos]]"
prev: "[[Plans/Active/_Index|Planos Ativos]]"
next: "[[Plans/Archive/_Index|Arquivo]]"
related: ["[[Docs/Architecture/Visão Geral]]", "[[Rules/01 - Princípios de Engenharia]]", "[[Plans/Archive/02 - Corpus Bíblico Canônico]]", "[[Plans/Archive/01 - Fundação da API Go]]"]
---

# Catálogo Bíblico

## Objetivo

Expor o corpus canônico (`bible/corpus/v1`) pela API, para que o web possa
explorar e ler livros, capítulos e versículos sem exigir conta — a leitura
pública já está permitida por
[[Docs/Architecture/Autenticação e Sessões|Autenticação e Sessões]]. Esta
entrega introduz o primeiro domínio de leitura real, com Postgres como fonte
de consulta e o corpus versionado como fonte de carga.

## Arquitetura

O módulo `catalog` expõe duas queries (`ListBooks` e `GetChapter`) e um
command (`PublishBook`). As entidades vivem em `catalog/domain`, as portas
definidas pela aplicação em `catalog/application/ports`, os casos de uso em
`application/queries` e `application/commands`, e o adapter SQL em
`catalog/postgres`.

Os repositórios retornam entidades de domínio, não DTOs HTTP. A camada de
transporte converte essas entidades para respostas JSON. `PublishBook` recebe
um input externo, constrói e valida `Book` e `Verse`, e executa a escrita por
uma porta transacional que fornece um writer vinculado à transação.

### Estado atual

As etapas de estrutura, API HTTP, repositório PostgreSQL, transação por livro,
gravação de `catalog_version`, tratamento de falhas do seed e cobertura de CI
estão implementadas.

> [!note] Histórico da decisão
> Essa camada foi implementada, removida e reimplementada nesta mesma
> entrega. Primeiro existia como `sql.Executor`/`PgxExecutor`
> (`internal/adapters/sql`), justificada por portabilidade de banco. Uma
> avaliação arquitetural externa (`AVALIACAO_ARQUITETURAL.md`, 5 de agosto
> de 2026) apontou que essa portabilidade era ilusória — o texto das
> queries (`$1`, `ON CONFLICT`, `"order"`) já é específico de Postgres — e
> a camada foi removida, com o repositório do catálogo passando a depender de
> `pgx` diretamente via uma interface privada (`dbtx`, com tipos do `pgx`
> na própria assinatura). Isso reabriu a discussão: o motivo certo pra
> isolar o driver nunca foi portabilidade de banco, é a segunda cláusula do
> Princípio da Inversão de Dependência — abstrações não devem depender de
> detalhes, e `dbtx` (com `pgx.Rows`/`pgx.Row`/`pgconn.CommandTag` na
> assinatura) violava exatamente isso. A porta foi recriada como
> `internal/ports/dbexec` (nome novo, evita a colisão com
> `database/sql` que o `internal/adapters/sql` original tinha), com tipos
> só nossos. Ver [[Rules/01 - Princípios de Engenharia]], seção "Isolar
> driver/framework atrás de uma porta".

```go
// internal/ports/dbexec/executor.go — porta, sem tipo de pgx
package dbexec

type Row interface {
    Scan(dest ...any) error
}

type Rows interface {
    Row
    Next() bool
    Err() error
    Close()
}

type Executor interface {
    QueryRow(ctx context.Context, sql string, args ...any) Row
    Query(ctx context.Context, sql string, args ...any) (Rows, error)
    Exec(ctx context.Context, sql string, args ...any) error
}
```

```go
// internal/adapters/postgres/pgx_executor.go — adapter, conhece pgx
type pgxConn interface {
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type PgxExecutor struct{ conn pgxConn }

func NewPgxExecutor(conn pgxConn) PgxExecutor { return PgxExecutor{conn: conn} }
```

```go
// internal/catalog/postgres/repository.go
type Repository struct {
    db dbexec.Executor
}

func NewRepository(db dbexec.Executor) *Repository {
    return &Repository{db: db}
}
```

`pgxConn` é satisfeita tanto por `*pgxpool.Pool` (leitura) quanto por
`pgx.Tx` (seed, etapa 8), então o mesmo `PgxExecutor` serve pros dois — sem
precisar de dois adapters. O repositório continua fora de
`internal/adapters/` (fica em `internal/catalog`, junto do domínio que
serve) porque o texto das queries é Postgres-específico de qualquer forma;
só a execução em si (`pgx`) fica isolada em `PgxExecutor`.

O corpus em `bible/corpus/v1` continua sendo a única fonte canônica do
conteúdo bíblico (ver
[[Plans/Archive/02 - Corpus Bíblico Canônico|Corpus Bíblico Canônico]]).
Postgres aqui é uma **projeção de leitura** desse corpus, carregada por um
comando de seed idempotente — não um lugar onde o conteúdo é editado. Se o
corpus for revalidado ou ganhar uma v2, o seed roda de novo; ele não lê nem
escreve fora do banco em tempo de request.

Os casos de uso não conhecem SQL, `pgx` ou `chi`. O handler HTTP não decide o
que é um "capítulo válido" — só traduz path params e erros do domínio para
HTTP.

## Stack

- Go 1.26.5
- `github.com/jackc/pgx/v5` e `github.com/jackc/pgx/v5/pgxpool`
- CLI `migrate` (`github.com/golang-migrate/migrate/v4/cmd/migrate`) para
  criar e aplicar migrations — instalado à parte
  (`go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`),
  não é dependência do binário da API. `cmd/api` não migra sozinho: assume
  que o schema já foi aplicado antes de subir. Aplicar migration é passo de
  dev/deploy, não uma responsabilidade em tempo de execução do servidor —
  evita acoplar o boot HTTP a estado "dirty" de migration ou a corrida entre
  múltiplas instâncias migrando ao mesmo tempo.
- Postgres local via `infra/docker-compose.yml` (serviço único, vinculado só
  em `127.0.0.1`, sem outras dependências de infra ainda)
- `net/http`, `github.com/go-chi/chi/v5`, `log/slog`, `go test`. `chi`
  fica isolado atrás de `internal/ports/httprouter.Router` — só
  `router.go` importa `chi` diretamente; `health_routes.go` e
  `catalog_routes.go` usam `net/http` puro, incluindo `r.PathValue(...)`
  (Go 1.22+) no lugar de `chi.URLParam` pra ler path params.

## Critérios de aceitação

1. `infra/docker-compose.yml` sobe um Postgres vazio utilizável localmente,
   vinculado só em `127.0.0.1`.
2. `migrate -path internal/adapters/postgres/migrations -database "$DATABASE_URL" up`
   aplica o schema; rodar de novo não reaplica (`no change`).
3. `go run ./cmd/seed-catalog` carrega os 73 livros de `bible/corpus/v1` no
   banco; rodar duas vezes não duplica dado (idempotente).
4. `GET /books` responde `200` com os 73 livros, ordenados por `order`.
5. `GET /books/{bookId}/chapters/{number}` responde `200` com os versículos
   do capítulo em ordem, preservando repetições de numeração (`9.1`, `9.2`)
   como partes distintas do mesmo versículo, e com `book_name` preenchido.
6. `GET /books/{bookId}/chapters/{number}` responde `404` para `bookId` ou
   `number` inexistentes, e `400` para `number` não-numérico ou `<= 0`.
7. `go build ./...`, `go vet ./...`, `go test ./...` e `go test -race ./...`
   passam no diretório `api/`, localmente e no CI.
8. `cmd/api` recusa subir se o Postgres estiver inacessível (`Ping` no
   boot), e encerra de forma graciosa (`SIGTERM`/`SIGINT` drenam requisições
   em curso antes de fechar).

## Estrutura de arquivos

```text
api/
  cmd/seed/catalog/main.go                      # seed; usa PublishBook com transação por livro
  internal/config/config.go                    # + DatabaseURL
  internal/ports/dbexec/
    executor.go                                 # porta Executor/Row/Rows/CopyFrom, sem tipo de pgx
  internal/ports/httprouter/
    router.go                                   # porta Router, sem tipo de chi
  internal/adapters/postgres/
    migrations/
      000001_create-catalog.up.sql
      000001_create-catalog.down.sql
      000002_catalog-version.up.sql
      000002_catalog-version.down.sql
    pgx_executor.go                              # PgxExecutor implementa dbexec.Executor usando pgx
  internal/catalog/
    domain/                                     # Book, Chapter, Verse e invariantes
    application/
      commands/publish_book.go                  # command + validação e transação
      queries/list_books.go                     # query de leitura
      queries/get_chapter.go                    # query de leitura
      ports/readers.go                          # portas de leitura, escrita e transação
    postgres/
      queries.go                                # SQL específico do PostgreSQL
      repository.go                             # repositório que retorna entidades de domínio
      transaction.go                            # transaction manager do seed
  integration/catalog/repository_test.go        # integração, pula sem DATABASE_URL
  internal/transport/httpapi/
    router.go                                   # NewRouter(deps Dependencies) monta chi.NewRouter() + registradores; único arquivo que importa chi
    dependencies.go                             # Dependencies, CatalogDependencies — agrupadas por domínio
    health_routes.go                            # registerHealthRoutes(httprouter.Router, health.CheckHealth)
    health_routes_test.go
    catalog_routes.go                           # registerCatalogRoutes(httprouter.Router, CatalogDependencies)
    catalog_routes_test.go

infra/
  docker-compose.yml                            # Postgres local

.github/workflows/
  api-ci.yml                                    # build, vet, test, -race, integração Postgres
```

## Etapas

### 1. Postgres local

- Criar `infra/docker-compose.yml` com um serviço `postgres:16`, porta
  `127.0.0.1:5432:5432` (não expor em todas as interfaces), credenciais de
  desenvolvimento e volume nomeado.
- Confirmar `docker compose -f infra/docker-compose.yml up -d` e
  `psql "$DATABASE_URL" -c '\dt'` funcionam contra um banco vazio.
- Commit sugerido: `chore(infra): add local postgres compose`.

### 2. Configuração do banco

- Adicionar `DatabaseURL string` a `config.Config`, lido de
  `DATABASE_URL` via `lookup`. Sem default: ausência é erro
  (`ErrDatabaseURLNotSet`), diferente de `PORT`/`ENVIRONMENT`.
- Cobrir no `config_test.go`: `DATABASE_URL` ausente falha; presente carrega.
- Commit sugerido: `feat(api): require DATABASE_URL in configuration`.

### 3. Migrations

- Instalar o CLI uma vez: `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`.
- Gerar o par de arquivos com a própria ferramenta, não na mão:

  ```sh
  migrate create -ext sql -dir internal/adapters/postgres/migrations -seq create_catalog
  ```

  Isso cria `000001_create-catalog.up.sql` e `000001_create-catalog.down.sql`.
  Preencher o `.up.sql` com constraints que protejam as invariantes do
  domínio, não só os tipos das colunas:

  ```sql
  CREATE TABLE books (
      id TEXT PRIMARY KEY,
      "order" SMALLINT NOT NULL UNIQUE CHECK ("order" > 0),
      name VARCHAR(255) NOT NULL UNIQUE CHECK (length(trim(name)) > 0),
      testament VARCHAR(3) NOT NULL CHECK (testament IN ('old', 'new')),
      chapter_count SMALLINT NOT NULL CHECK (chapter_count > 0)
  );

  CREATE TABLE verses (
      book_id TEXT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
      chapter SMALLINT NOT NULL CHECK (chapter > 0),
      number SMALLINT NOT NULL CHECK (number > 0),
      text TEXT NOT NULL CHECK (length(trim(text)) > 0),
      part SMALLINT NOT NULL CHECK (part > 0),
      PRIMARY KEY (book_id, chapter, number, part)
  );
  ```

  Sem `IF NOT EXISTS`: numa migration versionada e rastreada pelo
  `migrate`, `IF NOT EXISTS` esconderia um schema pré-existente e
  incompatível em vez de falhar visivelmente. `"order"` e `name` levam
  `UNIQUE` cada um separadamente (não uma constraint composta) — cada coluna
  precisa ser única por si só, não só o par. A FK de `verses.book_id` usa
  `ON DELETE CASCADE`: apagar um livro apaga seus versículos junto. Não há
  índice adicional em `(book_id, chapter)` — a chave primária
  `(book_id, chapter, number, part)` já cobre esse prefixo com um índice
  btree; um índice separado seria redundante.

  E o `.down.sql`:

  ```sql
  DROP TABLE IF EXISTS verses;
  DROP TABLE IF EXISTS books;
  ```

- Aplicar com o CLI (não com código Go — `cmd/api` não migra sozinho):

  ```sh
  migrate -path internal/adapters/postgres/migrations -database "$DATABASE_URL" up
  ```

- Rodar o comando duas vezes seguidas contra o Postgres local: a segunda vez
  deve responder `no change`, sem reaplicar.
- Commit sugerido: `feat(api): add catalog migrations`.

### 4. Caso de uso `ListBooks`

- Escrever `list_books_test.go` antes de `list_books.go` para exigir:

  ```go
  type Book struct {
      ID           string
      Order        int
      Name         string
      Testament    string
      ChapterCount int
  }

  type BookRepository interface {
      ListBooks(ctx context.Context) ([]Book, error)
  }

  type ListBooks struct {
      repo BookRepository
  }

  func NewListBooks(repo BookRepository) ListBooks
  func (uc ListBooks) Execute(ctx context.Context) ([]Book, error)
  ```

  `repo` fica privado — sem construtor exportado, nada fora do pacote
  `catalog` (nem os testes, nem o `cmd/api`) consegue montar a struct.

- Testar com um `BookRepository` fake (map/slice em memória), sem Postgres.
- Commit sugerido: `feat(api): add list books use case`.

### 5. Caso de uso `GetChapter`

- Escrever `get_chapter_test.go` antes de `get_chapter.go` para exigir:

  ```go
  type Verse struct {
      Number int
      Part   int
      Text   string
  }

  type Chapter struct {
      BookID   string
      BookName string
      Number   int
      Verses   []Verse
  }

  var ErrChapterNotFound = errors.New("chapter not found")

  type ChapterRepository interface {
      FindChapter(ctx context.Context, bookID string, number int) (Chapter, error)
  }

  type GetChapter struct {
      repo ChapterRepository
  }

  func NewGetChapter(repo ChapterRepository) GetChapter
  func (uc GetChapter) Execute(ctx context.Context, bookID string, number int) (Chapter, error)
  ```

- Cobrir capítulo existente, capítulo inexistente (`ErrChapterNotFound`) e
  preservação de partes repetidas (`9.1`, `9.2`) na ordem certa.
- Commit sugerido: `feat(api): add get chapter use case`.

### 6. Porta `dbexec`, adapter `PgxExecutor` e repositório do catálogo

- `internal/ports/dbexec/executor.go`: porta `Executor`/`Row`/`Rows`, sem
  nenhum tipo do pacote `pgx` (ver exemplo completo na seção Arquitetura).
- `internal/adapters/postgres/pgx_executor.go`: `PgxExecutor` implementa
  `dbexec.Executor`; `pgxConn` (privada) é satisfeita tanto por
  `*pgxpool.Pool` quanto por `pgx.Tx`, então o mesmo `PgxExecutor` serve pro
  pool (leitura) e pra transação (seed, etapa 8).
- `internal/catalog/postgres/queries.go`: texto das queries SQL.

  ```sql
  -- ListBooksQuery
  SELECT id, "order", name, testament, chapter_count FROM books ORDER BY "order"

  -- FindChapterVersesQuery — junta com books pra preencher Chapter.BookName;
  -- sem o JOIN, o contrato de Chapter não é cumprido pela persistência
  SELECT v.book_id, v.chapter, v.number, v.text, v.part, b.name
  FROM verses v
  JOIN books b ON b.id = v.book_id
  WHERE v.book_id = $1 AND v.chapter = $2
  ORDER BY v.number, v.part
  ```

- `internal/catalog/postgres/repository.go`: `Repository` depende só de
  `dbexec.Executor` (não de `pgx`) e implementa as portas de
  `application/ports` — uma struct, três responsabilidades de persistência.
  Os métodos retornam entidades de `catalog/domain`. Cada método:
  - faz `defer rows.Close()` logo após um `Query` bem-sucedido — não
    depender só do fechamento automático do `pgx` ao esgotar `Next()`,
    porque um retorno antecipado por erro de `Scan` deixaria a conexão
    presa até o GC;
  - checa o erro de cada `rows.Scan(...)`;
  - checa `rows.Err()` **depois** do loop `for rows.Next()`, nunca antes —
    checar antes do loop não verifica nada de útil, porque nenhuma linha
    foi lida ainda nesse ponto.
- `api/integration/catalog/repository_test.go`: teste de integração contra o
  Postgres local (`t.Skip` se `DATABASE_URL` não estiver setado), cobrindo
  `ListBooks`, `FindChapter` (incluindo `BookName` preenchido, ordem de
  partes repetidas e capítulo inexistente). Constrói o repositório com
  `catalogpostgres.NewRepository(postgres.NewPgxExecutor(pool))` — é o teste que
  teria pego o bug de `BookName` se existisse desde o início; testes com
  stub não pegam, porque o stub já devolve o campo certo por construção.
- Commit sugerido: `feat(api): add catalog repository`.

### 7. Rotas HTTP

`router.go` não recebe mais handlers inline conforme os domínios crescem —
cada domínio ganha seu próprio arquivo de registro de rotas, e `NewRouter`
só monta. Evita o mesmo arquivo virando um acúmulo de handlers de `/health`,
`/books`, e futuramente auth, sync, lembretes, etc.

`health_routes.go` e `catalog_routes.go` recebem `httprouter.Router`
(`internal/ports/httprouter`), não `chi.Router` — a porta tem só
`Get(pattern string, handler http.HandlerFunc)`, que `chi.Router` já
satisfaz estruturalmente, sem adapter nenhum no meio. Só `router.go`
importa `chi` (pra construir o `chi.NewRouter()` concreto). Dentro dos
handlers, path params usam `r.PathValue("bookId")` — stdlib desde Go 1.22,
já populado automaticamente pelo `chi` — no lugar de `chi.URLParam`, então
nenhum dos dois arquivos de rota importa `chi`.

`NewRouter` também não recebe um parâmetro por caso de uso — isso escalaria
com o número de casos de uso (ilimitado), não com o número de domínios
(limitado). As dependências ficam agrupadas por domínio em
`internal/transport/httpapi/dependencies.go`:

```go
type Dependencies struct {
    Health  health.CheckHealth
    Catalog CatalogDependencies
}

type CatalogDependencies struct {
    ListBooks  *catalog.ListBooks
    GetChapter *catalog.GetChapter
}
```

Um caso de uso novo num domínio existente (ex.: `catalog` ganhar
`SearchVerses`) só mexe em `CatalogDependencies`, não em `NewRouter`. Um
domínio novo soma um campo em `Dependencies` e uma chamada a mais em
`NewRouter` — nunca mais que isso.

- Extrair o que já existe pra `health_routes.go`.
- Escrever `catalog_routes_test.go` antes de `catalog_routes.go`, exigindo:

  ```text
  GET /books
  status: 200
  body: array de livros, ordenados por order

  GET /books/{bookId}/chapters/{number}
  status: 200
  body: capítulo com verses em ordem, book_name preenchido

  GET /books/xx/chapters/1        -> 404 (livro inexistente)
  GET /books/gn/chapters/9999     -> 404 (capítulo inexistente)
  GET /books/gn/chapters/0        -> 400 (number inválido)
  GET /books/gn/chapters/-1       -> 400
  GET /books/gn/chapters/abc      -> 400
  ```

- Implementar `catalog_routes.go`: sem lógica de negócio no handler — só
  tradução de path params, validação de forma (`number` precisa ser inteiro
  `> 0`, senão `400` antes de chamar o caso de uso), chamada ao caso de uso
  e status HTTP (`catalog.ErrChapterNotFound` → `404`; qualquer outro erro →
  `500`).
- Reescrever `router.go` como composição pura.
- Commit sugerido: `refactor(api): split router registration by domain` (a
  extração de `/health`), seguido de `feat(api): expose catalog endpoints`
  (as rotas novas).

### 8. Seed do corpus

O seed não é upsert por linha — é **substituir o livro inteiro dentro de uma
transação**. Upsert por si só não publica uma projeção fiel: se uma revisão
do corpus remover ou renumerar um versículo, o upsert não apaga a linha
antiga, ela fica órfã no banco pra sempre. Substituir o livro inteiro
(apagar tudo daquele `book_id` e reinserir do zero, na mesma transação)
garante que o banco sempre reflete exatamente o corpus publicado, livro por
livro.

#### 8.1. Migration da versão publicada

Nova migration (`migrate create ... -seq catalog_version`), registrando qual
hash do `manifest.json` está publicado — sem isso não dá pra auditar se o
banco está desatualizado em relação ao corpus:

```sql
-- 000002_catalog-version.up.sql
CREATE TABLE catalog_version (
    id            SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    corpus_sha256 TEXT NOT NULL CHECK (length(trim(corpus_sha256)) > 0),
    published_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

`CHECK (id = 1)` força linha única — a tabela guarda só a versão *atual*,
não histórico. `corpus_sha256` vem de `bibleSha256` em
`bible/corpus/v1/manifest.json`.

```sql
-- 000002_catalog-version.down.sql
DROP TABLE IF EXISTS catalog_version;
```

#### 8.2. Escrita no repositório

> [!note] Implementado com uma diferença do desenho original
> `BookInput` e `VerseInput` existem no command como contrato de entrada; o
> repositório e a porta continuam recebendo `domain.Book`/`domain.Verse`, não
> DTOs HTTP. E
> o loop de `Exec` por versículo foi trocado por `dbexec.Executor.CopyFrom`
> (streaming via `COPY` do Postgres) antes mesmo de existir carga real pra
> medir — motivo real: rodar o seed contra um banco externo com N
> milhares de `Exec` (um por versículo) paga N round-trips de rede; `COPY`
> paga um. Ver [[Rules/01 - Princípios de Engenharia]] pra `CopyFrom` como
> parte da porta `dbexec.Executor`.

`internal/catalog/postgres/repository.go` ganha um método de escrita, e o
command (`internal/catalog/application/commands/publish_book.go`) define as
portas que ele
implementa — mesmo padrão de sempre:

```go
// internal/catalog/application/commands/publish_book.go
type PublishBook struct { transactions ports.TransactionManager }

func (uc PublishBook) Execute(ctx context.Context, input PublishBookInput) error {
    book, err := domain.NewBook(/* dados do input */)
    if err != nil { return err }
    // Constrói os Verse, valida e chama ReplaceBook dentro da transação.
    return uc.transactions.WithinTransaction(ctx, func(ctx context.Context, writer ports.CatalogWriter) error {
        return writer.ReplaceBook(ctx, book, verses)
    })
}
```

`Repository.ReplaceBook` apaga os versículos do livro antes de mexer nele
(evita problema de FK), faz upsert do livro, e carrega os versículos novos
via `CopyFrom` — sem loop, sem `Exec` por linha:

```go
func (r *Repository) ReplaceBook(ctx context.Context, book Book, verses []Verse) error {
    if err := r.db.Exec(ctx, DeleteBookVersesQuery, book.ID); err != nil {
        return err
    }
    if err := r.db.Exec(ctx, UpsertBookQuery, book.ID, book.Order, book.Name, book.Testament, book.ChapterCount); err != nil {
        return err
    }

    rows := make([][]any, len(verses))
    for i, verse := range verses {
        rows[i] = []any{book.ID, verse.Chapter, verse.Number, verse.Text, verse.Part}
    }

    _, err := r.db.CopyFrom(ctx, "verses", verseColumns, rows)
    return err
}
```

Isso exigiu estender a porta `dbexec.Executor` com um quarto método, sem
vazar tipo de `pgx` na assinatura (`table string`, `columns []string`,
`rows [][]any` — só tipos nossos):

```go
type Executor interface {
    QueryRow(ctx context.Context, sql string, args ...any) Row
    Query(ctx context.Context, sql string, args ...any) (Rows, error)
    Exec(ctx context.Context, sql string, args ...any) error
    CopyFrom(ctx context.Context, table string, columns []string, rows [][]any) (int64, error)
}
```

`PgxExecutor.CopyFrom` traduz pra `pgx.CopyFrom(ctx, pgx.Identifier{table}, columns, pgx.CopyFromRows(rows))`.

Queries em `query.go` (`InsertVerseQuery` não existe mais — substituída pelo
`CopyFrom`):

```sql
-- DeleteBookVersesQuery
DELETE FROM verses WHERE book_id = $1

-- UpsertBookQuery
INSERT INTO books (id, "order", name, testament, chapter_count)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
    "order" = EXCLUDED."order",
    name = EXCLUDED.name,
    testament = EXCLUDED.testament,
    chapter_count = EXCLUDED.chapter_count
```

Medido de verdade: corpus inteiro (73 livros, 35.624 versículos) publicado
em ~1,1s contra Postgres local, idempotente (rodar duas vezes seguidas
mantém as mesmas contagens).

#### 8.3. `cmd/seed/catalog/main.go`

> [!note] Estado implementado
> O comando lê `bible/corpus/v1/bible.json` e chama `PublishBook.Execute` por
> livro. O command cria e valida o domínio, e o
> `catalog/postgres.TransactionManager` abre, confirma ou reverte uma
> transação por livro. A escrita usa `CopyFrom`, portanto um livro não fica
> parcialmente publicado se a operação falhar.
>
> O seed grava `catalog_version` somente depois que todos os livros são
> publicados. Qualquer falha interrompe o processo e produz status de saída
> diferente de zero; a versão não é registrada nesse caso.

- Não depende de `bible/tools` (módulo Go separado) — decodifica o JSON já
  validado do corpus v1 diretamente, sem reimportar as regras de
  normalização.
- Idempotente por natureza pro que já está implementado: rodar duas vezes
  seguidas com o mesmo corpus produz o mesmo estado, porque cada livro é
  substituído por inteiro, não acumulado — confirmado rodando duas vezes
  contra o Postgres local (`SELECT count(*) FROM books` = 73, `SELECT
  count(*) FROM verses` = 35624 nas duas vezes).
- `SELECT corpus_sha256 FROM catalog_version` recebe o `bibleSha256` do
  manifesto após um seed completo; a operação usa upsert para suportar novas
  versões do corpus.

### 9. Composition root

- Atualizar `cmd/api/main.go`:
  - abrir o pool com `pgxpool.New(ctx, cfg.DatabaseURL)`;
  - chamar `pool.Ping(ctx)` com timeout antes de aceitar tráfego — falha de
    conexão vira `slog.Error` + `os.Exit(1)`, não erro silencioso na
    primeira requisição;
  - montar `catalog/postgres.Repository` com `postgres.NewPgxExecutor(pool)`,
    `queries.NewListBooks(repo)`, `queries.NewGetChapter(repo)` e
    `httpapi.NewRouter(httpapi.Dependencies{...})`;
  - `cmd/api` não aplica migration — assume que o schema já foi migrado via
    CLI (etapa 3);
  - registrar `signal.Notify` pra `SIGINT`/`SIGTERM` e chamar
    `srv.Shutdown(ctx)` com timeout — drena requisições em curso em vez de
    derrubar conexões abruptamente; `defer pool.Close()` fecha o pool
    depois do shutdown.
- Rodar `migrate ... up` (se ainda não tiver rodado), depois `go run
  ./cmd/api`, subir o Postgres local e verificar
  `curl http://localhost:8080/books` e
  `curl http://localhost:8080/books/gn/chapters/1`. Confirmar que
  `SIGTERM` loga "shutting down" e o processo sai com código `0`.
- Commit sugerido: `feat(api): wire catalog into the server entrypoint`.

### 10. CI da API

- `.github/workflows/api-ci.yml`, no molde de `vault-lint.yml` e
  `bible-integrity.yml`: um job `build-and-test` (`go build`, `go vet`,
  `go test`, `go test -race`, sem Postgres) e um job `integration` que sobe
  `postgres:16` como serviço do GitHub Actions, aplica as migrations via
  CLI `migrate` e roda `go test ./...` de novo com `DATABASE_URL` setado —
  cobrindo `repository_test.go`, que só roda com banco disponível.
- O workflow cobre build, vet, testes unitários, race detector e integração
  contra PostgreSQL.

## Fora deste plano

- Busca textual, referências cruzadas, múltiplas traduções.
- Cache (Redis) das respostas do catálogo.
- Autenticação, progresso, sincronização e compartilhamento.
- `/readyz` separado de `/health`, métricas, tracing, SLOs — necessários
  antes de produção com múltiplas réplicas, mas fora do escopo desta
  entrega (ver `AVALIACAO_ARQUITETURAL.md`, achado 9).
- Licenciamento e proveniência do corpus bíblico distribuído — risco
  jurídico real, não uma decisão técnica; precisa de decisão do dono do
  produto antes de qualquer publicação (ver `AVALIACAO_ARQUITETURAL.md`,
  achado 12).
- Deploy e qualquer client web ou mobile consumindo esses endpoints.

Essas partes dependem de decisões e contratos que merecem entregas próprias.

---

◀ [[Plans/Active/_Index|Planos Ativos]] · próxima: [[Plans/Archive/_Index|Arquivo]] ▶
