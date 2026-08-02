---
title: "02 - Segurança"
section: Rules
type: rule
status: active
tags: [versum, rules, security]
up: "[[Rules/_Index|Regras]]"
prev: "[[Rules/01 - Princípios de Engenharia]]"
next: "[[Rules/03 - Padrão do Vault]]"
---

# 02 — Segurança

## Regra

- Armazenar tokens de uso único e sessões somente como hashes quando possível.
- Limitar tentativas de login, sync e geração de imagens.
- Validar autorização no servidor em toda operação pessoal.
- Nunca registrar segredos, tokens, e-mail, URLs assinadas ou progresso em logs.
- Usar transações e restrições do banco para invariantes críticos.
- Degradar cache de forma segura; não permitir que Redis se torne a única fonte
  de autorização ou estado.

## Por quê

O aplicativo guarda hábitos de leitura e identificação pessoal. Segurança e
privacidade devem ser propriedades do fluxo, não uma correção posterior.

---

◀ [[Rules/01 - Princípios de Engenharia|Princípios]] · próxima: [[Rules/03 - Padrão do Vault|Padrão do vault]] ▶
