# LGPD Compliance Plan

## 1. LGPD Data Map

### Personal Data Collected

| Data | Table(s) | Purpose | Legal Basis (LGPD) |
|------|----------|---------|--------------------|
| Email | `users`, `magic_links` | Authentication, communication | Art. 7º, V (contract execution) |
| Name | `profiles` | User identification/profile | Art. 7º, V (contract execution) |
| Username | `profiles` | Public profile identifier | Art. 7º, V (contract execution) |
| Bio | `profiles` | User-provided profile content | Art. 7º, I (consent) |
| Picture URL | `profiles` | Avatar display | Art. 7º, I (consent) |
| IP address | `sessions`, Redis | Audit, rate limiting | Art. 7º, IX (legitimate interest) |
| User-Agent | `sessions` | Device identification | Art. 7º, IX (legitimate interest) |
| Reading behavior | `journey_readings`, `discovery_readings` | App functionality | Art. 7º, V (contract execution) |
| Annotations | `marks` | User-generated content | Art. 7º, I (consent) |
| Likes | `likes` | User preference | Art. 7º, I (consent) |
| Magic link requests | `magic_links` | Authentication flow | Art. 7º, V (contract execution) |

### Third-Party Processors

| Processor | Data Received | Purpose |
|-----------|---------------|---------|
| Resend | Email address | Email delivery (magic links) |
| PostgreSQL (hosted) | All personal data | Database storage |
| Redis (hosted) | IP address (transient) | Rate limiting cache |

---

## 2. Gap Analysis

### Already Compliant ✅

| Requirement | Implementation |
|-------------|----------------|
| Security measures (Art. 46) | Argon2 hashing, CSPRNG, HttpOnly cookies, HSTS, rate limiting |
| Data minimization (Art. 6, III) | Only essential data collected |
| Access control | Session-based auth with refresh/rotation |

### Missing 🔴

| # | Requirement | Priority | What's Needed |
|---|-------------|----------|---------------|
| 1 | Consent mechanism (Art. 7º, I + Art. 8º) | **High** | Onboarding consent flow + `consent_logs` table |
| 2 | Data deletion right (Art. 18, VI) | **High** | `DELETE /api/v1/users/@me` endpoint |
| 3 | Data access/portability (Art. 18, II + V) | **High** | `GET /api/v1/users/@me/export` endpoint |
| 4 | Data retention & purge (Art. 15-16) | **Medium** | Cron job + TTL for expired/revoked data |
| 5 | Privacy policy (Art. 9) | **Medium** | Public page with DPO contact, purposes, rights |
| 6 | Consent records (Art. 8º, §2º) | **High** | `consent_logs` table with audit trail |
| 7 | DPO contact (Art. 41) | **Medium** | Email/contact in privacy policy |
| 8 | Data breach notification (Art. 33-36) | **Medium** | Incident response plan + Discord alert |

---

## 3. Implementation Plan

### Phase 1 — Consent & Audit Trail

**What:** `consent_logs` table + onboarding consent flow

#### Table: `consent_logs`

```typescript
// apps/api/src/modules/users/db/consent-logs.table.ts
export const consentLogs = pgTable("consent_logs", {
  id: uuid("id").primaryKey().defaultRandom(),
  userId: uuid("user_id").notNull().references(() => users.id, { onDelete: "cascade" }),
  purpose: varchar("purpose", { length: 100 }).notNull(),
  granted: boolean("granted").notNull(),
  ip: text("ip").notNull(),
  userAgent: text("user_agent").notNull(),
  createdAt: timestamp("created_at", { precision: 3, withTimezone: true }).notNull().defaultNow(),
});
```

**Flow:**
1. After magic link authentication, before profile creation → show consent screen
2. User checks consent checkboxes (profile data, annotations, terms)
3. `POST /api/v1/consent` records each purpose as a row in `consent_logs`
4. API blocks profile creation / annotations if consent not given

**Repository methods:**
- `createConsentLog(params)` — insert row
- `getConsentLogsByUserId({ userId })` — retrieve all consents
- `hasConsent({ userId, purpose })` — check if consent exists and is granted

---

### Phase 2 — Data Export (Portability)

**Endpoint:** `GET /api/v1/users/@me/export`

