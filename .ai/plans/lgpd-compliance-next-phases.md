# LGPD Compliance — Próximas Fases

> Branch: `feat/lgpd-compliance`
>
> Plano de implementação para completar a adequação à LGPD.

---

## Sumário

- [Fase 3 — Exportação de Dados (Portabilidade)](#fase-3--exportação-de-dados-portabilidade)
- [Fase 4 — Tela de Onboarding com Consentimento](#fase-4--tela-de-onboarding-com-consentimento)
- [Fase 5 — Página de Política de Privacidade](#fase-5--página-de-política-de-privacidade)
- [Fase 6 — Job de Retenção e Purga](#fase-6--job-de-retenção-e-purga)
- [Fase 7 — DPO + Plano de Notificação de Violação](#fase-7--dpo--plano-de-notificação-de-violação)
- [Estrutura de Testes](#estrutura-de-testes)

---

## Fase 3 — Exportação de Dados (Portabilidade)

**Base legal:** Art. 18, II + V (LGPD) — direito de acesso e portabilidade

**Endpoint:** `GET /api/v1/users/@me/export`

### Funcionamento

1. AuthMiddleware valida a sessão (cookie)
2. Controller chama `UserServiceV1.exportUserData({ userId })`
3. Service busca todos os dados do usuário em paralelo:
   - `users` — email, createdAt
   - `profiles` — username, name, bio, pictureUrl
   - `sessions` — createdAt, ip, userAgent, expiresAt (excluir tokenHash)
   - `journey_readings` + `discovery_readings` — histórico de leitura
   - `marks` — anotações
   - `likes` — favoritos
   - `consent_logs` — histórico de consentimento
4. Retorna JSON estruturado com todos os dados (200)
5. Nenhum dado sensível (token hashes) é exposto

### Resposta (JSON)

```json
{
  "exportedAt": "2026-05-15T10:00:00.000Z",
  "user": {
    "email": "user@example.com",
    "createdAt": "2026-01-01T00:00:00.000Z"
  },
  "profile": {
    "username": "john",
    "name": "John Doe",
    "bio": "Bio text",
    "pictureUrl": "https://..."
  },
  "sessions": [
    {
      "createdAt": "...",
      "ip": "127.0.0.1",
      "userAgent": "Mozilla/...",
      "expiresAt": "..."
    }
  ],
  "readingHistory": {
    "journey": [],
    "discovery": []
  },
  "annotations": [],
  "likes": [],
  "consentLogs": []
}
```

### Implementação

#### 1. Schemas (`users/schemas/v1/users.v1.common.schema.ts`)

- `exportUserDataResponseSchema` — zod.object com todos os campos acima
- Usar `.openapi()` para documentação

#### 2. Repositórios — novos métodos de query

| Repositório | Método | Descrição |
|---|---|---|
| `user.repository.ts` | `findByIdWithAllRelations` | User + profile + sessions + readings + marks + likes + consentLogs |
| OU usar queries separadas nos repositórios existentes |

**Decisão de design:** usar queries separadas nos repositórios existentes, evitando um método monolítico. O service orquestra as chamadas.

Métodos já existentes que podem ser reutilizados:
- `UserRepository.findById` ✅
- `ProfileRepository.findByUserId` ✅
- `AuthRepository.getSessionsByUserId` ✅
- `ConsentLogsRepository.getConsentLogsByUserId` ✅

Métodos novos necessários:
- `JourneyReadingRepository.findByUserId` — se não existir
- `DiscoveryReadingRepository.findByUserId` — se não existir
- `MarksRepository.findByUserId` — se não existir
- `LikesRepository.findByUserId` — se não existir

#### 3. Service (`user.v1.service.ts`)

```ts
async exportUserData({ userId }: { userId: string }): Promise<ExportUserDataResponse>
```

- Busca usuário (lança erro se não existir)
- Busca profile, sessions, readings, marks, likes, consentLogs em paralelo com `Promise.all`
- Monta objeto de resposta
- Retorna dados

#### 4. Controller (`users.v1.controller.ts`)

```ts
exportUserData = async (c: Context) => {
  const session = c.get("session") as Session;
  const data = await this.service.exportUserData({ userId: session.userId });
  return c.json(data, 200);
};
```

#### 5. Rota (`users.v1.route.ts`)

```ts
const exportMeRoute = createRoute({
  method: "get",
  path: "/@me/export",
  tags: ["Users"],
  security: [{ cookieAuth: [] }],
  responses: {
    200: { content: { "application/json": { schema: exportUserDataResponseSchema } } },
    ...createErrorResponses([401, 404, 429, 500]),
  },
});
```

**Ordem de registro:** após `router.use("/@me", authMiddleware.validateSession)`, registrar a rota.

### Testes necessários

- `user.v1.service.test.ts` — 3 novos casos:
  - `exportUserData` retorna todos os dados do usuário
  - `exportUserData` lança erro quando usuário não existe
  - `exportUserData` retorna dados parciais quando não há profile/sessions/etc

---

## Fase 4 — Tela de Onboarding com Consentimento

**Base legal:** Art. 7º, I + Art. 8º (LGPD) — consentimento do titular

### Fluxo completo

```
Login via magic link
  → Redireciona para /onboarding
  → Tela exibe checkboxes de consentimento
  → Usuário marca as opções e confirma
  → POST /api/v1/consent (registra no consent_logs)
  → Redireciona para criação de perfil
  → ProfileService verifica hasConsent("profile_content") antes de criar
```

### O que deve ser implementado

#### 1. Tela de onboarding (`apps/client/src/app/(private)/onboarding/page.tsx`)

**Checkboxes:**
- [ ] Permitir armazenar conteúdo do perfil (bio, foto) — `profile_content`
- [ ] Permitir salvar minhas anotações — `annotations`
- [ ] Permitir salvar meus favoritos — `likes`
- [ ] Aceito os Termos de Uso e Política de Privacidade — `terms`

**Botão:** "Confirmar e continuar"

**Comportamento:**
- Ao confirmar, chama `POST /api/v1/consent`
- Em caso de sucesso, redireciona para criação de perfil (`/onboarding/profile`)
- Se o usuário já tiver consentido (checar via `GET /api/v1/consent`), pular onboarding

#### 2. Guardas de rota (client)

Atualizar o sistema de guards para verificar se o consentimento foi dado antes de permitir acesso a:
- Criação de perfil
- Funcionalidades de anotação
- Funcionalidades de like

Isso pode ser feito via:
- Rota protegida que consulta `GET /api/v1/consent`
- Ou via um campo adicional no `GET /users/@me`

#### 3. Verificações de consentimento nos serviços (backend)

Já implementado:
- `ConsentLogServiceV1.hasConsent` ✅

Onde integrar:
- `ProfileServiceV1.createProfile` — verificar `hasConsent("profile_content")` antes de criar
- Quando houver serviço de marks — verificar `hasConsent("annotations")`
- Quando houver serviço de likes — verificar `hasConsent("likes")`

### Testes necessários

- `consent-log.v1.service.test.ts` — já existe com 6 testes ✅
- `profile.v1.service.test.ts` — adicionar teste: `createProfile` bloqueia sem consentimento `profile_content`

---

## Fase 5 — Página de Política de Privacidade

**Base legal:** Art. 9 (LGPD) — direito de informação

### Rota

`GET /privacy` — página pública (rota no client: `apps/client/src/app/(public)/privacy/page.tsx`)

### Conteúdo da página

1. **Dados coletados** (tabela do plano original):
   - Email, nome, username, bio, picture URL
   - IP, User-Agent, comportamento de leitura
   - Anotações, likes, magic link requests

2. **Finalidade do processamento** para cada tipo de dado

3. **Base legal** (LGPD Art. 7º) para cada tipo

4. **Terceiros com acesso a dados:**
   - Resend (email)
   - PostgreSQL (hospedagem)
   - Redis (cache de rate limiting)

5. **Direitos do titular** (LGPD Art. 18):
   - Acessar dados: `GET /api/v1/users/@me/export`
   - Corrigir dados: `PATCH /api/v1/profiles/@me`
   - Excluir conta: `DELETE /api/v1/users/@me`
   - Revogar consentimento: `POST /api/v1/consent` com `granted: false`

6. **DPO** — contato

7. **Períodos de retenção** (tabela com prazos)

8. **Política de cookies** — informar que usa cookie de sessão (HttpOnly)

### Implementação

**Client:** Página Next.js estática em `apps/client/src/app/(public)/privacy/page.tsx`

- Componente de página simples com markdown/JSX
- Links diretos para os endpoints da API

### Testes necessários

- Teste de renderização do componente (se houver padrão de testes no client)
- Atualmente o client não tem testes — não criar se não houver padrão

---

## Fase 6 — Job de Retenção e Purga

**Base legal:** Art. 15-16 (LGPD) — eliminação dos dados após o término da finalidade

### Tabelas e condições

| Tabela | Condição de Purge | Frequência |
|--------|-------------------|------------|
| `magic_links` | `expiresAt < now - 30 days` AND `usedAt IS NOT NULL` | Diária |
| `magic_links` | `expiresAt < now - 7 days` AND `usedAt IS NULL` (abandonados) | Diária |
| `sessions` | `revokedAt IS NOT NULL` AND `revokedAt < now - 90 days` | Semanal |

### Implementação

#### 1. Repositórios — novos métodos de deleção em lote

**AuthRepository:**
```ts
deleteExpiredMagicLinks(): Promise<number>  // retorna qtd deletada
deleteExpiredSessions(): Promise<number>    // retorna qtd deletada
```

#### 2. Serviço de purge

Criar `apps/api/src/modules/auth/services/purge.service.ts`:

```ts
export class PurgeService {
  async purgeExpiredMagicLinks(): Promise<{ deleted: number }>
  async purgeExpiredSessions(): Promise<{ deleted: number }>
  async runDailyPurge(): Promise<{ magicLinks: number; sessions: number }>
}
```

#### 3. Script/CLI

Criar comando CLI (`apps/api/src/cli/modules/purge/`) ou integração com cron externo:

```
bun run src/cli/index.ts purge --daily    # executa purge diário
bun run src/cli/index.ts purge --weekly   # executa purge semanal
```

Ou expor como endpoint administrativo protegido.

### Testes necessários

- Criar `purge.service.test.ts` com testes:
  - `purgeExpiredMagicLinks` deleta links expirados e usados
  - `purgeExpiredMagicLinks` deleta links abandonados (não usados, expirados > 7d)
  - `purgeExpiredSessions` deleta sessões revogadas há mais de 90 dias
  - `runDailyPurge` executa ambas as purgas e retorna totais

---

## Fase 7 — DPO + Plano de Notificação de Violação

**Base legal:** Art. 41 (DPO) + Art. 33-36 (notificação de violação)

### O que implementar

#### 1. Contato do DPO

Adicionar na página de privacidade (Fase 5) um email de contato do DPO.

#### 2. Plano de resposta a incidentes (documentação)

Criar `.ai/docs/incident-response-plan.md` com:

- Definição de violação de dados pessoais
- Fluxo de detecção e contenção
- Responsáveis
- Prazo de notificação à ANPD (Art. 48 — 2 dias úteis)
- Template de notificação aos titulares

#### 3. Alerta de segurança

Integrar com Discord webhook para alertar em casos como:
- Múltiplas tentativas de login falhas (rate limit excedido)
- Tentativas de acesso a tokens inválidos
- Qualquer erro 500 inesperado no auth

### Testes necessários

- Documentação — sem testes
- Se houver implementação de webhook de alerta, testar a função de notificação

---

## Estrutura de Testes

### Status atual dos testes

| Arquivo | Testes | Status |
|---------|--------|--------|
| `consent-log.v1.service.test.ts` | 6 (recordConsents, getUserConsents, hasConsent) | ✅ Completo |
| `user.v1.service.test.ts` | 11 (create, getById, getByEmail, delete) | ✅ Completo (delete adicionado) |

### Testes a criar nesta fase

| Fase | Arquivo | Novos testes |
|------|---------|--------------|
| 3 | `user.v1.service.test.ts` | 3: exportUserData (sucesso, sem dados, usuário inexistente) |
| 4 | `profile.v1.service.test.ts` | 1: createProfile bloqueado sem consentimento |
| 6 | `purge.service.test.ts` (novo) | 4: purge mágic links (usados + abandonados), purge sessions, runDailyPurge |

### Padrão de testes (seguir existente)

```ts
// Mock do repositório
const createMockRepository = () => ({
  metodo: vi.fn<() => Promise<Tipo>>(),
});

// Service com dependências mockadas via construtor
service = new MeuServiceV1({
  repository: mockRepository as unknown as RepositorioReal,
});

// Asserções
expect(result).toEqual(expected);
expect(mockRepository.metodo).toHaveBeenCalledWith({ ... });
```

---

## Ordem de Implementação Sugerida

```
Fase 3 — Export (GET /@me/export)       ← obrigação legal, alta prioridade
Fase 4 — Onboarding UI                   ← obrigação legal, alta prioridade
Fase 5 — Privacy Page                    ← transparência, média prioridade
Fase 6 — Purge Job                       ← minimização, média prioridade
Fase 7 — DPO + Breach Notification       ← governança, média prioridade
```

Cada fase deve ser implementada com seus respectivos testes antes de passar para a próxima.
