# Avaliação Arquitetural do Versum

**Data da avaliação:** 7 de setembro de 2026
**Escopo:** branch `feat/auth`, documentação do `Obsidian Vault/`, implementação da API Go, corpus bíblico, migrations, testes e automações existentes.
**Objetivo:** fornecer evidências e opções para que a liderança técnica decida como simplificar a arquitetura, corrigir inconsistências e preparar o produto para evolução e escala.

## Resumo executivo

A stack e a arquitetura macro são adequadas ao produto. Go, PostgreSQL, `pgx`, `chi`, SQLite no mobile e uma API stateless formam uma base simples e capaz de atender uma escala significativa sem microserviços.

O principal problema não é a escolha das tecnologias. É a distribuição da complexidade:

- o projeto ainda precisa controlar a quantidade de abstrações à medida que novos
  módulos entram no monólito;
- sincronização, observabilidade e operação continuam incompletas;
- autenticação está em planejamento implementável, com as migrations iniciais
  criadas, mas ainda não possui handlers e casos de uso;
- o catálogo está implementado, testado e separado como entrega concluída;
- a API já possui CI próprio, embora a cobertura operacional ainda possa crescer.

A recomendação principal é manter um **monólito modular Go**, assumir PostgreSQL explicitamente, usar interfaces semânticas pequenas definidas pelos consumidores e remover abstrações de infraestrutura sem um segundo caso real. Antes de expandir o backend, devem ser formalizados os contratos de sincronização e autenticação.

## Conclusão geral

| Dimensão | Avaliação | Síntese |
| :-- | :-- | :-- |
| Stack | Adequada | Go, PostgreSQL, `pgx`, `chi`, Expo/React Native e SQLite são escolhas coerentes. |
| Arquitetura macro | Adequada | Monorepo e monólito modular são suficientes; não há justificativa para microserviços. |
| Idiomatismo Go | Em evolução | A porta `dbexec` é pequena e neutra; novas abstrações devem continuar sendo justificadas por casos de uso reais. |
| Limites de domínio | Coerentes no catálogo | O repositório usa `dbexec`, enquanto o adapter concreto PostgreSQL fica isolado em `internal/adapters/postgres`. |
| Escalabilidade | Potencialmente boa, ainda não demonstrada | A API pode ser stateless, porém faltam readiness, observabilidade, estratégia de eventos e processo de deploy. |
| Sincronização | Risco alto | O diferencial central do produto ainda não possui semântica completa nem modelo operacional. |
| Segurança | Em implementação | O contrato de autenticação, cifragem de dados pessoais e migrations estão documentados; o fluxo HTTP ainda falta. |
| Testes e entrega | Baseline funcional | `go test ./...`, `go vet ./...`, race e CI da API existem; faltam ampliar verificações operacionais. |
| Corpus | Implementado | Seed transacional por livro, publicação de versão e hash do manifesto estão implementados; licenciamento continua pendente. |

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

### 1. A portabilidade de banco não deve ser um objetivo implícito

**Severidade:** alta  
**Categoria:** idiomatismo, complexidade e limites arquiteturais

O catálogo agora assume PostgreSQL explicitamente. A antiga abstração genérica
de SQL foi substituída por `internal/ports/dbexec`, uma porta pequena sem tipos
de `pgx`; o adapter `postgres.PgxExecutor` conhece o driver concreto.

Essa afirmação não se sustenta: outro banco exigiria mudanças nas queries, nos tipos, na semântica de transação e possivelmente no tratamento de erros.

A porta não promete trocar de banco sem alterar queries ou semântica:

```go
type Executor interface {
    QueryRow(ctx context.Context, query string, args ...any) Row
    Query(ctx context.Context, query string, args ...any) (Rows, error)
    Exec(ctx context.Context, query string, args ...any) error
}
```

O risco restante é manter a porta genérica maior do que os casos de uso exigem.
As queries continuam PostgreSQL-específicas por usarem `$1`, `ON CONFLICT` e
`COPY`.

