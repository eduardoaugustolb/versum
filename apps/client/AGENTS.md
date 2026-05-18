@../AGENTS.md

# Client App — Contexto Específico

## Stack
- **Framework:** Next.js 16 (App Router)
- **UI:** React 19, Tailwind CSS v4, shadcn/ui, phosphor-icons
- **Data Fetching:** TanStack Query v5, Orval (codegen de OpenAPI)
- **Forms:** TanStack Form
- **Animações:** GSAP, SplitText

## ⚠️ Next.js 16 — Breaking Changes
APIs, convenções e estrutura de arquivos DIFEREM de versões anteriores do Next.js.
Sempre consulte `node_modules/next/dist/docs/` antes de escrever código.

## Convenções Específicas
- **Componentes:** PascalCase, em pastas próprias (`UserProfile/index.tsx`)
- **Features:** Toda funcionalidade em `src/features/<nome>/` com estrutura self-contained
- **DAL (Data Access Layer):** Código gerado pelo Orval fica em `src/dal/orval/` (fetch, tanstackQuery, zod)
- **Route Guards:** Em `src/app/(private)/routes/` — middleware de autenticação e onboarding
- **UI Components:** Em `src/components/ui/` — componentes shadcn customizados
