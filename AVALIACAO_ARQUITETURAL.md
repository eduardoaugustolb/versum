# Avaliação Arquitetural do Versum

**Data da avaliação:** 5 de agosto de 2026  
**Escopo:** documentação do `Obsidian Vault/`, implementação da API Go, corpus bíblico, infraestrutura local, testes e automações existentes.  
**Objetivo:** fornecer evidências e opções para que a liderança técnica decida como simplificar a arquitetura, corrigir inconsistências e preparar o produto para evolução e escala.

## Resumo executivo

A stack e a arquitetura macro são adequadas ao produto. Go, PostgreSQL, `pgx`, `chi`, SQLite no mobile e uma API stateless formam uma base simples e capaz de atender uma escala significativa sem microserviços.

O principal problema não é a escolha das tecnologias. É a distribuição da complexidade:

- o projeto introduz cedo abstrações e cerimônias de Clean/Hexagonal Architecture;
- a abstração genérica de SQL não entrega a portabilidade prometida;
- regras arquiteturais, planos e implementação divergem entre si;
- sincronização, autenticação, observabilidade e operação, que determinam a robustez real do produto, permanecem subespecificadas;
- a implementação atual da API ainda não está integrada de ponta a ponta e não compila;
- não existe CI para impedir regressões da API.

A recomendação principal é manter um **monólito modular Go**, assumir PostgreSQL explicitamente, usar interfaces semânticas pequenas definidas pelos consumidores e remover abstrações de infraestrutura sem um segundo caso real. Antes de expandir o backend, devem ser formalizados os contratos de sincronização e autenticação.

## Conclusão geral

| Dimensão | Avaliação | Síntese |
| :-- | :-- | :-- |
| Stack | Adequada | Go, PostgreSQL, `pgx`, `chi`, Expo/React Native e SQLite são escolhas coerentes. |
| Arquitetura macro | Adequada | Monorepo e monólito modular são suficientes; não há justificativa para microserviços. |
| Idiomatismo Go | Precisa de revisão | Há abstrações inspiradas em arquiteturas orientadas a classes e duplicação parcial das APIs de `pgx`/`database/sql`. |
| Limites de domínio | Inconsistentes | As regras dizem que adapters ficam fora do domínio, mas o SQL concreto está no pacote `catalog`. |
| Escalabilidade | Potencialmente boa, ainda não demonstrada | A API pode ser stateless, porém faltam readiness, shutdown, observabilidade, estratégia de eventos e processo de deploy. |
| Sincronização | Risco alto | O diferencial central do produto ainda não possui semântica completa nem modelo operacional. |
| Segurança | Parcial | Existem bons princípios, mas autenticação e sessões não estão especificadas em nível implementável. |
| Testes e entrega | Insuficientes | Há testes unitários, mas faltam integração PostgreSQL, CI da API e validação do composition root. |
| Corpus | Boa integridade mecânica | Manifest e hashes são positivos, mas seed atômico, versão no banco e licenciamento precisam ser resolvidos. |

## Arquitetura documentada

A arquitetura pretendida é:

```text
Next.js web -----------------+
                             |
React Native/Expo Android ---+--> API Go/chi --> PostgreSQL
  +-- SQLite + outbox        |          +-----> Redis
                             |          +-----> S3
                             |          +-----> e-mail/FCM
                             |
                             +--> worker Go para jobs assíncronos
```

Responsabilidades documentadas:

- API Go: contrato público, autenticação, regras de negócio e sincronização;
- worker Go: lembretes, geração assíncrona e manutenção;
- PostgreSQL: estado sincronizado e projeção operacional do catálogo;
- corpus versionado: autoridade do conteúdo bíblico;
- Redis: cache, rate limit e locks curtos;
- S3: imagens de compartilhamento;
- Android: conteúdo local, estado de leitura e outbox em SQLite.

Referências:

- `Obsidian Vault/Docs/Architecture/Visão Geral.md:18-50`
- `Obsidian Vault/Docs/Architecture/Visão Geral.md:88-95`
- `Obsidian Vault/Docs/Architecture/Sincronização Offline.md:18-29`