Consequências da decisão atual:

- mais interfaces e adapters para manter;
- o adapter pode ser substituído em testes sem vazar `pgx` para a aplicação;
- não há promessa de portabilidade ilusória entre bancos.

Evidências:

- `Obsidian Vault/Plans/Archive/02 - Catálogo Bíblico.md:33-53`
- `Obsidian Vault/Docs/Architecture/Visão Geral.md:74-81`
- `api/internal/ports/dbexec/executor.go:1-20`
- `api/internal/adapters/postgres/pgx_executor.go:1-30`

**Estado:** decisão aplicada no catálogo. Só introduzir suporte a outro banco
quando existir um requisito concreto.

### 2. A separação de driver foi corrigida no catálogo

**Severidade:** alta  
**Categoria:** limites e manutenibilidade

As regras determinam que:

- portas são definidas no caso de uso que depende delas;
- adapters externos ficam fora do domínio;
- dependências apontam para dentro;
- SQL é conhecido apenas pelo adapter PostgreSQL.

O catálogo mantém suas queries junto da funcionalidade, mas a execução é feita
pela porta `dbexec`. O domínio e os casos de uso não importam `pgx`; somente o
adapter PostgreSQL conhece o driver.

Evidências:

- `Obsidian Vault/Rules/01 - Princípios de Engenharia.md:16-25`
- `Obsidian Vault/Rules/01 - Princípios de Engenharia.md:39-82`
- `api/internal/catalog/postgres/repository.go:1-20`
- `api/internal/ports/dbexec/executor.go:1-20`

A decisão registrada é um vertical slice pragmático com isolamento do driver:

| Abordagem | Vantagem | Custo |
| :-- | :-- | :-- |
| Vertical slice pragmático | Menos pacotes e indireção; SQL próximo da funcionalidade | `catalog` deixa de ser um domínio puro e passa a incluir infraestrutura |
| Ports & Adapters | Direção de dependência explícita; domínio sem SQL | Mais um pacote e wiring, justificável onde há regra de negócio real |

**Estado:** corrigido para o catálogo. O módulo de autenticação deve seguir a
mesma direção.

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

### 4. A baseline da API foi restaurada

**Severidade:** resolvida; baixa como risco residual
**Categoria:** implementação e entrega

O composition root agora cria o pool PostgreSQL, confirma a conexão com
`Ping`, injeta catálogo e health no router, fecha o pool e trata shutdown
gracioso:

- `api/cmd/api/main.go:22`
- `api/internal/transport/httpapi/router.go:9`
- `api/internal/transport/httpapi/dependencies.go:8-16`

A verificação atual com `go test ./...` passa em todos os pacotes da API.

Evidências: `api/cmd/api/main.go`, `api/internal/transport/httpapi/router.go` e
`.github/workflows/api-ci.yml`.

**Estado:** resolvido para a baseline atual. Readiness, observabilidade e
integração real com um banco de CI continuam como trabalho operacional.

### 5. O contrato de capítulo foi corrigido na persistência

**Severidade:** resolvida; baixa como risco residual
**Categoria:** corretude

`FindChapter` agora faz join com `books`, preenche `BookName`, fecha as rows e
valida o erro após a iteração.

Evidências:

- `api/internal/catalog/domain/chapter.go`
- `api/internal/catalog/postgres/queries.go:5-11`
- `api/internal/catalog/postgres/repository.go:71-108`

Os testes de aplicação e transporte cobrem o contrato; a integração com banco
continua sendo a verificação adicional desejável para o ambiente de CI.

**Estado:** resolvido no repositório do catálogo.

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

### 7. O seed publica uma projeção fiel e atômica

**Severidade:** alta  
**Categoria:** dados e operação

O corpus é a autoridade do conteúdo, enquanto PostgreSQL é uma projeção. O seed
substitui cada livro dentro de uma transação e publica o hash do manifesto em
`catalog_version`.

