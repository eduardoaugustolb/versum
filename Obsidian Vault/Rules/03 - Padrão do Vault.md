---
title: "03 - Padrão do Vault"
section: Rules
type: rule
status: active
tags: [versum, rules, docs, vault]
up: "[[Rules/_Index|Regras]]"
prev: "[[Rules/02 - Segurança]]"
next: "[[Plans/_Index]]"
---

# 03 — Padrão do Vault

## Estrutura

- `_Index.md` é um mapa de navegação, nunca uma duplicação dos documentos.
- `PRD.md` descreve o produto.
- `Docs/` guarda conhecimento estável; `Docs/Decisions/` guarda ADRs.
- `Rules/` guarda padrões obrigatórios e duradouros.
- `Plans/Active/` guarda trabalho em andamento; planos concluídos vão para
  `Plans/Archive/`.

## Metadados e links

Toda nota adicionada ao vault possui frontmatter com `title`, `section`, `type`,
`status`, `tags` e `up`. Notas em uma trilha possuem `prev` e `next` e mostram
os mesmos links no rodapé. Use wikilinks para notas internas e links Markdown
somente para URLs externas.

## Critério de criação

Criar uma nota quando ela responde uma pergunta que continuará útil depois da
entrega atual. Detalhes de PR, experimento e rascunho não entram no vault.

## Modelos

Use os modelos em `.obsidian/templates/` para notas de decisão, plano e guia.

---

◀ [[Rules/02 - Segurança|Segurança]] · próxima: [[Plans/_Index|Planos]] ▶
