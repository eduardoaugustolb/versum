---
title: PRD
tags:
  - versum
  - product
  - mvp
status: approved
---

# PRD — Versum

## Visão

Versum é um leitor bíblico para quem quer atravessar a Bíblia inteira no próprio
ritmo. A experiência é simples e contemplativa: a pessoa lê, retoma de onde
parou, pode ficar offline e compartilha uma passagem quando isso fizer sentido.

O produto não é uma rede social. Não há feed, likes, ranking, sequência pública
ou qualquer métrica que transforme a vida espiritual em performance.

## Problema

Muitas pessoas desejam ler a Bíblia por inteiro, mas não conseguem manter a
continuidade. Aplicativos tradicionais frequentemente funcionam apenas online,
perdem o ponto da leitura ou trocam profundidade por conteúdo rápido.

## Proposta de valor

> Leia a Bíblia inteira no seu ritmo, com seu progresso salvo e disponível onde
> você estiver.

## Público inicial

Pessoas católicas e cristãs que desejam ler a Bíblia inteira com constância,
especialmente quem nunca conseguiu concluí-la.

## Primeiro lançamento

### Jornadas do usuário

1. A pessoa entra por um link seguro enviado por e-mail.
2. No web, ela explora e lê o catálogo bíblico online.
3. No Android, ela baixa conteúdo para continuar lendo sem conexão.
4. Ao avançar, o aplicativo registra o progresso localmente e o sincroniza ao
   voltar para a internet.
5. Em outro dispositivo autenticado, ela retoma o estado já confirmado.
6. Ela escolhe dias e horário para um lembrete diário de leitura.
7. Ela compartilha um link ou uma imagem de capítulo ou versículo.

### Limites do lançamento

Ficam fora do escopo: feed, seguidores, curtidas, comentários, grupos, perfil
social, avatar, reflexões editoriais e iOS.

## Privacidade

O progresso de leitura é privado por padrão. Uma futura opção de compartilhar
qualquer dado deve ser explícita, granular e reversível. Progresso não é prova
de devoção e não participa de ranking.

## Métricas úteis

- pessoas que retomam a leitura após 7 e 30 dias;
- capítulos lidos por pessoas ativas;
- taxa de sincronização bem-sucedida;
- downloads offline concluídos;
- lembretes entregues e abertos, sem usar esses dados para pressionar o usuário.

## Próximos documentos

- [[Docs/Arquitetura|Arquitetura]]
