---
title: "Sincronização Offline"
section: Docs
subsection: Architecture
type: architecture
status: approved
tags: [versum, docs, architecture, offline, sync]
up: "[[Docs/Architecture/_Index|Arquitetura]]"
prev: "[[Docs/Architecture/Autenticação e Sessões]]"
next: "[[Rules/_Index]]"
related: ["[[Docs/Decisions/002 - Progresso por Eventos]]", "[[Rules/01 - Princípios de Engenharia]]"]
---

# Sincronização Offline

🏠 [[_Index|Home]] › 📚 [[Docs/_Index|Documentação]] › 📐 [[Docs/Architecture/_Index|Arquitetura]] › **Sincronização offline**

O Android usa SQLite para conteúdo baixado, estado de leitura e uma outbox de
eventos pendentes. Cada evento contém ID único, dispositivo, sequência local,
referência de leitura e instante de criação.

Ao sincronizar, a API processa os eventos numa transação. Uma restrição de
unicidade torna reenvios idempotentes. Eventos são aditivos: um dispositivo
atrasado acrescenta histórico, mas não reduz o ponto confirmado mais avançado.
A API responde com o estado canônico para reconciliar a cópia local.

O app tenta sincronizar após uma ação, ao recuperar conectividade e ao voltar ao
primeiro plano. O sistema operacional pode suspender o processo; por isso não
há promessa de sincronização contínua e uma desinstalação pode perder eventos
ainda não enviados.

---

◀ [[Docs/Architecture/Autenticação e Sessões|Autenticação e sessões]] · próxima: [[Rules/_Index|Regras]] ▶
