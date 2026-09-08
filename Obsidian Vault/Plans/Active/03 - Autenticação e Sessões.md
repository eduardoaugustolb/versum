---
title: "03 - Autenticação e Sessões"
section: Plans
subsection: Active
type: implementation-plan
status: planned
date: 2026-09-07
tags: [versum, plans, api, go, auth, security, sessions]
up: "[[Plans/Active/_Index|Planos Ativos]]"
prev: "[[Plans/Active/_Index|Planos Ativos]]"
next: "[[Plans/Archive/_Index|Arquivo]]"
related: ["[[Docs/Architecture/Autenticação e Sessões]]", "[[Docs/Decisions/003 - Outbox para Magic Links]]", "[[Rules/02 - Segurança]]", "[[Plans/Archive/02 - Catálogo Bíblico]]"]
---

# Autenticação e Sessões

## Objetivo

Implementar a identidade do MVP sem senha: a pessoa solicita um magic link,
consome o token uma única vez e passa a ter sessões revogáveis por dispositivo.
O catálogo continua público; progresso, sincronização, downloads offline,
push e qualquer dado pessoal passam a depender de uma sessão autenticada.

## Decisões do MVP

- O e-mail é a identidade da pessoa. O endereço é normalizado antes de ser
  persistido, cifrado em repouso e não é exposto em respostas públicas.
- Dados pessoais devem ser cifrados em repouso, inclusive em réplicas e
  backups. A chave não fica no banco nem no repositório; seu armazenamento,
  rotação e acesso ficam sob um serviço de gestão de chaves configurado por
  ambiente.
- Para localizar uma identidade sem descriptografar a tabela inteira, manter
  um índice cego separado: HMAC versionado do e-mail normalizado. O HMAC não
  substitui o valor cifrado nem pode ser usado para reconstruir o endereço.
- O token do magic link é aleatório, curto e de uso único. O banco armazena
  apenas seu hash, com expiração e consumo atômico.
- O endpoint de solicitação responde de forma indistinguível para e-mails
  existentes e inexistentes, evitando enumeração de contas.
- A sessão tem identidade própria, pertence a um dispositivo e pode ser
  revogada individualmente ou em conjunto.
- No web, a sessão usa cookie `httpOnly`, `Secure` em produção e política
  `SameSite` definida pelo fluxo. O Android trocará o link por uma sessão
  revogável em etapa própria do cliente.
- Tokens, cookies, e-mails e URLs de autenticação nunca entram em logs.
- O envio de e-mail fica atrás de uma porta; o primeiro adapter pode ser um
  transportador local para desenvolvimento e um provedor configurável para
  produção.
- Token, evento de entrega e transições de sessão que exigem consistência são
  coordenados por uma unidade de trabalho específica de identidade. A entrega
  de magic link usa outbox transacional conforme a
  [[Docs/Decisions/003 - Outbox para Magic Links|ADR 003]].

## Escopo técnico

### 1. Modelo e migrations

Criar as tabelas e restrições para:

- `users`: identidade estável, e-mail cifrado, índice cego para lookup, versão
  da chave e datas de criação/atualização;
- `login_tokens`: hash do token, usuário, expiração, consumo e índices para
  consulta e limpeza;
- `sessions`: hash do segredo da sessão, usuário, dispositivo, expiração,
  revogação e datas de uso.
- `outbox_events`: tipo, payload cifrado, versão da chave, tentativas,
  disponibilidade para retry, estado de processamento, erro redigido e datas
  de criação/processamento.

As invariantes críticas devem ser garantidas pelo PostgreSQL. Tokens consumidos
e sessões revogadas não podem voltar a ser válidos sob concorrência.
`LoginToken`, evento de outbox e qualquer criação associada devem confirmar ou
falhar juntos na mesma transação.

O módulo deve definir também ciclo de vida dos dados pessoais: coleta mínima,
retenção limitada ao necessário, exportação e exclusão da conta, com limpeza
de tokens, sessões, dispositivos e registros derivados. Logs, métricas,
traces, réplicas e backups devem passar pela mesma revisão de minimização e
proteção; nenhum segredo ou dado pessoal deve aparecer em texto puro.

### 2. Domínio e portas

Organizar o módulo `auth` por funcionalidade, com entidades e erros de domínio
independentes de HTTP, PostgreSQL, e-mail e cookies. Definir portas pequenas
para relógio, aleatoriedade, hash, repositório transacional, emissor de
mensagens e armazenamento de sessão. O relógio retorna instantes em UTC; datas
civis com fuso do usuário só recebem value object próprio quando houver regra de
calendário no domínio.

