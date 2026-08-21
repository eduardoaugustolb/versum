---
title: "01 - Princípios de Engenharia"
section: Rules
type: rule
status: active
tags: [versum, rules, engineering, architecture]
up: "[[Rules/_Index|Regras]]"
prev: "[[Rules/_Index|Regras]]"
next: "[[Rules/02 - Segurança]]"
---

# 01 — Princípios de Engenharia

## Regra

- Organizar código por funcionalidade, não por camada global.
- Definir portas no caso de uso que delas depende.
- Cada caso de uso resolve uma única operação de negócio: um método público
  (`Execute`, ou nome equivalente) por struct, sem acumular ações não
  relacionadas.
- Manter adapters externos fora do domínio.
- O repositório conhece a query/SQL que executa, mas não conhece o caso de uso
  que o chama: a dependência aponta do caso de uso para a porta, nunca o
  contrário. Isso permite múltiplos adapters (Postgres, mock de teste etc.)
  para a mesma porta sem que o caso de uso saiba que SQL existe.
- Isolar qualquer driver ou framework concreto atrás de uma porta pequena e
  neutra — sem nenhum tipo dele na própria assinatura da porta, não só um
  wrapper fino. Isso vale para banco (`dbexec.Executor`, sem tipo de `pgx`)
  e para o roteador HTTP (`httprouter.Router`, sem tipo de `chi`) do mesmo
  jeito: quem usa a porta (`catalog/application`, `catalog_routes.go`) nunca
  importa o pacote do driver; só o adapter concreto
  (`postgres.PgxExecutor`, `router.go`) importa. Trocar de tecnologia vira
  mudar o que é injetado na composição, não reescrever quem usa a porta —
  mesmo quando a troca nunca chega a acontecer de verdade, o isolamento
  ainda evita que o driver vaze pra código que deveria ser só regra de
  domínio ou tradução HTTP.
- Tratar PostgreSQL como fonte de verdade e cache como dado descartável.
- Criar testes unitários para regras e integração para banco, contrato e
  concorrência.
- Preferir a menor abstração que mantém uma dependência substituível.

## Por quê

O projeto deve ser fácil de explicar, testar e evoluir. Seguimos Clean
Architecture (regra de dependência: as camadas internas não conhecem as
externas) e Hexagonal Architecture / Ports & Adapters (a porta é definida pelo
caso de uso, adapters plugam nela). As duas são um meio para limites claros,
não uma meta de quantidade de interfaces ou camadas.

## Exemplo

O caso de uso define a porta (`Repository`) e depende só dela. O adapter
Postgres implementa a porta e é o único lugar que conhece SQL:

```go
// internal/health/check.go — caso de uso define a porta
package health

type Repository interface {
    FindStatus(ctx context.Context) (Status, error)
}

type CheckHealth struct {
    repo Repository
}

func (uc CheckHealth) Execute(ctx context.Context) (Status, error) {
    return uc.repo.FindStatus(ctx)
}
```

```go
// internal/adapters/postgres/health_repository.go — adapter implementa a porta
package postgres

type HealthRepository struct {
    db *pgxpool.Pool
}

func (r HealthRepository) FindStatus(ctx context.Context) (health.Status, error) {
    var state string
    err := r.db.QueryRow(ctx, "SELECT state FROM health_check LIMIT 1").Scan(&state)
    if err != nil {
        return health.Status{}, err
    }
    return health.Status{State: state}, nil
}
```

A cadeia de conhecimento vai numa direção só: `postgres.HealthRepository` conhece
SQL e o pool de conexão; `health.CheckHealth` conhece só a interface
`Repository`. Nenhum dos dois conhece o outro lado — o composition root em
`cmd/api` é quem liga `postgres.HealthRepository{}` a `health.CheckHealth{repo: ...}`.

## Isolar driver/framework atrás de uma porta

O mesmo padrão se aplica a qualquer biblioteca externa que o código realmente
aciona, não só bancos. Duas instâncias reais no `api/`:

```go
// internal/ports/dbexec/executor.go — porta neutra, sem tipo de pgx
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
// internal/ports/httprouter/router.go — porta neutra, sem tipo de chi
package httprouter

type Router interface {
    Get(pattern string, handler http.HandlerFunc)
}
```

`internal/adapters/postgres.PgxExecutor` implementa `dbexec.Executor` e é o
único lugar que importa `pgx`. `internal/transport/httpapi/router.go`
constrói o `chi.NewRouter()` concreto e é o único lugar que importa `chi` —
`health_routes.go` e `catalog_routes.go` recebem `httprouter.Router` e usam
só `net/http` puro (`r.PathValue(...)` em vez de `chi.URLParam`, disponível
desde Go 1.22 e populado pelo `chi` automaticamente).

No catálogo, leituras e escritas também são separadas em `application/queries`
e `application/commands`. Isso é uma aplicação pragmática de CQRS: queries não
alteram estado; commands validam e alteram estado. Não implica dois bancos nem
dois serviços.

Nos dois casos, a pergunta que decide se vale a porta não é "vamos trocar
essa tecnologia algum dia" — normalmente a resposta é não. É "esse código
(`catalog/application`, `catalog_routes.go`) deveria conhecer só a regra que
resolve, ou também os detalhes de uma biblioteca de terceiro que não tem
nada a ver com o que ele decide?". Quando a porta cabe em poucos métodos e o
adapter é fino, o isolamento se paga mesmo sem uma segunda implementação
concreta aparecer.

---

◀ [[Rules/_Index|Regras]] · próxima: [[Rules/02 - Segurança|Segurança]] ▶
