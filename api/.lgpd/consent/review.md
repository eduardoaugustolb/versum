# Revisão de consentimento — 2026-10-01

Diagnóstico F3, cenário A. Não é implementação nem aprovação de base legal.

## Evidência atual

`internal/database/postgres/migrations/000003_create-auth.up.sql:1` contém contas, tokens e sessões; não contém versões de política ou consentimentos. `internal/transport/httpapi/identityaccess/routes.go:58` aceita somente e-mail. `internal/identityaccess/application/commands/request_magic_link.go:65` cria usuário sem consultar consentimento. O Plano 04 do Vault, linhas 55–65, mantém esse trabalho pendente.

## Ajustes necessários no plano

- Resolver a base de autenticação antes de tornar `conta_acesso` um checkbox obrigatório; consultar [bases propostas](../legal-basis.md). Ciência da política, aceite contratual e consentimento são registros distintos (LGPD arts. 7º–9º).
- Para finalidades consentidas, registrar titular, finalidade, versão do texto específico, instante do evento, ação e origem verificável; rejeitar finalidade desconhecida e associação a outro titular. O consentimento destacado do progresso deve permanecer separado (arts. 8º e 11, I).
- Ordenar concessão/revogação por instante do evento e desempate determinístico; `granted_at` sozinho não ordena revogações. Bloquear novos tratamentos e propagar a decisão a filas e operadores (arts. 8º, §5º, e 18, IX).
- Append-only protege integridade durante a retenção, mas não significa guardar eternamente. Definir finalidade, base e descarte da prova, inclusive após exclusão; reconciliar com a cascata prevista no Vault (arts. 15–16).
- IP, user-agent e hashes vinculáveis continuam exigindo minimização; não coletar tudo apenas porque compõe uma evidência. Hash isolado não comprova qual texto foi apresentado: conservar a versão do texto e a referência verificável da interface (arts. 6º, III, e 8º, §2º).

## Validação prevista para implementação

Concessão por finalidade; revogação prevalecendo no mesmo instante; idempotência; isolamento entre titulares; versão inexistente; recusa opcional sem bloquear conta; reconsentimento quando necessário; exportação e eliminação conforme retenção aprovada.

Stack observada: Go/PostgreSQL. Não há motivo para introduzir Prisma ou Better Auth somente para cumprir este plano.
