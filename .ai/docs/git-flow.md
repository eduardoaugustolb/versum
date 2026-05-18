# Git Flow — Versum

## Estrutura de Branches

```
main (produção)
  └── development (integração)
       ├── feat/<nome>          # Novas features
       ├── fix/<nome>           # Correções de bug
       ├── refactor/<nome>      # Refatorações
       ├── docs/<nome>          # Documentação
       ├── chore/<nome>         # Manutenção
       ├── release/<versao>     # Preparação de release
       └── hotfix/<nome>        # Correção urgente em produção
```

## Fluxo Padrão

1. Crie branch a partir de `development`
2. Desenvolva na branch
3. Abra PR (Pull Request) para `development`
4. Após revisão, faça merge em `development`
5. Releases: `development` → `release/x.y.z` → `main`

## Pós-Merge (PR aprovado)

Sempre que um PR for aprovado e mergeado no remoto:

```bash
# 1. Atualizar branches base
git checkout development
git pull origin development

git checkout main
git pull origin main

# 2. Limpar branch de feature (local e remoto)
git branch -d <nome-da-branch>
git push origin --delete <nome-da-branch>

# 3. Voltar para development para iniciar próxima tarefa
git checkout development
```

**Segurança:** `git fetch` / `git pull` são operações seguras. Elas NUNCA sobrescrevem arquivos não trackeados (`.env`, `.certs`, `node_modules` estão no `.gitignore`).

## Sandbox para Agentes de IA

Branches criadas por IA **devem** usar prefixo `ai-` para fácil identificação:

| Tipo | Exemplo |
|------|---------|
| Feature | `feat/ai-add-dark-mode` |
| Fix | `fix/ai-login-redirect` |
| Refactor | `refactor/ai-extract-service` |
| Docs | `docs/ai-update-readme` |
| Chore | `chore/ai-remove-dead-code` |

### Regras para Agentes de IA

1. **Sempre** crie branch a partir de `development` — **nunca trabalhe direto em `development`**
2. Antes de começar, **verifique a branch atual** com `git branch --show-current`
3. Se estiver em `development`, **PARE** e crie uma sub-branch imediatamente
4. Nunca force push (`git push --force`)
5. Commits devem seguir **Conventional Commits** (veja abaixo)
6. Execute `biome check` e `tsc --noEmit` antes de cada commit
7. Ao finalizar, abra PR para `development` com descrição clara
8. Nunca faça merge direto em `main` ou `development`

### Fluxo Correto ao Iniciar uma Tarefa

```bash
# 1. Verificar branch atual
git branch --show-current
# → development  (PRECISO CRIAR UMA SUB-BRANCH)

# 2. Atualizar repositório local sem perder alterações
git stash
git pull origin development --rebase
git stash pop

# 3. Criar sub-branch a partir de development
git checkout -b feat/ai-<descricao> development

# 4. Desenvolver, commitar, abrir PR
```

## Conventional Commits

```
<tipo>(<escopo opcional>): <descrição>

[corpo opcional]

[rodapé opcional]
```

### Tipos

| Tipo | Uso |
|------|-----|
| `feat` | Nova funcionalidade |
| `fix` | Correção de bug |
| `refactor` | Mudança que não altera comportamento |
| `docs` | Documentação |
| `chore` | Manutenção, build, dependências |
| `test` | Testes |
| `style` | Formatação, estilo de código |
| `perf` | Performance |

### Exemplos

```
feat(api): add version suffix to discovery module
fix(client): remove unused onNext param from ConsentStepView
refactor(api): extract session validation to middleware
docs(ai): unify AGENTS.md across all apps
chore: remove empty directories
test(api): add coverage for purge service
```

### Escopos

| Escopo | Aplica-se a |
|--------|-------------|
| `api` | `apps/api/` |
| `client` | `apps/client/` |
| `landing` | `apps/landing-page/` |
| `logger` | `packages/logger/` |
| `ai` | `.ai/`, `AGENTS.md` |
| `ci` | `.github/workflows/` |
