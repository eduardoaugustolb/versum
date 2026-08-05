---
title: "02 - Catálogo Bíblico"
section: Plans
subsection: Active
type: implementation-plan
status: approved
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

O domínio `catalog` expõe dois casos de uso pequenos: `ListBooks` e
`GetChapter`. Cada um define sua própria porta de repositório. Um único
`Repository`, dentro do próprio pacote `catalog`, implementa as duas portas,
consultando tabelas `books` e `verses`.

> [!note] Histórico da decisão
> Essa camada foi implementada, removida e reimplementada nesta mesma
> entrega. Primeiro existia como `sql.Executor`/`PgxExecutor`
> (`internal/adapters/sql`), justificada por portabilidade de banco. Uma
> avaliação arquitetural externa (`AVALIACAO_ARQUITETURAL.md`, 5 de agosto
> de 2026) apontou que essa portabilidade era ilusória — o texto das
> queries (`$1`, `ON CONFLICT`, `"order"`) já é específico de Postgres — e
> a camada foi removida, com `catalog.Repository` passando a depender de
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
// internal/catalog/repository.go
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
  cmd/seed-catalog/main.go
  internal/config/config.go                    # + DatabaseURL
  internal/ports/dbexec/
    executor.go                                 # porta Executor/Row/Rows, sem tipo de pgx
  internal/ports/httprouter/
    router.go                                   # porta Router, sem tipo de chi
  internal/adapters/postgres/
    migrations/
      000001_create-catalog.up.sql
      000001_create-catalog.down.sql
    pgx_executor.go                              # PgxExecutor implementa dbexec.Executor usando pgx
  internal/catalog/
    book.go                                     # tipo Book
    chapter.go                                  # tipos Chapter, Verse
    list_books.go                               # caso de uso + porta BookRepository
    list_books_test.go
    get_chapter.go                              # caso de uso + porta ChapterRepository
    get_chapter_test.go
    query.go                                    # texto das queries SQL
    repository.go                               # Repository — depende de dbexec.Executor
    repository_test.go                          # integração, pula sem DATABASE_URL
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
- `internal/catalog/query.go`: texto das duas queries.

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

- `internal/catalog/repository.go`: `Repository` depende só de
  `dbexec.Executor` (não de `pgx`) e implementa `BookRepository` e
  `ChapterRepository` — uma struct, dois métodos. Cada método:
  - faz `defer rows.Close()` logo após um `Query` bem-sucedido — não
    depender só do fechamento automático do `pgx` ao esgotar `Next()`,
    porque um retorno antecipado por erro de `Scan` deixaria a conexão
    presa até o GC;
  - checa o erro de cada `rows.Scan(...)`;
  - checa `rows.Err()` **depois** do loop `for rows.Next()`, nunca antes —
    checar antes do loop não verifica nada de útil, porque nenhuma linha
    foi lida ainda nesse ponto.
- `internal/catalog/repository_test.go`: teste de integração contra o
  Postgres local (`t.Skip` se `DATABASE_URL` não estiver setado), cobrindo
  `ListBooks`, `FindChapter` (incluindo `BookName` preenchido, ordem de
  partes repetidas e capítulo inexistente). Constrói o repositório com
  `catalog.NewRepository(postgres.NewPgxExecutor(pool))` — é o teste que
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

- `cmd/seed-catalog/main.go`: lê `bible/corpus/v1/books/*.json` (caminho via
  flag ou `BIBLE_CORPUS_PATH`, default `bible/corpus/v1`), decodifica cada
  livro e faz upsert (`INSERT ... ON CONFLICT DO UPDATE`) em `books` e
  `verses` dentro de uma transação por livro — abre a `pgx.Tx` e monta
  `catalog.NewRepository(postgres.NewPgxExecutor(tx))` por livro,
  reaproveitando o mesmo repositório da leitura (`pgxConn` aceita tanto
  pool quanto transação).
- Upsert por livro não é suficiente pra publicar uma projeção fiel: se um
  versículo for removido ou renumerado numa revisão do corpus, o upsert não
  o remove do banco — ele fica órfão. Antes de publicar, o seed deve
  truncar `verses`/`books` do livro dentro da mesma transação e reinserir
  do zero, em vez de só fazer upsert. Registrar `bibleSha256` do
  `manifest.json` numa tabela `catalog_version` (nova migration) pra
  auditoria de qual versão do corpus está publicada.
- Não depende de `bible/tools` (módulo Go separado) — decodifica o JSON já
  validado do corpus v1 diretamente, sem reimportar as regras de
  normalização.
- Rodar duas vezes seguidas contra o Postgres local e confirmar que
  `SELECT count(*) FROM books` continua em 73 e `SELECT count(*) FROM
  verses` bate com `manifest.json`.
- Commit sugerido: `feat(api): add catalog seed command`.

### 9. Composition root

- Atualizar `cmd/api/main.go`:
  - abrir o pool com `pgxpool.New(ctx, cfg.DatabaseURL)`;
  - chamar `pool.Ping(ctx)` com timeout antes de aceitar tráfego — falha de
    conexão vira `slog.Error` + `os.Exit(1)`, não erro silencioso na
    primeira requisição;
  - montar `catalog.NewRepository(postgres.NewPgxExecutor(pool))`,
    `catalog.NewListBooks(repo)`, `catalog.NewGetChapter(repo)` e
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
- Sem isso, uma incompatibilidade como a do `cmd/api`/`httpapi.Dependencies`
  fica invisível até alguém rodar localmente — hooks de pre-commit locais
  não substituem CI remoto num PR.
- Commit sugerido: `ci(api): add build, vet, test and integration workflow`.

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
