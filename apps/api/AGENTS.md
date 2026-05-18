@../AGENTS.md

# API App — Contexto Específico

## Stack
- **Runtime:** Bun (Node.js compatível)
- **Framework:** Hono + OpenAPIHono (with Zod OpenAPI)
- **ORM:** Drizzle ORM (PostgreSQL)
- **Auth:** Sessão em cookie httpOnly + Magic Link por e-mail
- **Email:** Resend

## Módulos

Cada módulo em `src/modules/<nome>/` deve seguir:

```
<nome>/
├── controllers/     # Handlers HTTP (Hono Context)
├── db/              # Schemas Drizzle + relations
├── repositories/    # Acesso a banco
├── routes/          # Definição de rotas OpenAPI
├── schemas/         # Schemas Zod OpenAPI
└── services/        # Lógica de negócio
```

## Convenções Específicas
- **Versão de rota:** Sempre usar `v1` — nome de arquivo: `auth.v1.controller.ts`, `auth.v1.service.ts`
- **ViewModels:** Sempre usar `SuccessViewModel.create()` nas respostas
- **Erros:** Usar classes de erro em `utils/app/errors/` (BadRequestError, NotFoundError, etc.)
- **Imports no módulo:** Usar `@/` alias ao importar de fora do módulo
- **Imports locais:** Usar path relativo SEM extensão `.ts`
