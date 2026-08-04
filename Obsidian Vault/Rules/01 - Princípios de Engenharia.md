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

---

◀ [[Rules/_Index|Regras]] · próxima: [[Rules/02 - Segurança|Segurança]] ▶
