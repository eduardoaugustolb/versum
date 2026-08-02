---
title: "Autenticação e Sessões"
section: Docs
subsection: Architecture
type: architecture
status: approved
tags: [versum, docs, architecture, auth, security]
up: "[[Docs/Architecture/_Index|Arquitetura]]"
prev: "[[Docs/Architecture/Visão Geral]]"
next: "[[Docs/Architecture/Sincronização Offline]]"
related: ["[[Rules/02 - Segurança]]"]
---

# Autenticação e Sessões

🏠 [[_Index|Home]] › 📚 [[Docs/_Index|Documentação]] › 📐 [[Docs/Architecture/_Index|Arquitetura]] › **Autenticação**

Magic link é a única forma de entrada do MVP. O token é de uso único, expira em
pouco tempo e é armazenado apenas como hash. No web, a sessão usa cookie
`httpOnly`. No Android, um deep link troca o magic link por uma sessão revogável
no armazenamento seguro do aparelho.

## Regras

- Uma sessão pertence a um dispositivo e pode ser revogada.
- Redirects de magic link usam allowlist.
- Tokens, credenciais e URLs assinadas nunca entram em logs.
- Download offline, progresso, push e qualquer dado pessoal exigem autenticação.
- A leitura pública do catálogo no web não exige conta.

---

◀ [[Docs/Architecture/Visão Geral|Visão geral]] · próxima: [[Docs/Architecture/Sincronização Offline|Sincronização offline]] ▶
