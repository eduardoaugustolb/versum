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
- Tratar PostgreSQL como fonte de verdade e cache como dado descartável.
- Criar testes unitários para regras e integração para banco, contrato e
  concorrência.
- Preferir a menor abstração que mantém uma dependência substituível.

## Por quê

O projeto deve ser fácil de explicar, testar e evoluir. Clean Architecture é um
meio para limites claros, não uma meta de quantidade de interfaces.

---

◀ [[Rules/_Index|Regras]] · próxima: [[Rules/02 - Segurança|Segurança]] ▶
