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
  internal/<domínio>/         casos de uso, portas e o repositório SQL do domínio
    <ação>.go                 um caso de uso por arquivo (ex.: check.go)
    <ação>_test.go
    repository.go             repositório SQL — depende de dbexec.Executor, não de pgx
  internal/ports/dbexec/         porta de execução SQL, neutra entre drivers (sem tipo de pgx)
  internal/ports/httprouter/     porta de registro de rota, neutra entre roteadores (sem tipo de chi)
  internal/adapters/<tecnologia>/   driver concreto que implementa a porta (postgres, redis, s3, email, fcm, discord)
  internal/transport/httpapi/       adapter HTTP — router.go monta o chi concreto; os demais arquivos usam só httprouter.Router e net/http
```

Todo driver ou framework de terceiro que o código realmente aciona — banco,
roteador HTTP, fila, o que for — fica atrás de uma porta pequena e neutra:
sem nenhum tipo do driver na própria assinatura da porta, só no adapter que
a implementa. `catalog.Repository` depende de `dbexec.Executor`
(implementado por `postgres.PgxExecutor`, que conhece `pgx`);
`catalog_routes.go`/`health_routes.go` dependem de `httprouter.Router`
(satisfeita diretamente por `chi.Router`, sem adapter — só `router.go`
conhece `chi`). Casos de uso não conhecem SQL nem `pgx`; o handler HTTP não
conhece regra de negócio nem `chi`.

O critério não é "vamos trocar essa tecnologia algum dia" — normalmente a
resposta é não, e o texto das queries continua específico de Postgres de
qualquer forma. É se o código que resolve a regra de negócio ou traduz
HTTP deveria também conhecer os detalhes de uma biblioteca de terceiro que
não tem nada a ver com o que ele decide. Ver
[[Plans/Active/02 - Catálogo Bíblico]] pro histórico dessa decisão (foi
revertida e revisada duas vezes antes de chegar nesse formato).

`internal/adapters/<tecnologia>/` guarda o que exige infraestrutura própria
sem uma porta neutra equivalente ainda definida: migrations SQL, Redis, S3,
e-mail, FCM, Discord.

O padrão de caso de uso e porta — com exemplo de código — está detalhado em
[[Rules/01 - Princípios de Engenharia]].

## Responsabilidades

| Componente | Responsabilidade |
| :-- | :-- |
| API | regras de negócio, autenticação, sincronização e contrato público |
| Worker | jobs de lembrete, geração assíncrona e manutenção |
| Web | leitura online e links compartilháveis |
| Mobile | leitura offline, outbox e notificações locais |

---

◀ [[Docs/Architecture/_Index|Arquitetura]] · próxima: [[Docs/Architecture/Autenticação e Sessões|Autenticação e sessões]] ▶
