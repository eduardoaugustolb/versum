# AGENTS.md — Versum

**Leia isso antes de fazer qualquer alteração no projeto.**

---

## Stack

| Camada | Tecnologia |
|--------|-----------|
| Frontend (cliente) | Next.js 16, React 19, Tailwind CSS v4, shadcn/ui |
| Backend (API) | Hono (OpenAPIHono), Drizzle ORM, Zod |
| Database | PostgreSQL |
| Auth | E-mail Magic Link + Infinity Session (cookie httpOnly) |
| Package Manager | Bun |
| Linter/Formatter | Biome |
| Test | Vitest |
| Codegen | Orval (OpenAPI → TanStack Query + Zod) |

---

## Documentação Obrigatória

Sempre leia estes arquivos antes de codificar:

- `.ai/prd.md` — Product Requirements Document
- `.ai/docs/naming-convention.md` — Database, API e código
- `.ai/docs/api-development.md` — Como criar endpoints
- `.ai/docs/git-flow.md` — Branches, commits e sandbox
- `.ai/rules/` — Regras de desenvolvimento

---

## Estrutura do Projeto

```
versum/
├── .ai/                 # Documentação e regras globais
├── apps/
│   ├── api/             # Backend REST (Hono)
│   ├── client/          # Frontend (Next.js App Router)
│   └── landing-page/    # Landing page marketing
├── packages/
│   └── logger/          # @versum/logger
└── AGENTS.md            # Este arquivo
```

---

## Regras Obrigatórias para Agentes de IA

### 1. Versionamento (Git Flow) — REGRA ABSOLUTA
- **Branch base:** `development` — **NUNCA trabalhe diretamente nela**
- **Nomenclatura:** `feat/<descricao>`, `fix/<descricao>`, `refactor/<descricao>`, `docs/<descricao>`, `chore/<descricao>`
- **Commits:** Conventional Commits — `tipo(escopo): descrição`
- **Sandbox:** Branches criadas por IA usam prefixo `feat/ai-` ou `fix/ai-`
- **PR:** Sempre abrir PR para `development` ao finalizar
- ⚠️ **Violar esta regra = quebra o fluxo de revisão e integração contínua**

### 2. Antes de Codificar — CHECKLIST OBRIGATÓRIO
1. Leia `AGENTS.md` (este arquivo) — **todo, do início ao fim**
2. Leia `.ai/docs/git-flow.md` — **entenda o fluxo de branches**
3. Leia `.ai/prd.md`
4. Leia `.ai/docs/` pertinente ao que vai fazer
5. Leia `.ai/rules/` pertinente
6. **Verifique a branch atual** com `git branch --show-current`
   - Se estiver em `development`, **PARE** e crie uma sub-branch
   - Se estiver em outra branch sem prefixo `ai-`, considere criar uma nova
7. Crie branch a partir de `development`: `git checkout -b <tipo>/ai-<descricao> development`
8. **Atualize o repositório local** sem comprometer alterações:
   ```bash
   git stash
   git pull origin development --rebase
   git stash pop
   ```
9. Verifique exemplos existentes no código antes de criar algo novo

### 3. Regras de Código
- **Database:** `snake_case` nas colunas, `camelCase` nas propriedades Drizzle
- **API Responses:** `camelCase`, sempre usar ViewModels (`SuccessViewModel.create()`)
- **API módulos:** Seguir padrão `controllers/db/repositories/routes/schemas/services`
- **Versão de API:** Usar sufixo `.v1.` em todos os arquivos de módulo (ex: `auth.v1.controller.ts`)
- **Imports:** Usar `@/` alias (configurado em tsconfig)
- **Components React:** PascalCase (`UserProfile.tsx`)
- **Testes:** Colocar ao lado do arquivo que testam (`service.test.ts`)

### 4. Importante
- Este é um projeto **Bun** — nunca use npm/yarn/pnpm
- Next.js 16 tem breaking changes — consulte `node_modules/next/dist/docs/` se necessário
- Nunca commite `.env`, `.certs`, `node_modules`
- Execute `biome check` e `tsc --noEmit` antes de cada commit