## Pontos positivos

- A visão de produto é clara, privada e sem escopo social desnecessário.
- Go e PostgreSQL são adequados para uma API de leitura, autenticação e sincronização.
- `pgx` e `chi` oferecem baixo overhead e pouca complexidade acidental.
- Os handlers atuais são stateless, favorecendo escala horizontal.
- O código propaga `context.Context` até a persistência.
- As consultas usam parâmetros, sem SQL montado por concatenação.
- Interfaces de repositório pequenas e próximas dos consumidores são uma boa direção.
- SQLite com outbox é apropriado para um aplicativo Android offline-first.
- PostgreSQL como autoridade e cache como dado descartável é um princípio correto.
- O corpus possui geração determinística, manifesto, hashes e verificações de integridade.
- Não há concorrência ou estado global desnecessário no código atual.
- As dependências externas da API são reduzidas.

## Achados prioritários

### 1. A abstração SQL não entrega portabilidade real

**Severidade:** alta  
**Categoria:** idiomatismo, complexidade e limites arquiteturais

O plano reconhece que placeholders e construções como `$1` e `ON CONFLICT` são específicos do PostgreSQL, mas afirma que trocar de banco exigiria apenas outro `<Driver>Executor`, sem alterar `catalog.Repository`.

Essa afirmação não se sustenta: outro banco exigiria mudanças nas queries, nos tipos, na semântica de transação e possivelmente no tratamento de erros.

A implementação replica parcialmente APIs de `database/sql` e `pgx`:

```go
type Executor interface {
    QueryRow(ctx context.Context, query string, args ...any) Row
    Query(ctx context.Context, query string, args ...any) (Rows, error)
    Exec(ctx context.Context, query string, args ...any) error
}
```

Consequências:

- mais interfaces e adapters para manter;
- perda de funcionalidades nativas do `pgx`, como `CommandTag`;
- ausência de portabilidade efetiva das queries;
- suporte a transações precisa de novas abstrações;
- pacote próprio chamado `sql`, com potencial de colisão com `database/sql`;
- dificuldade adicional para diagnosticar erros específicos do PostgreSQL.

Evidências:

- `Obsidian Vault/Plans/Active/02 - Catálogo Bíblico.md:33-53`
- `Obsidian Vault/Docs/Architecture/Visão Geral.md:74-81`
- `api/internal/adapters/sql/executor.go:5-20`
- `api/internal/catalog/query.go:13-14`

**Recomendação:** assumir PostgreSQL explicitamente e usar `pgx` diretamente no adapter de persistência. Manter interfaces semânticas como `BookRepository` e `ChapterRepository` no pacote consumidor. Só introduzir suporte a outro banco quando existir um requisito concreto.

### 2. Regras arquiteturais e implementação se contradizem

**Severidade:** alta  
**Categoria:** limites e manutenibilidade

As regras determinam que:

- portas são definidas no caso de uso que depende delas;
- adapters externos ficam fora do domínio;
- dependências apontam para dentro;
- SQL é conhecido apenas pelo adapter PostgreSQL.

Entretanto, `catalog.Repository` e as queries concretas estão em `internal/catalog`, e esse pacote importa `internal/adapters/sql`.

Evidências:

- `Obsidian Vault/Rules/01 - Princípios de Engenharia.md:16-25`
- `Obsidian Vault/Rules/01 - Princípios de Engenharia.md:39-82`
- `api/internal/catalog/repository.go:1-16`
- `api/internal/catalog/query.go:1-14`

Existem duas abordagens coerentes:

| Abordagem | Vantagem | Custo |
| :-- | :-- | :-- |
| Vertical slice pragmático | Menos pacotes e indireção; SQL próximo da funcionalidade | `catalog` deixa de ser um domínio puro e passa a incluir infraestrutura |
| Ports & Adapters | Direção de dependência explícita; domínio sem SQL | Mais um pacote e wiring, justificável onde há regra de negócio real |

**Recomendação:** escolher uma abordagem e documentá-la sem exceções contraditórias. Para este projeto, um adapter PostgreSQL explícito e interfaces pequenas no consumidor oferece o melhor equilíbrio.

