---
title: "002 - Progresso por Eventos"
section: Docs
subsection: Decisions
type: adr
status: accepted
date: 2026-08-02
tags: [versum, docs, adr, sync, offline]
up: "[[Docs/Decisions/_Index|Decisões]]"
prev: "[[Docs/Decisions/001 - Plataforma Inicial]]"
next: "[[Rules/_Index]]"
related: ["[[Docs/Architecture/Sincronização Offline]]"]
---

# 002 — Progresso por Eventos

## Contexto

Dois dispositivos podem ler sem conexão e sincronizar em ordens diferentes. Uma
atualização que simplesmente substitui o progresso permite duplicação e pode
fazer o estado regredir.

## Decisão

Representar avanços como eventos idempotentes. O banco aceita cada `eventId`
uma única vez e mantém uma projeção canônica monotônica: o ponto confirmado mais
avançado nunca diminui.

## Consequências

- Retries são seguros.
- Eventos atrasados não removem progresso.
- O servidor continua sendo a referência entre dispositivos.
- Há mais dados e uma projeção explícita para manter e observar.

## Alternativas descartadas

- Last write wins: pode apagar avanço válido de outro dispositivo.
- Sincronizar apenas um cursor: não preserva histórico nem torna retries
  idempotentes.
- Confiar no cliente para o estado final: não oferece proteção contra conflito
  ou manipulação em massa.

---

◀ [[Docs/Decisions/001 - Plataforma Inicial|ADR 001]] · próxima: [[Rules/_Index|Regras]] ▶
