---
title: "003 - Outbox para Magic Links"
section: Docs
subsection: Decisions
type: adr
status: accepted
date: 2026-09-08
tags: [versum, docs, adr, auth, security, outbox]
up: "[[Docs/Decisions/_Index|Decisões]]"
prev: "[[Docs/Decisions/002 - Progresso por Eventos]]"
next: "[[Rules/_Index]]"
related: ["[[Docs/Architecture/Autenticação e Sessões]]"]
---

# 003 — Outbox para Magic Links

## Contexto

Solicitar um magic link envolve persistir um `LoginToken` e entregar seu token
bruto ao usuário por e-mail. Fazer essas operações separadamente cria duas
falhas graves: o token pode ser persistido sem que o e-mail seja enviado, ou o
e-mail pode ser enviado para um token cuja transação não confirmou.

O token bruto não pode ser salvo na tabela de tokens: apenas seu hash é usado
para validar o link. Porém, o worker de entrega precisa recuperar o segredo uma
única vez para montar a URL enviada ao usuário.

## Decisão

Usar uma unidade de trabalho específica de identidade para, na mesma transação
PostgreSQL:

1. criar ou localizar o usuário conforme a política do caso de uso;
2. criar o `LoginToken` contendo somente o hash;
3. gravar um evento de outbox para entrega do magic link;
4. confirmar a transação.

O evento guarda `user_id`, `login_token_id` e um payload cifrado que contém o
token bruto necessário para entrega. Nem o token bruto, nem a URL completa, nem
o e-mail em texto puro podem aparecer em logs, métricas, traces ou colunas sem
cifragem. O payload informa a versão da chave de cifragem para permitir rotação.

O worker processa a outbox de forma assíncrona: reserva eventos pendentes com
`FOR UPDATE SKIP LOCKED`, decripta o payload somente durante a entrega, obtém o
e-mail protegido do usuário e envia a mensagem. Em sucesso, registra a entrega;
em falha, aumenta as tentativas e agenda retry com backoff. Depois do limite de
tentativas, o evento exige observação e tratamento operacional.

## Garantias e limites

- A criação do token e o agendamento do e-mail são atômicos.
- Entrega é *at-least-once*: um provedor pode aceitar a mensagem e a confirmação
  do worker falhar. O template e o provedor devem tolerar duplicidade.
- O token continua sendo aleatório, de uso único, com expiração curta e salvo
  somente como hash na tabela principal.
- Consumo deve ser atômico e impedir reutilização concorrente.
- A outbox não substitui rate limit, allowlist de redirects, logs redigidos,
  gestão externa de chaves, controle de acesso ou monitoramento de falhas.

## Consequências

- O caso de uso não chama o provedor de e-mail dentro da transação.
- Repositórios de usuário, token e outbox devem operar sobre a mesma `pgx.Tx`
  durante a unidade de trabalho.
- É necessário um worker, schema de outbox, política de retry e métricas de
  eventos pendentes/falhos.
- O payload cifrado adiciona rotação de chave e limpeza segura após entrega ou
  expiração.

## Alternativas descartadas

- Enviar e-mail antes do commit: pode entregar link inválido.
- Enviar e-mail após o commit no processo HTTP: perde a entrega em crash entre
  commit e chamada ao provedor.
- Salvar token bruto junto do hash: amplia a exposição de segredos em dumps,
  backups e acessos indevidos ao banco.
- Usar somente fila externa sem registro transacional: mantém a janela de falha
  entre banco e publicação da mensagem.

---

◀ [[Docs/Decisions/002 - Progresso por Eventos|ADR 002]] · próxima: [[Rules/_Index|Regras]] ▶