### 3. Há cerimônia excessiva para operações simples

**Severidade:** média  
**Categoria:** idiomatismo Go

O padrão determina um struct e um método `Execute` para cada operação. `ListBooks` e `GetChapter`, no estado atual, apenas delegam ao repositório.

Esse padrão pode ser útil quando há autorização, transação, múltiplas dependências, invariantes ou orquestração. Em consultas pass-through, ele aumenta:

- quantidade de arquivos;
- construtores;
- ponteiros e wiring;
- fakes e testes de delegação;
- superfície de mudança para uma operação simples.

**Recomendação:** usar funções ou serviços pequenos para operações simples. Introduzir um caso de uso explícito quando houver comportamento de negócio observável, não como regra obrigatória para toda consulta.

### 4. A API atual não compila e não está conectada ao PostgreSQL

**Severidade:** crítica no estado avaliado  
**Categoria:** implementação e entrega

`NewRouter` recebe `httpapi.Dependencies`, mas `cmd/api` passa diretamente um `health.CheckHealth`:

- `api/cmd/api/main.go:22`
- `api/internal/transport/httpapi/router.go:9`
- `api/internal/transport/httpapi/dependencies.go:8-16`

Resultado confirmado por `go test ./...`:

```text
cannot use health.NewCheckHealth() (value of struct type health.CheckHealth)
as httpapi.Dependencies value in argument to httpapi.NewRouter
```

Também não existe no composition root:

- criação do `pgxpool.Pool`;
- `Ping` no boot;
- criação de `PgxExecutor` e do repositório;
- injeção de `ListBooks` e `GetChapter`;
- fechamento do pool;
- shutdown gracioso.

Como o working tree estava em implementação ativa durante a avaliação, essa falha pode ser transitória. Ainda assim, a ausência de CI permitiu que uma incompatibilidade simples permanecesse sem barreira automatizada remota.

**Recomendação:** finalizar o composition root antes de adicionar novas funcionalidades e criar CI obrigatório para `api/**`.

### 5. O contrato de capítulo não é cumprido pela persistência

**Severidade:** alta  
**Categoria:** corretude

`Chapter` expõe `BookName`, mas `FindChapter` consulta apenas `verses` e nunca preenche esse campo.

Evidências:

- `api/internal/catalog/chapter.go:10-15`
- `api/internal/catalog/query.go:14`
- `api/internal/catalog/repository.go:40-67`

Os testes dos casos de uso e handlers usam stubs que já devolvem `BookName`, por isso não detectam a divergência entre contrato e SQL.

**Recomendação:** consultar o nome do livro no repositório ou removê-lo do contrato. Adicionar teste de integração real com PostgreSQL.

### 6. A semântica de sincronização está incompleta

**Severidade:** alta  
**Categoria:** domínio, concorrência e escalabilidade

A sincronização é o principal problema distribuído do produto. A documentação define eventos idempotentes e uma projeção monotônica, mas não especifica:

- o significado exato de “progresso”;
- a diferença entre posição atual e maior avanço histórico;
- a ordem canônica entre livros, capítulos e partes;
- leitura fora de ordem;
- releitura;
- reset explícito;
- marcação e desmarcação de capítulos;
- chave única `(user_id, device_id, sequence)`;
- detecção de lacunas na sequência local;
- confiança ou não no relógio do cliente;
- tamanho máximo do lote de sincronização;
- concorrência na atualização da projeção;
- retries parciais;
- versionamento do protocolo;
- retenção, compactação ou arquivamento dos eventos.

Evidências:

- `Obsidian Vault/Docs/Architecture/Sincronização Offline.md:18-29`
- `Obsidian Vault/Docs/Decisions/002 - Progresso por Eventos.md:19-34`

Um usuário pode voltar a Gênesis depois de chegar a Salmos. Nesse cenário, “retomar leitura” e “maior ponto alcançado” representam estados diferentes e não podem compartilhar a mesma projeção.

**Recomendação:** não implementar event sourcing genérico. Especificar operações tipadas, registrá-las em uma inbox idempotente e atualizar projeções transacionalmente. Definir cursores de servidor, limites de lote e retenção antes da implementação.

