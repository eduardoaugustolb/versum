---
title: "02 - Corpus Bíblico Canônico"
section: Plans
subsection: Archive
type: implementation-plan
status: completed
date: 2026-08-02
tags: [versum, plans, bible, go, data]
up: "[[Plans/Archive/_Index|Arquivo de Planos]]"
prev: "[[Plans/Archive/01 - Fundação da API Go]]"
next: "[[Plans/Archive/_Index|Arquivo]]"
related: ["[[PRD]]", "[[Rules/01 - Princípios de Engenharia]]"]
---

# Corpus Bíblico Canônico

## Objetivo

Transformar os arquivos de entrada em um corpus versionado e validado pelo
Versum, adequado tanto para o seed completo quanto para downloads por livro.
O corpus v1 publicado é a única fonte canônica mantida pelo Versum.

## Estrutura

```text
bible/
  corpus/v1/
    manifest.json              # versão, contagens e hashes
    bible.json                 # todos os 73 livros
    books/
      01-gn.json
      ...
      73-ap.json
```

Cada livro canônico terá `id`, `order`, `name`, `testament` e capítulos com
versículos em ordem. O texto será aparado e o prefixo inicial `[n]` será
removido, reproduzindo a regra usada pelo seed anterior.

## Garantias

- entrada com exatamente 73 livros, na ordem de `listalivros.json`;
- capítulos positivos e sequenciais; versículos positivos e não decrescentes;
- repetições de número preservadas como partes do mesmo versículo (`9.1`,
  `9.2`), sem descartar texto de origem;
- nome e total de capítulos coerentes com o catálogo bruto;
- texto não vazio, sem espaços externos e sem prefixo de numeração;
- o arquivo consolidado e cada livro publicado são verificados por SHA-256.

> [!warning]
> A fonte possui lacunas e repetições de numeração em alguns capítulos. Elas
> são preservadas para não inventar, deslocar ou descartar referências; o
> manifesto registra as quantidades para auditoria.

## Etapas

1. Criar testes de normalização e de rejeição de dados inconsistentes.
2. Implementar `bible/tools/cmd/bible-normalize`, sem dependências externas.
3. Gerar `bible/corpus/v1/` e validar o resultado integralmente.
4. Validar o corpus versionado por meio de `bible/tools/cmd/bible-integrity` no CI.

> [!note]
> O corpus em `bible/corpus/v1/` é a fonte de dados do Versum. O verificador
> não depende de uma cópia bruta dos dados.

---

◀ [[Plans/Archive/01 - Fundação da API Go|Fundação da API Go]] · [[Plans/Archive/_Index|Arquivo]] ▶
