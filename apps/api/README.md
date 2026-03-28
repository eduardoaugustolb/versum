# Versum API

API RESTful para o projeto Versum, construída com Bun, Hono e Drizzle ORM.

## Quick Start

```bash
# Instalar dependências
bun install

# Rodar em desenvolvimento
bun run src/index.ts

# Rodar CLI
bun run cli
```

## Estrutura do Projeto

```
src/
├── cli/               # Ferramentas de linha de comando
├── modules/           # Módulos da aplicação (auth, bible, users)
├── view-models/       # Modelos de resposta da API
├── infrastructure/   # Configurações de infraestrutura
└── utils/            # Utilitários
```

## Variáveis de Ambiente

```env
DATABASE_URL=postgresql://user:password@host:5432/database
PORT=3000
```

## CLI

A CLI oferece ferramentas para gerenciamento do banco de dados e geração de documentação:

```bash
bun run cli
```

Opções disponíveis:
- **Bíblia** - Seed e gerenciamento da Bíblia
- **Database** - Operações de banco de dados
- **Gerar Docs OpenAPI** - Gera documentação da API em markdown

## Documentação da API

A documentação completa está disponível em:
- [`src/cli/output/api-documentation.md`](./src/cli/output/api-documentation.md)

### Visão Geral

- **Versão:** 1.0.0
- **OpenAPI:** 3.0.0
- **Formato de resposta:** JSON
- **Convenções:** camelCase para propriedades de API

### Autenticação

A API utiliza autenticação via **Magic Link**:
1. Usuário solicita link mágico via `/api/v1/auth/magic-link`
2. Link contém token enviado por email
3. Token é validado em `/api/v1/auth/magic-link?token=...`
4. Sessão criada via cookie `__Host-session`

### Endpoints Públicos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/public/bible/books` | Listar livros da Bíblia |
| `GET` | `/api/v1/public/bible/books/{dynamicId}` | Obter livro por slug/nome |
| `GET` | `/api/v1/public/bible/books/{dynamicId}/chapters` | Listar capítulos |
| `GET` | `/api/v1/public/bible/books/{dynamicId}/chapters/{number}` | Obter capítulo |
| `GET` | `/api/v1/public/bible/books/{dynamicId}/chapters/{number}/verses` | Listar versículos |
| `GET` | `/api/v1/public/bible/books/{dynamicId}/chapters/{number}/verses/{verse}` | Obter versículo |

### Autenticação

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/auth/magic-link` | Enviar magic link |
| `GET` | `/api/v1/auth/magic-link` | Validar token e autenticar |
| `POST` | `/api/v1/auth/logout` | Encerrar sessão |

### Usuários

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/users/@me` | Obter usuário autenticado |
| `PATCH` | `/api/v1/users/@me` | Atualizar usuário autenticado |
| `GET` | `/api/v1/users/{username}` | Obter usuário público |

## Schemas Principais

### User

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `string` (uuid) | ID único do usuário |
| `username` | `string` | Nome de usuário único |
| `name` | `string` | Nome de exibição |
| `email` | `string` | E-mail do usuário |
| `bio` | `string` | Biografia (nullable) |
| `pictureUrl` | `string` | URL da foto de perfil (nullable) |
| `createdAt` | `string` | Data de criação |

### Book

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `string` (uuid) | ID único do livro |
| `name` | `string` | Nome do livro |
| `slug` | `string` | Slug do livro |
| `niceName` | `string` | Nome amigável |
| `testament` | `string` | "OLD" ou "NEW" |
| `totalChapters` | `integer` | Número total de capítulos |

### PaginationViewModel

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `currentPage` | `integer` | Página atual |
| `totalPages` | `integer` | Número total de páginas |
| `totalItems` | `integer` | Número total de itens |
| `itemsPerPage` | `integer` | Itens por página |
| `hasNextPage` | `boolean` | Existe próxima página |
| `hasPrevPage` | `boolean` | Existe página anterior |

## Scripts

```bash
bun run src/index.ts       # Iniciar servidor
bun run cli                # Abrir CLI
bun run lint              # Verificar código
bun run typecheck         # Verificar tipos TypeScript
```