### 7. O seed planejado não publica uma projeção fiel e atômica

**Severidade:** alta  
**Categoria:** dados e operação

O corpus é a autoridade do conteúdo, enquanto PostgreSQL é uma projeção. O seed planejado usa upserts e transações por livro.

Problemas:

- upsert não remove versículos excluídos ou renumerados;
- falha no meio pode deixar uma mistura de versões;
- não há versão ou hash do corpus registrado no banco;
- validar apenas 73 livros não prova que todos os versículos correspondem ao manifesto;
- não existe atualmente um comando funcional e versionado para carregar o catálogo.

Evidências:

- `Obsidian Vault/Plans/Active/02 - Catálogo Bíblico.md:59-65`
- `Obsidian Vault/Plans/Active/02 - Catálogo Bíblico.md:92-93`
- `Obsidian Vault/Plans/Active/02 - Catálogo Bíblico.md:462-474`
- `api/internal/adapters/postgres/migrations/000001_create-catalog.up.sql:1-18`

**Recomendação:** carregar dados em staging, validar contagens e hashes e publicar a versão atomicamente. Registrar versão e hash do manifesto no banco. Como alternativa, servir o catálogo imutável como artefato estático com cache/CDN e manter PostgreSQL apenas para estado pessoal e busca, se necessária.

### 8. O schema não protege invariantes suficientes

**Severidade:** média  
**Categoria:** integridade de dados

Faltam restrições para números positivos, texto não vazio e coerência básica:

```sql
CHECK (chapter_count > 0)
CHECK (chapter > 0)
CHECK (number > 0)
CHECK (part > 0)
CHECK (length(trim(text)) > 0)
```

Outros pontos:

- `CREATE TABLE IF NOT EXISTS` pode esconder um schema incompatível em migrations versionadas;
- o índice `(book_id, chapter)` é provavelmente redundante com a chave primária;
- a coluna `"order"` exige quoting permanente e poderia ser `canonical_order`;
- não há dimensão de tradução ou versão para múltiplos corpus futuros;
- não existe restrição que relacione o maior capítulo com `chapter_count`.

Evidência:

- `api/internal/adapters/postgres/migrations/000001_create-catalog.up.sql:1-18`

**Recomendação:** endurecer o schema e realizar as validações que não cabem em constraints durante a publicação atômica do corpus.

### 9. O serviço ainda não está preparado para operação horizontal

**Severidade:** alta antes de produção  
**Categoria:** operação e escalabilidade

Embora os handlers sejam stateless, faltam:

- liveness e readiness separadas;
- `Ping` do PostgreSQL com timeout no boot;
- shutdown gracioso;
- `ReadHeaderTimeout` e política explícita de timeouts;
- request ID;
- recuperação e observação de panic;
- logs estruturados por requisição;
- métricas de latência, status e saturação do pool;
- tracing ou correlação entre API e worker;
- métricas de backlog, retries e idade do evento mais antigo;
- SLOs, alertas e dashboards;
- processo reproduzível de migration e seed.

O health check atual sempre devolve `ok`, mesmo sem banco disponível:

- `api/internal/health/check.go:13-15`

**Recomendação:** implementar `/livez` e `/readyz`, lifecycle completo do servidor e instrumentação mínima antes de múltiplas réplicas ou deploy de produção.

### 10. Serviços externos estão sendo antecipados

**Severidade:** média  
**Categoria:** complexidade prematura

A visão geral já reserva adapters para Redis, S3, FCM e Discord.

Avaliação por tecnologia:

| Tecnologia | Avaliação |
| :-- | :-- |
| PostgreSQL | Necessário e adequado |
| `pgx` | Adequado e idiomático |
| `chi` | Adequado, simples e de baixo acoplamento |
| SQLite no mobile | Necessário para offline-first |
| S3 | Justificável para imagens compartilhadas |
| FCM | Justificável caso existam notificações remotas |
| Redis | Adiar até existir necessidade medida de cache, rate limit distribuído ou locks |
| Discord | Sem requisito de produto documentado |
| Worker separado | Criar quando existirem jobs assíncronos concretos |

