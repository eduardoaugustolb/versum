# PRD — Versum

## 1. Visão geral

**Versum** é um aplicativo de leitura bíblica com dinâmica de feed moderno, inspirado em experiências como TikTok, porém com foco total em leitura, constância e reflexão — sem mecânicas de dopamina barata.

O app permite que usuários leiam versículos ou capítulos de forma **aleatória (Descoberta)** ou **sequencial (Jornada)**, interajam apenas por curtidas e acompanhem seu progresso ao longo do tempo.

---

## 2. Problema que o Versum resolve

* Dificuldade em manter constância na leitura bíblica
* Apps tradicionais são engessados e pouco intuitivos
* Experiências modernas focam em entretenimento, não em profundidade

O Versum resolve isso oferecendo:

* Leitura fluida em formato de feed
* Progresso claro
* Experiência minimalista
* Total foco no texto

---

## 3. Objetivo do produto

Criar um ambiente digital que:

* Incentive a leitura diária da Bíblia
* Gere constância sem pressão
* Elimine distrações e complexidade
* Respeite o ritmo individual do usuário

---

## 4. Público-alvo

* Jovens e adultos cristãos
* Pessoas acostumadas a consumir conteúdo em feed
* Usuários que querem ler a Bíblia, mas não mantêm constância

---

## 5. Proposta de valor

> **"Leia a Bíblia no ritmo de hoje, sem perder a essência."**

---

## 6. Funcionalidades principais

### 6.1 Autenticação

* Login via **e-mail (Magic Link)**
* Sessão infinita (Infinity Session)
* Perfil simples

Campos de perfil:

* Nome
* Username
* Foto

---

### 6.2 Modos de leitura

#### 📌 Modo Descoberta (For You)

* Versículos ou capítulos exibidos de forma aleatória
* Conteúdo não se repete
* Ideal para leitura diária

---

#### 🛤️ Modo Jornada (Sequencial)

* Leitura contínua do Gênesis ao Apocalipse
* Sem pular capítulos
* Progresso salvo automaticamente

Indicadores:

* Livro atual
* Capítulo atual
* % de progresso geral

---

### 6.3 Interações

O usuário **não pode marcar textos nem salvar trechos**.

Interações permitidas:

* Curtir versículos individuais
* Curtir capítulos completos

Objetivo:

* Reduzir complexidade
* Evitar fragmentação do texto
* Manter foco total na leitura

---

### 6.4 Curtidas

* Curtidas são **privadas por padrão**
* O usuário pode **alterar a visibilidade** das curtidas
* Não existem **números públicos por padrão**
* O usuário pode **optar por exibir métricas**
* Curtidas influenciam o feed pessoal

Biblioteca do usuário:

* Versículos curtidos
* Capítulos curtidos

---

### 6.5 Rede social

* Seguir usuários
* Feed "Amigos" mostrando apenas:

  * Versículos curtidos
  * Capítulos curtidos

Regras:

* Sem comentários
* Sem reposts
* Sem métricas públicas

---

## 7. Estrutura de dados (alto nível)

### Usuário

* id
* username
* nome
* foto_url
* created_at

### Versículo

* id
* livro
* capitulo
* versiculo
* texto

### Capítulo

* id
* livro
* capitulo

### Curtida

* user_id
* alvo_id (versículo ou capítulo)
* tipo (VERSICULO | CAPITULO)
* created_at

### Leitura

* user_id
* alvo_id
* modo (descoberta | jornada)
* lido_em

---

## 8. Requisitos não funcionais

* Alta performance
* Scroll suave
* Interface limpa
* Dark mode por padrão
* Baixo consumo de dados

---

## 9. Stack

| Camada | Tecnologia |
|--------|-----------|
| Frontend | **Next.js 16** (App Router) + **React 19** |
| Backend / API REST | **Hono** (OpenAPIHono) + **Drizzle ORM** |
| Database | **PostgreSQL** |
| Auth | **E-mail Magic Link** + Sessão infinita (cookie httpOnly) |
| Package Manager | **Bun** |
| Linter/Formatter | **Biome** |
| Test | **Vitest** |
| Codegen | **Orval** (OpenAPI → TanStack Query + Zod) |
| Email | **Resend** |
| Cache | Redis (opcional) |
| UI | Tailwind CSS v4 + shadcn/ui + phosphor-icons |
| Animações | GSAP + SplitText + Lenis |
| Seed bíblico | JSON normalizado em `src/assets/json/bible.json` |

---

## 10. Métricas de sucesso (MVP)

* Usuários ativos diários
* Versículos ou capítulos lidos por sessão
* Taxa de retorno (D1 / D7)
* Curtidas por usuário

---

## 11. Fora do escopo

* Comentários
* Marcações de texto
* Salvamento de trechos
* Monetização
* Anúncios

---

## 12. Visão futura

* Leitura offline
* Estatísticas pessoais de leitura
* Versões diferentes da Bíblia

---

**Versum não é sobre acumular interações.**

É sobre constância.