Definir uma `IdentityAccessUnitOfWork` que disponibiliza repositórios de
usuário, token e outbox vinculados à mesma transação. A outbox recebe um payload
cifrado contendo o token bruto de entrega; a tabela de `login_tokens` mantém
apenas o hash. O worker é o único consumidor autorizado a decriptar esse
payload para construir o link e enviar o e-mail.

Casos de uso previstos:

- `RequestMagicLink`: normaliza o e-mail, cria ou encontra a identidade,
  invalida tokens anteriores conforme a política e grava o evento de entrega
  na outbox, sem chamar o provedor de e-mail na transação;
- `ConsumeMagicLink`: valida o token, marca-o como consumido em operação
  atômica e cria uma sessão;
- `AuthenticateSession`: valida a sessão e retorna a identidade;
- `RevokeSession`: encerra a sessão atual;
- `RevokeAllSessions`: encerra todas as sessões da identidade.

### 3. HTTP e composição

Adicionar endpoints para solicitar e consumir magic links, consultar a sessão
atual e sair. O handler apenas traduz entrada, cookie e erros para HTTP; as
regras ficam nos casos de uso.

Adicionar middleware de autenticação que injeta a identidade no contexto para
rotas privadas. `/books` e capítulos continuam acessíveis sem autenticação.

Redirects do magic link devem aceitar somente destinos de uma allowlist
configurada. O contrato precisa impedir open redirect e não deve devolver o
token em logs ou respostas de erro.

### 4. Worker de outbox e operação

Implementar worker separado para reservar eventos pendentes com `FOR UPDATE
SKIP LOCKED`, decriptar o payload somente em memória durante a entrega e chamar
o emissor de e-mail. A confirmação de entrega marca o evento como processado;
falhas usam retry com backoff e limite de tentativas.

O desenho é *at-least-once*: duplicidade de envio é possível após falha entre o
provedor aceitar a mensagem e a confirmação no banco. Templates, provedor e
observabilidade devem tratar esse risco. Métricas e alertas mínimos incluem
eventos pendentes, atrasados, falhos e tentativas esgotadas.

### 5. Testes e verificação

Cobrir, no mínimo:

- normalização de e-mail e respostas indistinguíveis;
- token inválido, expirado, consumido e consumo concorrente;
- expiração e revogação de sessão;
- isolamento entre dispositivos e revogação global;
- cookie seguro e middleware em rotas públicas e privadas;
- redirects fora da allowlist;
- integração com PostgreSQL para constraints e transações;
- atomicidade entre criação do token e gravação da outbox;
- reserva concorrente de eventos, retry com backoff e recuperação após crash;
- payload da outbox cifrado, rotação de chave e ausência do token bruto em
  tabelas, logs, métricas e traces;
- comportamento seguro diante de entrega duplicada;
- ausência de segredos nos logs.
- cifragem e decifragem com chave versionada, incluindo rotação sem perda de
  acesso;
- impossibilidade de recuperar o e-mail a partir do índice cego;
- retenção, exportação e exclusão dos dados pessoais.

Executar `go test ./...`, `go vet ./...` e os testes de integração com banco
antes de marcar o plano como concluído.

## Fora deste plano

- senha, OAuth, login social ou MFA;
- recuperação de conta por suporte manual;
- implementação do cliente Android;
- progresso, sincronização offline, push e lembretes;
- rate limiting distribuído e observabilidade avançada, além das interfaces
  necessárias para não deixar o fluxo aberto.

O desenho técnico não substitui a definição jurídica de bases legais,
encarregado, prazos de retenção e atendimento a titulares; esses itens devem
ser registrados antes da operação em produção.

## Critérios de conclusão

- Uma pessoa consegue solicitar e consumir um magic link válido uma única vez.
- Um token nunca é persistido em texto puro nem aparece em logs.
- A criação do token e o agendamento da entrega são atômicos; o worker consegue
  recuperar falhas sem expor o token bruto fora do payload cifrado da outbox.
- A sessão criada pode ser validada, expirada e revogada por dispositivo.
- Rotas privadas rejeitam requisições sem sessão válida.
- O catálogo público permanece funcionando sem conta.
- Migrations, testes, configuração e documentação de operação estão versionados.

---

◀ [[Plans/Active/_Index|Planos Ativos]] · próxima: [[Plans/Archive/_Index|Arquivo]] ▶