Evidência:

- `Obsidian Vault/Docs/Architecture/Visão Geral.md:39-50`

**Recomendação:** não provisionar nem abstrair dependências sem caso de uso implementado. PostgreSQL pode atender a primeira versão de jobs, leases e controles simples.

### 11. Autenticação ainda não constitui um protocolo implementável

**Severidade:** alta antes da implementação  
**Categoria:** segurança

As decisões iniciais são adequadas: magic link de uso único, hash no banco, cookie `httpOnly`, sessão por dispositivo e redirects allowlisted.

Ainda faltam decisões sobre:

- entropia e algoritmo de hash;
- consumo atômico do token;
- expiração exata e tolerância de relógio;
- proteção contra enumeração de e-mail;
- rate limit por IP, e-mail e dispositivo;
- comportamento diante de scanners automáticos de links;
- cookies `Secure`, `SameSite`, `Path` e proteção CSRF;
- rotação, expiração e revogação global de sessões;
- Android App Links verificados;
- armazenamento e rotação de chaves;
- auditoria sem exposição de dados pessoais.

Evidência:

- `Obsidian Vault/Docs/Architecture/Autenticação e Sessões.md:18-29`

**Recomendação:** escrever um threat model e um contrato detalhado de autenticação antes de iniciar handlers e migrations.

### 12. O corpus possui risco de licenciamento

**Severidade:** alta para publicação  
**Categoria:** jurídico e distribuição

O repositório distribui o texto bíblico, mas não documenta claramente:

- tradução utilizada;
- fonte original;
- titular dos direitos;
- licença do conteúdo;
- direito de redistribuição e transformação;
- atribuição exigida;
- processo para correções editoriais.

A licença MIT do código não concede automaticamente direitos sobre o corpus.

**Recomendação:** separar licença de código e licença de dados, registrar proveniência e confirmar o direito de distribuição antes da publicação web ou em lojas.

## Outros problemas concretos da implementação

### Result sets não são fechados explicitamente

`ListBooks` e `FindChapter` não executam `defer rows.Close()` após uma consulta bem-sucedida:

- `api/internal/catalog/repository.go:21-37`
- `api/internal/catalog/repository.go:46-67`

O `pgx` fecha as linhas quando a iteração termina normalmente, mas um retorno antecipado por erro de `Scan` pode manter a conexão ocupada por mais tempo.

### Validação HTTP é incompleta

Capítulos `0` e negativos chegam à persistência. Valores incompatíveis com `SMALLINT` podem resultar em erro interno em vez de erro de validação.

Também faltam:

- limite e formato de `bookId`;
- contrato padronizado de erros;
- versionamento da API;
- OpenAPI ou contrato equivalente;
- ETag e `Cache-Control` para conteúdo imutável;
- limites de body e headers.

### Testes unitários escondem erros de integração

Os testes com stubs verificam a delegação dos casos de uso e o mapeamento HTTP, mas não cobrem:

- queries reais;
- tipos e scans do PostgreSQL;
- migrations `up` e `down`;
- preenchimento de `BookName`;
- ordenação real dos versículos;
- fechamento de rows;
- startup e shutdown;
- readiness;
- composition root.

### PostgreSQL local é publicado em todas as interfaces

`infra/docker-compose.yml` publica `5432:5432` e usa credenciais triviais. Para desenvolvimento local, é mais seguro usar `127.0.0.1:5432:5432`. Ambientes compartilhados devem usar secrets externos e não publicar o banco diretamente.

### Não existe CI para a API

Os workflows atuais cobrem corpus e vault, mas não `api/**`. Hooks locais não substituem CI remoto.

Um pipeline mínimo deveria executar:

```text
go build ./...
go vet ./...
go test ./...
go test -race ./...
```

Uma etapa separada deve subir PostgreSQL e executar migrations, seed e testes de integração.

## Arquitetura recomendada

### Direção

Manter um monólito modular Go, organizado por capacidades de negócio, com adapters concretos e poucos níveis de indireção.

