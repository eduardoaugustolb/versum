---
title: "Visão Geral da Arquitetura"
section: Docs
subsection: Architecture
type: architecture
status: approved
tags: [versum, docs, architecture]
up: "[[Docs/Architecture/_Index|Arquitetura]]"
prev: "[[Docs/Architecture/_Index|Arquitetura]]"
next: "[[Docs/Architecture/Autenticação e Sessões]]"
related: ["[[Docs/Apps/_Index]]", "[[Docs/Decisions/001 - Plataforma Inicial]]"]
---

# Visão Geral da Arquitetura

🏠 [[_Index|Home]] › 📚 [[Docs/_Index|Documentação]] › 📐 [[Docs/Architecture/_Index|Arquitetura]] › **Visão geral**

> [!info] Contexto
> O sistema serve leitura online no web e leitura offline no Android, mantendo o
> servidor como referência do estado sincronizado.

## Estrutura do repositório

```text
versum/
  api/             API Go, worker e migrations
  web/             leitor online em Next.js
  mobile/          aplicativo Android em React Native/Expo
  infra/           ambiente local e deploy
  Obsidian Vault/  documentação humana
```

Os clientes são projetos independentes. Não existe workspace JavaScript nem um
pacote de regras compartilhado entre TypeScript e Go. A API publica o contrato
consumido pelos clientes.

## Backend

A API usa Go, `net/http`, `chi`, PostgreSQL via `pgx`, Redis e S3 compatível.
Ela segue Clean Architecture e Hexagonal Architecture (Ports & Adapters),
organizadas por funcionalidade:

- casos de uso definem portas pequenas;
- adapters HTTP traduzem requests e respostas;
- adapters externos implementam Postgres, Redis, S3, e-mail, FCM e Discord;
- `cmd/api` e `cmd/worker` montam as dependências.

Casos de uso não conhecem bibliotecas HTTP, banco, cache, storage ou push.
PostgreSQL é a fonte de verdade. Redis serve somente para cache, rate limit e
locks curtos. S3 guarda imagens de compartilhamento.

### Estrutura interna da API

```text
api/
  cmd/api/                    composition root da API HTTP
  cmd/worker/                 composition root do worker assíncrono
  internal/<domínio>/         casos de uso + portas, por funcionalidade
    <ação>.go                 um caso de uso por arquivo (ex.: check.go)
    <ação>_test.go
  internal/adapters/<tecnologia>/   implementações de portas (postgres, redis, s3, email, fcm, discord)
  internal/transport/httpapi/       adapter HTTP (handlers, router)
```

`internal/adapters/<tecnologia>` é agrupado por tecnologia, não por domínio:
adapters são implementação plugável e compartilhada entre domínios. Por
exemplo, `internal/adapters/postgres/` pode conter `health_repository.go`,
`auth_repository.go` etc., cada um implementando a porta do respectivo
domínio. O padrão de caso de uso, porta e adapter — com exemplo de código —
está detalhado em [[Rules/01 - Princípios de Engenharia]].

## Responsabilidades

| Componente | Responsabilidade |
| :-- | :-- |
| API | regras de negócio, autenticação, sincronização e contrato público |
| Worker | jobs de lembrete, geração assíncrona e manutenção |
| Web | leitura online e links compartilháveis |
| Mobile | leitura offline, outbox e notificações locais |

---

◀ [[Docs/Architecture/_Index|Arquitetura]] · próxima: [[Docs/Architecture/Autenticação e Sessões|Autenticação e sessões]] ▶
