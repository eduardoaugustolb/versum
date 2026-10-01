---
title: "Sincronização Offline"
section: Docs
subsection: Architecture
type: architecture
status: approved
tags: [versum, docs, architecture, offline, sync]
up: "[[Docs/Architecture/_Index|Arquitetura]]"
prev: "[[Docs/Architecture/Autenticação e Sessões]]"
next: "[[Docs/Architecture/Privacidade e Consentimento]]"
related: ["[[Docs/Decisions/002 - Progresso por Eventos]]", "[[Rules/01 - Princípios de Engenharia]]"]
---

# Sincronização Offline

🏠 [[_Index|Home]] › 📚 [[Docs/_Index|Documentação]] › 📐 [[Docs/Architecture/_Index|Arquitetura]] › **Sincronização offline**

O desenho futuro do Android usa SQLite para conteúdo baixado, estado de leitura e uma outbox de eventos pendentes. Cada evento contém ID único, dispositivo, sequência local,
referência de leitura e instante de criação.

Ao sincronizar, a API processa os eventos numa transação. Uma restrição de
unicidade torna reenvios idempotentes. Eventos são aditivos: um dispositivo
atrasado acrescenta histórico, mas não reduz o ponto confirmado mais avançado.
A API responde com o estado canônico para reconciliar a cópia local.

O app tenta sincronizar após uma ação, ao recuperar conectividade e ao voltar ao
primeiro plano. O sistema operacional pode suspender o processo; por isso não
há promessa de sincronização contínua e uma desinstalação pode perder eventos
ainda não enviados.

## Ciclo de vida e privacidade

Este módulo não está implementado na API atual. Histórico identificável pode
revelar convicção religiosa: avaliar a hipótese do art. 11 antes de coletar.
Justificar a necessidade de cada evento e sua duração; histórico aditivo não
autoriza retenção ilimitada (LGPD arts. 6º, III, e 15–16).

Revogação e exclusão precisam alcançar SQLite, outbox local, API e projeções.
Eventos antigos de conta excluída não podem recriar dados. Definir tombstones
ou outro mecanismo com retenção mínima justificada, rejeição de eventos e
restauração de backup que respeite exclusões. Consentimento, quando aplicável,
deve ser validado no servidor antes de aceitar novos eventos (arts. 8º, §5º,
18 e 46).

Downloads públicos e progresso pessoal devem ser separados no mapa; backup do
SO e cópias locais precisam ser avaliados no plano Android. Ver
[[Plans/Active/04 - Privacidade e Consentimento|Plano 04]].

---

◀ [[Docs/Architecture/Autenticação e Sessões|Autenticação e sessões]] · próxima: [[Docs/Architecture/Privacidade e Consentimento|Privacidade e consentimento]] ▶