Problemas:

- a publicação ainda precisa de validação operacional de contagens e hashes no
  ambiente de deploy;
- o seed percorre o corpus inteiro, mas a execução e o rollback devem ser
  cobertos pelo pipeline de integração;
- licenciamento e proveniência continuam pendentes.

Evidências:

- `Obsidian Vault/Plans/Archive/02 - Catálogo Bíblico.md:59-65`
- `Obsidian Vault/Plans/Archive/02 - Catálogo Bíblico.md:92-93`
- `Obsidian Vault/Plans/Archive/02 - Catálogo Bíblico.md:462-474`
- `api/internal/adapters/postgres/migrations/000001_create-catalog.up.sql:1-18`

**Estado:** a estratégia principal foi implementada. Falta automatizar as
validações de publicação no pipeline e documentar licenciamento.

### 8. O schema do catálogo protege as invariantes principais

**Severidade:** média  
**Categoria:** integridade de dados

O schema do catálogo já possui restrições para números positivos, texto não
vazio, testamento válido e referências entre livros e versículos:

```sql
CHECK (chapter_count > 0)
CHECK (chapter > 0)
CHECK (number > 0)
CHECK (part > 0)
CHECK (length(trim(text)) > 0)
```

Pontos que ainda merecem decisão futura:

- não há dimensão de tradução ou versão para múltiplos corpus futuros;
- não existe restrição que relacione o maior capítulo com `chapter_count`.

Evidência:

- `api/internal/adapters/postgres/migrations/000001_create-catalog.up.sql:1-18`

**Estado:** fundações do schema foram endurecidas; validações de publicação
continuam no seed.

### 9. O serviço ainda não está preparado para operação horizontal

**Severidade:** alta antes de produção  
**Categoria:** operação e escalabilidade

Embora os handlers sejam stateless e o servidor já tenha timeouts, `Ping` no
boot e shutdown gracioso, ainda faltam:

- liveness e readiness separadas;
- request ID;
- recuperação e observação de panic;
- logs estruturados por requisição;
- métricas de latência, status e saturação do pool;
- tracing ou correlação entre API e worker;
- métricas de backlog, retries e idade do evento mais antigo;
- SLOs, alertas e dashboards;
- readiness real vinculada ao banco;
- processo reproduzível de migration e seed no deploy.

O health check atual é liveness simples e não representa readiness do banco:

- `api/internal/health/check.go:13-15`

**Recomendação:** separar `/livez` e `/readyz`, conectar readiness ao banco e
adicionar instrumentação mínima antes de múltiplas réplicas ou deploy de
produção.

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

### 11. Autenticação está em implementação

**Severidade:** alta antes da implementação  
**Categoria:** segurança

As decisões iniciais são adequadas: magic link de uso único, hash no banco, cookie `httpOnly`, sessão por dispositivo e redirects allowlisted.

O plano `03 - Autenticação e Sessões` agora define magic link de uso único,
consumo atômico, expiração, sessões revogáveis, allowlist de redirects,
proteção contra enumeração e cifragem de dados pessoais. As migrations
`000003_create-auth` já definem usuários, hashes de tokens, sessões e versões
de chaves.

Ainda falta implementar os casos de uso, handlers, middleware, envio de e-mail,
rate limiting e testes de concorrência. A associação explícita a dispositivo e
o cliente Android ficam para uma etapa futura.

Evidência:

- `Obsidian Vault/Docs/Architecture/Autenticação e Sessões.md:18-29`

**Recomendação:** implementar o plano atual, incluindo threat model, contrato
HTTP, gestão de chaves e testes antes de liberar rotas privadas.

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

### Result sets são fechados explicitamente

`ListBooks` e `FindChapter` executam `defer rows.Close()` e verificam `rows.Err()`
após a iteração:

