# API Development Guide

## Project Structure

```
apps/api/src/
├── modules/
│   ├── auth/
│   │   ├── db/           # Drizzle schemas
│   │   ├── repositories/
│   │   ├── services/
│   │   └── routes/       # Hono routes
│   ├── bible/
│   │   ├── db/
│   │   ├── repositories/
│   │   ├── services/
│   │   └── routes/
│   └── users/
│       └── ...
├── view-models/          # Response DTOs
└── infrastructure/
    └── db/               # DB config
```

## Creating a New Module

1. Create Drizzle schema in `modules/<name>/db/`
2. Create repository in `modules/<name>/repositories/`
3. Create service in `modules/<name>/services/`
4. Create routes in `modules/<name>/routes/`
5. Register routes in server

## Naming Convention

**CRITICAL:** Follow `.ai/docs/naming-convention.md`

- Database columns: `snake_case`
- API responses: `camelCase`

## View Models Pattern

Always wrap responses with View Models:

```typescript
import { SuccessViewModel } from "@/view-models/default/success.view-model";
import { PaginationViewModel } from "@/view-models/default/pagination.view-model";

return SuccessViewModel.create(data, PaginationViewModel.create({ ... }));
```

## OpenAPI Documentation

Update `apps/client/openapi.yaml` whenever you add/modify endpoints:

- Use Portuguese for descriptions
- Define request/response schemas
- Property names MUST be `camelCase`