```text
api/
  cmd/api/
    main.go

  internal/catalog/
    service.go
    types.go
    repository.go

  internal/progress/
    service.go
    events.go
    repository.go

  internal/auth/
    service.go
    repository.go

  internal/postgres/
    catalog.go
    progress.go
    auth.go
    migrations/

  internal/httpapi/
    catalog.go
    progress.go
    auth.go
    middleware.go
    router.go
```

Uma alternativa igualmente válida é colocar os adapters abaixo da funcionalidade, por exemplo `internal/catalog/postgres`. O ponto essencial é manter uma direção de dependência consistente e eliminar a abstração genérica de SQL.

### Princípios propostos

- Interfaces são definidas por quem as consome.
- Interfaces descrevem operações semânticas, não APIs genéricas de driver.
- O adapter PostgreSQL usa `pgx` diretamente.
- Funções simples não precisam virar structs com `Execute`.
- Serviços são usados quando existe comportamento de negócio ou orquestração.
- Transações pertencem ao limite da operação de negócio.
- DTOs HTTP são separados quando o contrato público difere do modelo interno.
- PostgreSQL é a primeira escolha; Redis só entra com necessidade comprovada.
- Worker separado só existe quando houver jobs concretos e política de retry definida.
- Não há microserviço sem necessidade independente de escala, segurança, disponibilidade ou ownership.

## Estratégia de sincronização proposta

Antes de codificar, definir uma especificação com:

```text
SyncOperation
  operation_id
  user_id
  device_id
  device_sequence
  operation_type
  payload
  client_created_at
  received_at
  protocol_version
```

Garantias mínimas:

- unicidade de `(user_id, device_id, device_sequence)`;
- `operation_id` globalmente único;
- processamento e projeção na mesma transação;
- resposta com operações aceitas, rejeitadas e cursor do servidor;
- lotes limitados;
- retries idempotentes;
- nenhuma decisão baseada apenas no relógio do cliente;
- projeções separadas para posição atual, maior avanço e capítulos concluídos;
- política explícita de reset e correção;
- retenção ou compactação do log;
- métricas de lag, conflitos, retries e falhas.

Essa abordagem preserva as vantagens dos eventos sem adotar event sourcing como arquitetura global.

## Estratégia de escala proposta

### Fase inicial

- Uma API Go stateless.
- Um PostgreSQL com pool, backups e métricas.
- Catálogo com ETag e `Cache-Control`.
- Corpus estático ou projeção publicada atomicamente.
- Inbox idempotente para sincronização.
- Jobs simples no PostgreSQL, se necessários.
- Sem Redis até existir carga ou requisito que o justifique.

### Escala horizontal

- Múltiplas réplicas da API atrás de load balancer.
- Readiness vinculada à capacidade de atender tráfego.
- Shutdown e drenagem de conexões.
- Pool dimensionado por réplica e pelo limite total do PostgreSQL.
- Operações idempotentes e transacionais.
- Métricas e alertas antes de otimizações.

### Evolução orientada por evidências

- CDN para conteúdo público e imutável.
- Redis para rate limit ou cache somente se PostgreSQL/CDN não atenderem.
- Particionamento da tabela de eventos após volume e padrão de acesso conhecidos.
- Fila externa quando jobs em PostgreSQL não atenderem throughput ou isolamento.
- Serviços separados apenas por necessidade operacional ou organizacional comprovada.

## Decisões solicitadas ao Tech Lead

| ID | Decisão | Opção recomendada |
| :-- | :-- | :-- |
| D1 | Modelo arquitetural do backend | Monólito modular com ports pequenas, sem aplicação dogmática de Clean Architecture |
| D2 | Abstração de persistência | Remover `sql.Executor`; assumir PostgreSQL e usar `pgx` no adapter |
| D3 | Localização do repositório SQL | Adapter PostgreSQL fora do núcleo da funcionalidade, ou vertical slice assumido explicitamente |
| D4 | Padrão de casos de uso | Não exigir struct + `Execute` para operações pass-through |
| D5 | Autoridade do catálogo | Definir se será artefato estático ou projeção PostgreSQL versionada e atômica |
| D6 | Modelo de progresso | Separar posição atual, maior avanço e conclusão; formalizar operações de sync |
| D7 | Estratégia de eventos | Inbox idempotente + projeções, sem event sourcing global |
| D8 | Redis e worker | Adiar até existir caso de uso e requisito operacional concreto |
| D9 | Contrato da API | Adotar OpenAPI ou contrato formal equivalente antes de múltiplos clientes |
| D10 | Qualidade mínima | CI obrigatório com build, vet, testes, race e integração PostgreSQL |
| D11 | Operação | Readiness, shutdown, logs, métricas e processo de migration antes de produção |
| D12 | Corpus | Resolver licença, proveniência, versionamento e publicação atômica |