- `api/internal/catalog/postgres/repository.go:47-69`
- `api/internal/catalog/postgres/repository.go:71-108`

**Estado:** resolvido no repositório do catálogo.

### Validação HTTP é incompleta

Capítulos `0`, negativos e não numéricos são rejeitados no handler antes da
persistência. Ainda faltam limites explícitos para alguns parâmetros públicos.

Também faltam:

- limite e formato de `bookId`;
- contrato padronizado de erros;
- versionamento da API;
- OpenAPI ou contrato equivalente;
- ETag e `Cache-Control` para conteúdo imutável;
- limites de body e headers.

### Testes unitários escondem erros de integração

Os testes unitários, HTTP e de integração cobrem parte do contrato. Ainda é
necessário ampliar a cobertura para:

- queries reais;
- tipos e scans do PostgreSQL;
- migrations `up` e `down` das novas tabelas de autenticação;
- preenchimento de `BookName`;
- ordenação real dos versículos;
- fechamento de rows;
- readiness;
- fluxos completos de autenticação e rotação de chaves.

### PostgreSQL local é publicado em todas as interfaces

`infra/docker-compose.yml` publica `5432:5432` e usa credenciais triviais. Para desenvolvimento local, é mais seguro usar `127.0.0.1:5432:5432`. Ambientes compartilhados devem usar secrets externos e não publicar o banco diretamente.

### CI da API existe, mas precisa ampliar a cobertura operacional

O workflow `.github/workflows/api-ci.yml` cobre `api/**` com build, vet, testes
e teste de race. Hooks locais não substituem CI remoto, e o pipeline ainda pode
ganhar PostgreSQL real para migrations, seed e integração.

O pipeline atual executa:

```text
go build ./...
go vet ./...
go test ./...
go test -race ./...
```

Uma etapa futura deve subir PostgreSQL e executar migrations, seed e testes de
integração.

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

### Prioridade 0: manter a baseline executável

1. Manter `go test ./...`, `go vet ./...` e `go test -race ./...` verdes.
2. Cobrir migrations e seed em um PostgreSQL de CI.
3. Manter composition root, timeouts e shutdown gracioso sob teste.

### Prioridade 1: concluir autenticação

1. Implementar os casos de uso de magic link e sessões.
2. Implementar envio de e-mail, cookies, middleware e allowlist.
3. Integrar cifragem, HMAC cego, rotação de chaves e ciclo de vida LGPD.
4. Cobrir concorrência, expiração, revogação e ausência de segredos em logs.
5. Executar as migrations de auth em integração PostgreSQL.

### Prioridade 2: implementar sincronização

1. Implementar a semântica completa de progresso.
2. Implementar protocolo, idempotência e reconciliação de sync.
3. Formalizar contrato HTTP para web e mobile.
4. Resolver licença e proveniência do corpus.

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
| `api: go test ./...` | Passou em 7 de setembro de 2026 |
| `api: go vet ./...` | Passou em 7 de setembro de 2026 |
| `api: go test -race ./...` | Passou em 7 de setembro de 2026; integração PostgreSQL foi pulada sem `DATABASE_URL` |
| `bible/tools: go test ./...` | Passou na avaliação anterior |
| `tools/vaultlint: go test ./...` | Passou na avaliação anterior |

## Limitações da avaliação

- A análise representa a branch `feat/auth` em 7 de setembro de 2026; as migrations de autenticação estão versionadas no commit `6cd47f5`.
- Web, mobile, worker e sincronização ainda não possuem implementação suficiente para revisão de código.
- Autenticação possui contrato e migrations, mas não fluxo executável completo.
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
- `Obsidian Vault/Plans/Archive/02 - Catálogo Bíblico.md`
- `api/cmd/api/main.go`
- `api/internal/catalog/`
- `api/internal/ports/dbexec/`
- `api/internal/adapters/postgres/`
- `api/internal/transport/httpapi/`
- `infra/docker-compose.yml`
