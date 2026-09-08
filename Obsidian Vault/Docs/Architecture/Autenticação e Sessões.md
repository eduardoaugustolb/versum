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

Dados pessoais, incluindo e-mail e metadados de dispositivo, são cifrados em
repouso, inclusive em réplicas e backups. As chaves ficam fora do banco, sob
gestão por ambiente e com rotação. Quando a aplicação precisar localizar um
registro por um dado cifrado, usa um índice cego com HMAC versionado; o valor
original não é mantido em texto puro apenas para permitir busca.

Solicitações de magic link criam o token e um evento de entrega na mesma
transação. O worker envia o e-mail a partir da outbox; o token bruto necessário
para a URL existe apenas em payload cifrado do evento. Consulte a
[[Docs/Decisions/003 - Outbox para Magic Links|ADR da outbox]].

## Regras

- Uma sessão pertence a um dispositivo e pode ser revogada.
- Redirects de magic link usam allowlist.
- Tokens, credenciais e URLs assinadas nunca entram em logs.
- Eventos de outbox com segredo mantêm payload cifrado, versão de chave, retry
  com backoff e observabilidade de falhas.
- Dados pessoais não entram em logs, métricas, traces ou backups sem a mesma
  proteção de cifragem e controle de acesso.
- Download offline, progresso, push e qualquer dado pessoal exigem autenticação.
- A leitura pública do catálogo no web não exige conta.

---

◀ [[Docs/Architecture/Visão Geral|Visão geral]] · próxima: [[Docs/Architecture/Sincronização Offline|Sincronização offline]] ▶
