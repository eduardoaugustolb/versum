---
title: "001 - Plataforma Inicial"
section: Docs
subsection: Decisions
type: adr
status: accepted
date: 2026-08-02
tags: [versum, docs, adr, platform]
up: "[[Docs/Decisions/_Index|Decisões]]"
prev: "[[Docs/Decisions/_Index|Decisões]]"
next: "[[Docs/Decisions/002 - Progresso por Eventos]]"
related: ["[[Docs/Architecture/Visão Geral]]"]
---

# 001 — Plataforma Inicial

## Contexto

O produto exige leitura online no web e leitura offline confiável no celular. A
equipe conhece React e Next.js, está aprendendo Go e ainda não tem experiência
com desenvolvimento nativo.

## Decisão

Usar Go para API e worker, Next.js para o web e React Native/Expo para Android.
Os três projetos ficam no mesmo repositório, mas sem workspace JavaScript. A
API é a fronteira entre os clientes e as regras de domínio.

## Consequências

- TypeScript reduz a curva de aprendizado do aplicativo mobile.
- Android vem antes de iOS para manter o primeiro lançamento controlado.
- Go concentra regras de negócio, dados e infraestrutura.
- Não há biblioteca de domínio compartilhada entre os clientes e a API.

## Alternativas descartadas agora

- Flutter: excelente para offline-first, mas exigiria aprender Dart e Flutter.
- PWA como único cliente: não entrega o mesmo controle de armazenamento e push
  de um aplicativo Android.
- iOS desde o início: aumenta escopo, contas e configuração sem validar o
  produto primeiro.

---

◀ [[Docs/Decisions/_Index|Decisões]] · próxima: [[Docs/Decisions/002 - Progresso por Eventos|ADR 002]] ▶
