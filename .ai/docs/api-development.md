# API Development Guide

## Project Structure

```
apps/api/src/
├── modules/
│   ├── auth/
│   │   ├── controllers/     # Hono request handlers
│   │   ├── db/              # Drizzle schemas + relations
│   │   ├── helpers/         # (opcional) Funções auxiliares do módulo
│   │   ├── repositories/    # Acesso a banco de dados
│   │   ├── routes/          # Definição de rotas OpenAPI
│   │   ├── schemas/         # Schemas Zod OpenAPI (request/response)
│   │   └── services/        # Lógica de negócio
│   ├── bible/
│   ├── consent-logs/
│   ├── interactions/
│   ├── reading/
│   └── users/
├── view-models/              # Response DTOs
├── middlewares/               # Hono middlewares
└── infrastructure/
    ├── db/                   # DB config + schema central
    ├── redis/                # Redis client
    └── resend/               # Email provider
```

## Creating a New Module

1. Create `modules/<nome>/` com a estrutura: `controllers/`, `db/`, `repositories/`, `routes/`, `schemas/`, `services/`
2. Todos os arquivos do módulo DEVEM usar sufixo de versão: `auth.v1.controller.ts`
3. Crie Drizzle schema em `db/`
4. Crie repository em `repositories/`
5. Crie service em `services/`
6. Crie schemas Zod em `schemas/`
7. Crie routes em `routes/`
8. Registre as rotas em `src/modules/routes.ts`

## Naming Convention

**CRITICAL:** Follow `.ai/docs/naming-convention.md`

- Database columns: `snake_case`
- Drizzle properties: `camelCase`
- API responses: `camelCase`
- Arquivos de módulo: `<nome>.v1.<tipo>.ts` (ex: `auth.v1.controller.ts`, `bible.v1.service.ts`)

## View Models Pattern

**SEMPRE** use View Models nas respostas:

```typescript
import { SuccessViewModel } from "@/view-models/default/success.view-model";
import { PaginationViewModel } from "@/view-models/default/pagination.view-model";

return c.json(SuccessViewModel.create(data, PaginationViewModel.create({ ... })), 200);
```

NUNCA retorne objetos soltos:
```typescript
// ❌ ERRADO
return c.json({ profile }, 200);
return c.json({ message: "ok" }, 200);

// ✅ CORRETO
return c.json(SuccessViewModel.create(profile), 200);
return c.json(SuccessViewModel.create(undefined, undefined, "ok"), 200);
```

## Error Handling

Use as classes de erro em `utils/app/errors/`:

```typescript
import { BadRequestError, NotFoundError, UnauthorizedError } from "@/utils/app/errors";

// Dentro do controller/service:
throw new NotFoundError("Profile not found");
throw new BadRequestError("Invalid email");
throw new UnauthorizedError("Session expired");
```

Os erros são automaticamente capturados pelo error handler.

## OpenAPI Documentation

Update `apps/client/openapi.yaml` whenever you add/modify endpoints:

- Use Portuguese for descriptions
- Define request/response schemas
- Property names MUST be `camelCase`

## Imports

- **Dentro do módulo:** path relativo SEM extensão `.ts`
- **Fora do módulo:** `@/` alias (ex: `import { App } from "@/utils/app"`)
- **View Models:** `@/view-models/default/success.view-model`

## Testes

- Colocar arquivo de teste ao lado do arquivo testado: `auth.v1.service.test.ts`
- Usar Vitest (`describe`/`it`/`expect`)
- Mockar dependências externas (banco, redis, email)
