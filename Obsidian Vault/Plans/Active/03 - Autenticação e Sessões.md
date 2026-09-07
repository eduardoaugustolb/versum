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
related: ["[[Docs/Architecture/Autenticação e Sessões]]", "[[Rules/02 - Segurança]]", "[[Plans/Archive/02 - Catálogo Bíblico]]"]
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

## Escopo técnico

### 1. Modelo e migrations

Criar as tabelas e restrições para:

- `users`: identidade estável, e-mail cifrado, índice cego para lookup, versão
  da chave e datas de criação/atualização;
- `login_tokens`: hash do token, usuário, expiração, consumo e índices para
  consulta e limpeza;
- `sessions`: hash do segredo da sessão, usuário, dispositivo, expiração,
  revogação e datas de uso.

As invariantes críticas devem ser garantidas pelo PostgreSQL. Tokens consumidos
e sessões revogadas não podem voltar a ser válidos sob concorrência.

O módulo deve definir também ciclo de vida dos dados pessoais: coleta mínima,
retenção limitada ao necessário, exportação e exclusão da conta, com limpeza
de tokens, sessões, dispositivos e registros derivados. Logs, métricas,
traces, réplicas e backups devem passar pela mesma revisão de minimização e
proteção; nenhum segredo ou dado pessoal deve aparecer em texto puro.

### 2. Domínio e portas

Organizar o módulo `auth` por funcionalidade, com entidades e erros de domínio
independentes de HTTP, PostgreSQL, e-mail e cookies. Definir portas pequenas
para relógio, aleatoriedade, hash, repositório transacional, emissor de
mensagens e armazenamento de sessão.

Casos de uso previstos:

- `RequestMagicLink`: normaliza o e-mail, cria ou encontra a identidade,
  invalida tokens anteriores conforme a política e solicita o envio;
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

### 4. Testes e verificação

Cobrir, no mínimo:

- normalização de e-mail e respostas indistinguíveis;
- token inválido, expirado, consumido e consumo concorrente;
- expiração e revogação de sessão;
- isolamento entre dispositivos e revogação global;
- cookie seguro e middleware em rotas públicas e privadas;
- redirects fora da allowlist;
- integração com PostgreSQL para constraints e transações;
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
- A sessão criada pode ser validada, expirada e revogada por dispositivo.
- Rotas privadas rejeitam requisições sem sessão válida.
- O catálogo público permanece funcionando sem conta.
- Migrations, testes, configuração e documentação de operação estão versionados.

---

◀ [[Plans/Active/_Index|Planos Ativos]] · próxima: [[Plans/Archive/_Index|Arquivo]] ▶