## Plano recomendado

### Prioridade 0: restaurar uma baseline executável

1. Corrigir o composition root e restaurar `go test ./...`.
2. Criar o pool PostgreSQL e executar `Ping` com timeout no boot.
3. Injetar repositório e operações do catálogo.
4. Implementar shutdown gracioso e fechamento do pool.
5. Criar CI da API.

### Prioridade 1: corrigir fundações do catálogo

1. Decidir sobre a remoção de `sql.Executor`.
2. Corrigir `BookName` e fechar `Rows` explicitamente.
3. Endurecer schema e validação HTTP.
4. Implementar testes de integração PostgreSQL.
5. Implementar migration e seed reproduzíveis.
6. Publicar o corpus de forma atômica e registrar sua versão.

### Prioridade 2: definir contratos críticos

1. Especificar a semântica completa de progresso.
2. Especificar protocolo, idempotência e reconciliação de sync.
3. Especificar autenticação, sessão e threat model.
4. Formalizar contrato HTTP para web e mobile.
5. Resolver licença e proveniência do corpus.

### Prioridade 3: preparar operação

1. Separar liveness e readiness.
2. Adicionar request ID, recuperação de panic e logs estruturados.
3. Adicionar métricas de HTTP, PostgreSQL e sincronização.
4. Definir SLOs e alertas mínimos.
5. Testar migrations, rollback, concorrência e falhas parciais.

### Prioridade 4: evoluir somente por necessidade

1. Implementar worker quando houver jobs reais.
2. Introduzir S3 quando compartilhamento de imagens for implementado.
3. Introduzir FCM se lembretes remotos forem uma decisão de produto.
4. Introduzir Redis após medição de carga ou requisito distribuído concreto.
5. Considerar particionamento ou fila externa somente com evidência operacional.

## Verificações realizadas

| Verificação | Resultado |
| :-- | :-- |
| `api: go test ./...` | Falhou na compilação em `cmd/api/main.go:22` |
| `bible/tools: go test ./...` | Passou |
| `tools/vaultlint: go test ./...` | Passou |

## Limitações da avaliação

- A análise representa o working tree local em 5 de agosto de 2026, que continha mudanças modificadas e arquivos ainda não rastreados.
- A falha de compilação pode representar trabalho em andamento, mas continua relevante pela ausência de CI da API.
- Web, mobile, worker, autenticação e sincronização ainda não possuem implementação suficiente para revisão de código.
- Não foram realizados testes de carga, segurança ofensiva ou recuperação de desastre.
- A análise de licenciamento identifica ausência de documentação; não constitui parecer jurídico.

## Referências principais

- `Obsidian Vault/PRD.md`
- `Obsidian Vault/Rules/01 - Princípios de Engenharia.md`
- `Obsidian Vault/Rules/02 - Segurança.md`
- `Obsidian Vault/Docs/Architecture/Visão Geral.md`
- `Obsidian Vault/Docs/Architecture/Autenticação e Sessões.md`
- `Obsidian Vault/Docs/Architecture/Sincronização Offline.md`
- `Obsidian Vault/Docs/Decisions/002 - Progresso por Eventos.md`
- `Obsidian Vault/Plans/Active/02 - Catálogo Bíblico.md`
- `api/cmd/api/main.go`
- `api/internal/catalog/`
- `api/internal/adapters/sql/`
- `api/internal/adapters/postgres/`
- `api/internal/transport/httpapi/`
- `infra/docker-compose.yml`
