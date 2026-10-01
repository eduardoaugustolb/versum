---
title: "04 - Privacidade e Consentimento"
section: Plans
subsection: Active
type: implementation-plan
status: planned
date: 2026-09-27
tags: [versum, plans, api, go, privacy, lgpd, consentimento]
up: "[[Plans/Active/_Index|Planos Ativos]]"
prev: "[[Plans/Active/03 - Autenticação e Sessões]]"
next: "[[Plans/Archive/_Index|Arquivo]]"
related: ["[[Docs/Architecture/Privacidade e Consentimento]]", "[[Docs/Architecture/Autenticação e Sessões]]", "[[Docs/Architecture/Sincronização Offline]]", "[[Rules/02 - Segurança]]", "[[PRD]]"]
---

# Privacidade e Consentimento (LGPD)

## Objetivo

Concluir controles de privacidade necessários à autenticação e preparar o
lançamento, com base no comportamento de
[[Docs/Architecture/Privacidade e Consentimento|Privacidade e Consentimento]].
O projeto está em pré-alpha; esta sequência não declara bases aprovadas, política
publicada ou controles operacionais comprovados. Jurídico e controlador validam
as escolhas e textos finais.

## Decisões a validar antes da implementação

- [ ] Registrar finalidade, campos, responsável, base e retenção por atividade
  (arts. 6º, 7º/11 e 37): conta, entrega, segurança e provas de atendimento.
- [x] Definir conta/acesso comum pelo art. 7º, V: acordo de uso gratuito e
  pedido do titular. App aberto ao público, sem venda/captação; leitura sem
  conta, conta opcional. Termos, aviso e consentimentos são distintos.
- [ ] Implementar referência à versão dos termos no pedido e confirmar controle
  do e-mail antes de ativar conta; limpar pedidos/identidades não confirmados.
- [ ] Avaliar se associação ao serviço/progresso permite inferência religiosa
  (arts. 5º, II e 11); hipóteses comuns não substituem bases de sensíveis.
- [ ] Para segurança com art. 7º, IX, realizar avaliação de legítimo interesse
  (art. 10), sem usar essa base para dados sensíveis.
- [ ] Fixar duração/gatilho e justificativa de retenção para conta, tokens,
  sessões, outbox e evidências. 24h e a antiga faixa 30–90 dias são propostas
  do projeto, não prazos da lei; separar segredo de entrega de metadados.

## Estado técnico observado

- [x] E-mail cifrado AES-GCM e índice cego HMAC versionado.
- [x] Solicitação de magic link com token aleatório/hash, TTL e resposta neutra.
- [x] Transação de usuário/token/outbox com payload cifrado.
- [x] Rate limit por IP/e-mail no endpoint; índice Redis de e-mail ainda usa
  SHA-256 determinístico, a substituir por HMAC secreto.
- [ ] Consumo do magic link, sessão autenticada e worker de entrega concluídos.
- [ ] Token bruto cifrado disponível na outbox para entrega: payload atual contém
  apenas identificadores, portanto ainda não sustenta o envio do link.
- [ ] Ledger, DSAR e purge integrados. Estruturas/migrations e expiração de
  token não equivalem a esses fluxos implementados.

## Implementação em `feat/auth`

Complementar [[Plans/Active/03 - Autenticação e Sessões|Plano 03]], mantendo estes
critérios verificáveis de privacidade e segurança (arts. 6º e 46):

- [ ] Concluir geração/persistência do token bruto apenas cifrado na outbox,
  sem retorno indevido, log, trace ou URL de erro; invalidar links conforme
  política e limpar segredos assim que sua conservação deixar de ser necessária.
- [ ] Concluir consumo atômico de uso único e sessões revogáveis; autorização
  no servidor, hash dos segredos e cookies web com proteção apropriada.
- [ ] Concluir worker com lease/entrega após commit, retries limitados e erro
  redigido; definir descarte de eventos e metadados.
- [ ] Limitar tamanho do body e validar entradas. Trocar SHA-256 da chave Redis
  de e-mail por HMAC com segredo/versionamento e preservação do TTL.
- [ ] Documentar IP real atrás de proxy confiável; impedir cliente de escolher
  IP por header arbitrário e definir comportamento seguro de falha do cache.
- [ ] Implementar purge de token/sessão/outbox com prazos aprovados, testes de
  gatilho, falha/reexecução e cobertura de cache/filas; cascata do banco não cobre
  automaticamente outbox sem vínculo referencial.
- [ ] Completar re-cifragem de e-mail/outbox antes de remover chave antiga;
  prever retomada, validação e recuperação da rotação.
- [ ] Revisar erros/logs de dependências e testar ausência de e-mail, token,
  payload sensível e URLs de autenticação em logs/respostas públicas.
