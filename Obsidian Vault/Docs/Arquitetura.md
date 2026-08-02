---
title: Arquitetura
tags:
  - versum
  - architecture
  - backend
  - mobile
status: approved
---

# Arquitetura — Versum

## Estrutura

```text
versum/
  api/       API Go, worker e migrations
  web/       leitor online em Next.js
  mobile/    aplicativo Android em React Native/Expo
  infra/     ambiente local e deploy
  Obsidian Vault/
```

Os clientes são projetos independentes. Não existe workspace JavaScript nem um
pacote de regras compartilhado entre TypeScript e Go. A API publica o contrato
que os clientes consomem.

## Backend

A API usa Go, `net/http`, `chi`, PostgreSQL via `pgx`, Redis e S3 compatível.
Ela segue arquitetura hexagonal por funcionalidade:

- casos de uso definem portas pequenas;
- adapters HTTP traduzem requests e respostas;
- adapters externos implementam acesso a Postgres, Redis, S3, e-mail, FCM e
  Discord;
- `cmd/api` e `cmd/worker` montam as dependências.

Casos de uso não conhecem bibliotecas HTTP, banco, cache, storage ou push.
PostgreSQL é a fonte de verdade. Redis é usado apenas para cache, rate limit e
locks de curta duração. S3 guarda imagens de compartilhamento.

## Autenticação

Magic link é a única forma de entrada do MVP. O token é de uso único, expira em
pouco tempo e é armazenado somente como hash. No web, a sessão usa cookie
`httpOnly`. No Android, um deep link troca o magic link por uma sessão revogável
guardada no armazenamento seguro do aparelho.

Sessões pertencem a dispositivos e podem ser revogadas. Redirects são validados
por allowlist; tokens, credenciais e URLs assinadas nunca entram em logs.

## Sincronização offline

O Android usa SQLite para conteúdo baixado, estado de leitura e uma outbox de
eventos pendentes. Cada evento possui um ID único, dispositivo, sequência local,
referência de leitura e instante de criação.

Ao receber uma sincronização, a API processa eventos em transação. Uma restrição
de unicidade torna reenvios idempotentes. Eventos são aditivos: um dispositivo
atrasado pode acrescentar histórico, mas nunca reduzir o ponto mais avançado já
confirmado. A API devolve o estado canônico para reconciliação local.

O app tenta sincronizar após uma ação, quando recupera conectividade e quando é
aberto novamente. O sistema operacional pode suspender o processo, portanto não
há promessa de sincronização contínua; dados não sincronizados podem ser
perdidos em uma desinstalação.

## Notificações

O usuário configura dias e horários, com opção simples de pausar ou desligar.
Notificações locais ajudam quando possível; o worker envia push Android por FCM
como complemento. Jobs são idempotentes, respeitam fuso horário e removem tokens
inválidos.

## Compartilhamento e operações

Links canônicos do web oferecem metadados para redes sociais. Cards de capítulos
e versículos são versionados no S3. O Discord recebe apenas alertas técnicos do
seed e do worker, sem dados pessoais ou conteúdo de leitura.

## Critérios de qualidade

1. Repetir o mesmo evento de leitura não muda o progresso duas vezes.
2. Eventos de dispositivos diferentes e fora de ordem não reduzem o estado
   canônico.
3. Conteúdo já baixado continua legível sem internet.
4. Magic links não podem ser reutilizados e não aparecem em logs.
5. Um link compartilhado não revela identidade ou progresso de quem o enviou.
6. Testes cobrem casos de uso, adapters Postgres, contrato HTTP e cenários de
   duplicação, retry e concorrência na sincronização.

## Decisões pendentes

- versão e licença do texto bíblico;
- provedor de e-mail, Redis, S3 e hospedagem;
- biblioteca de SQLite e estado no aplicativo;
- suporte a iOS após a primeira versão Android.