**Response format:**
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
    "journey": [ ... ],
    "discovery": [ ... ]
  },
  "annotations": [ ... ],
  "likes": [ ... ],
  "consentLogs": [ ... ]
}
```

**Steps:**
1. AuthMiddleware validates session
2. Controller calls `UserServiceV1.exportUserData(userId)`
3. Service queries: user, profile, sessions, readings, marks, likes, consent_logs
4. Returns all data as JSON (200)
5. No sensitive data: token hashes are excluded

---

### Phase 3 — Data Deletion (Right to Erasure)

**Endpoint:** `DELETE /api/v1/users/@me`

**Steps:**
1. AuthMiddleware validates session
2. Optionally require re-authentication (send magic link for confirmation)
3. Controller calls `UserServiceV1.deleteUser(userId)`
4. Service deletes in order:
   - Sessions (manual - no cascade)
   - Profile (manual - no cascade)
   - User (cascade deletes: marks, likes, journey_readings, discovery_readings)
5. Clears session cookie
6. Returns 204 No Content

**Tables affected by cascade (automatic):**
- `marks` (onDelete: cascade)
- `likes` (onDelete: cascade)
- `journey_readings` (onDelete: cascade)
- `discovery_readings` (onDelete: cascade)

**Tables needing manual deletion:**
- `sessions` (no cascade from users)
- `profiles` (no cascade from users)
- `magic_links` (tied to email, not userId — can optionally purge by email)
- `consent_logs` (onDelete: cascade — will be added in Phase 1)

**Confirmation step:**
To prevent accidental deletion, either:
- Require magic link re-auth before deletion (send email, user must click)
- Or use a two-step process: `POST /api/v1/users/@me/delete-request` → email → `GET /confirm-deletion?token=...`

---

### Phase 4 — Data Retention & Purge

**Cron job (or scheduled task):**

| Table | Purge Condition | Frequency |
|-------|----------------|-----------|
| `magic_links` | `expiresAt < now - 30 days` AND `usedAt IS NOT NULL` | Daily |
| `sessions` | `revokedAt IS NOT NULL` AND `revokedAt < now - 90 days` | Weekly |
| `magic_links` | `expiresAt < now - 7 days` AND `usedAt IS NULL` (abandoned) | Daily |

**Redis:** Rate limiter keys have natural TTL via window expiration — no action needed.

---

### Phase 5 — Privacy Policy Page

**Route:** `GET /privacy` (public page)

**Content:**
1. What data is collected (table from section 1)
2. Purpose of processing
3. Legal basis for each type of data
4. Third-party processors (Resend, DB host, Redis)
5. User rights (LGPD Art. 18): access, correction, deletion, portability, revocation
6. DPO contact information
7. Cookie policy
8. Data retention periods
9. How to exercise rights (links to export/delete endpoints)

---

## 4. File Structure

```
apps/api/src/modules/users/
├── db/
│   ├── users.table.ts              (existing)
│   ├── profiles.table.ts           (existing)
│   ├── users.relations.ts          (existing — add consentLogs relation)
│   └── consent-logs.table.ts       [NEW]
├── repositories/
│   ├── user.repository.ts          (existing — add delete + export methods)
│   ├── user.types.repository.ts    (existing — add types)
│   ├── profile.repository.ts       (existing)
│   └── consent.repository.ts       [NEW]
├── services/
│   ├── user.v1.service.ts          (existing — add delete + export methods)
│   └── profile.v1.service.ts       (existing)
├── controllers/
│   ├── users.v1.controller.ts      (existing — add delete + export handlers)
│   └── profile.v1.controller.ts    (existing)
├── routes/
│   ├── index.ts                    (existing)
│   ├── users.v1.route.ts           (existing — add delete + export routes)
│   └── profile.v1.routes.ts        (existing)
├── schemas/
│   └── v1/
│       ├── users.v1.common.schema.ts    (existing — add export/delete schemas)
│       └── profiles.v1.common.schema.ts (existing)

apps/api/src/modules/auth/repositories/
├── auth.repository.ts              (existing — add bulk revoke method)
├── auth.types.repository.ts        (existing)

apps/api/src/infrastructure/db/
├── schema.ts                      (existing — add consent-logs export)

apps/client/src/app/
├── (public)/privacy/
│   └── page.tsx                   [NEW] — privacy policy page
├── (private)/onboarding/
│   └── page.tsx                   [NEW] — consent collection

.ai/plans/
└── lgpd-compliance.md             [THIS FILE]
```

---

## 5. Routes Summary (New)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/users/@me/export` | Cookie | Export all user data (JSON) |
| `DELETE` | `/api/v1/users/@me` | Cookie | Delete user account and all data |
| `POST` | `/api/v1/consent` | Cookie | Record user consent |
| `GET` | `/api/v1/consent` | Cookie | Get user's consent history |

---

## 6. Priority Order for Implementation

1. **Consent logs table + consent repository** (foundation)
2. **DELETE /api/v1/users/@me** (right to erasure — legal requirement)
3. **GET /api/v1/users/@me/export** (portability — legal requirement)
4. **POST /api/v1/consent + onboarding UI** (consent mechanism)
5. **Privacy policy page** (transparency)
6. **Retention/purge job** (data minimization)
7. **DPO contact + breach notification plan** (governance)

---

## 7. Notes

- **Sessions table:** `sessions` does NOT have `onDelete: cascade` from `users`. Must delete manually.
- **Profiles table:** `profiles` does NOT have `onDelete: cascade` from `users`. Must delete manually.
- **Magic links:** Tied to email (not userId). Can optionally purge on deletion by looking up user's email.
- **Cascade tables (automatic on user delete):** `marks`, `likes`, `journey_readings`, `discovery_readings`.