- [ ] Disponibilizar suporte ao canal dos titulares: consulta/exportação,
  correção e exclusão autenticadas, trilha mínima e informação de compartilhamento.
  Endpoints propostos `GET /me`, `GET /me/export`, `DELETE /me` não cobrem todos
  os direitos; atender arts. 18–19 inclusive por procedimento humano quando cabível.
- [ ] Se uma finalidade efetiva usar consentimento, implementar versão
  recuperável do texto e ledger mínimo append-only por finalidade; concessão e
  revogação sem checkbox genérico ou pré-marcado. Não acoplar auth a consentimento
  obrigatório de conta se a base aprovada for outra.
- [ ] Conciliar retenção da prova com exclusão; justificar IP/user-agent/snapshot
  antes de coletá-los e não conservar ledger indefinidamente.
- [ ] Fechar testes do fluxo pedido → entrega → consumo → sessão → revogação,
  incluindo concorrência, idempotência, autorização e falhas; rodar verificações
  do projeto após reconciliar testes locais incompletos.

## Gates antes de convites externos ou produção

- [ ] Aprovar inventário/bases, retenção e responsáveis; mapear desenvolvimento,
  Redis, filas, réplicas e backups (arts. 6º, 15–16 e 37).
- [ ] Publicar política com informações do art. 9º e termos, indicando versão;
  confirmar identidade/endereço/contatos do controlador e canais reais.
- [ ] Definir encarregado/enquadramento aplicável e publicar canal; dispensa
  eventual de designação não elimina obrigações de atendimento (art. 41;
  Res. 2/2022 e 18/2024).
- [ ] Aprovar runbook e responsável por incidentes: avaliação de notificabilidade,
  prazos aplicáveis e registro de todos os incidentes por 5 anos (art. 48;
  Res. 15/2024). Não comunicar incidente fictício como exercício.
- [ ] Confirmar fornecedores ativos, dados, papéis, contratos/instruções, regiões
  e mecanismo de transferências quando necessário (arts. 33–39; Res. 19/2024).
  SCCs são um mecanismo disponível, não a única alternativa.
- [ ] Verificar TLS, custódia/acesso às chaves, discos/réplicas/backups, restauração
  que reaplica exclusões e monitoramento sem PII indevida (art. 46).
- [ ] Decidir controles de menores pelo melhor interesse, público e riscos
  (art. 14; Lei 15.211/2025). Art. 9º, §1º sobre conteúdo impróprio não exige por
  si age-gate universal para todo leitor bíblico; avaliar art. 24 e regulamentação.
- [ ] Avaliar impacto e necessidade de RIPD; eventual gate interno é registrado
  como decisão do projeto, não obrigação automática por dado sensível (art. 38).

## Módulos futuros: requisitos para seus planos

| Módulo | Entregas antes de coletar os respectivos dados |
| :-- | :-- |
| Progresso e sync offline | Avaliar inferência religiosa/base art. 11; minimizar cursor/histórico; justificar sequência/dispositivo/timestamps; retenção local/remota, SQLite, exportação, exclusão, revogação e impedir recriação por eventos antigos; avaliar RIPD |
| Push/lembretes | Finalidade opcional e base próprias; token como dado pessoal; desligar envio/apagar vínculo e limpar fila na revogação; conteúdo discreto na tela bloqueada; provedor/região |
| Analytics e métricas PRD | Base por métrica antes da agregação; retenção de eventos, correlação 7/30d e abertura de lembretes; teste de reidentificação e limiares; anonimização arts. 5º, XI e 12, dado anonimizado art. 5º, III; sem ranking/pressão |
| Compartilhamento | Sem sessão em link/imagem; referrer e URLs seguros; imagem sem metadados pessoais; armazenamento e retenção definidos |
| Web | Política/termos/consentimentos distintos; remover token da URL/histórico; cookies protegidos e coleta de analytics condicionada à base aprovada |
| Android | Permissões mínimas, proteção local, logout/exclusão e política de backup do SO; deep link sem vazamento; dados de dispositivo e menores no mapa |

## Critérios de conclusão

- Bases/finalidades e retenção aprovadas pelo responsável identificado, com
  decisões e exceções registradas; nada marcado como aprovado só por estar no plano.
- Fluxo auth e controles acima testados; purge executável, direitos atendidos e
  ledger disponível apenas para finalidades cuja base seja consentimento.
- Gates externos verificáveis com evidências de publicação, canais, fornecedores
  e infraestrutura; artefato redigido não equivale a publicado/contrato assinado.
- Cada plano futuro incorpora sua linha da tabela e só autoriza coleta após
  cumprir as condições aplicáveis.

---

◀ [[Plans/Active/03 - Autenticação e Sessões|Autenticação e sessões]] · próxima: [[Plans/Archive/_Index|Arquivo]] ▶
