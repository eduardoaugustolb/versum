---
title: "01 - Fundação da API Go"
section: Plans
subsection: Active
type: implementation-plan
status: approved
date: 2026-08-02
tags: [versum, plans, api, go]
up: "[[Plans/Active/_Index|Planos Ativos]]"
prev: "[[Plans/Active/_Index|Planos Ativos]]"
next: "[[Plans/Archive/_Index|Arquivo]]"
related: ["[[Docs/Architecture/Visão Geral]]", "[[Rules/01 - Princípios de Engenharia]]"]
---

# Fundação da API Go

## Objetivo

Criar uma API Go executável, configurável e testada, com uma rota de saúde que
possa ser usada pelo ambiente de desenvolvimento e pelo deploy. Esta entrega
estabelece o formato dos módulos e do transporte HTTP, sem antecipar banco,
autenticação ou conteúdo bíblico.

## Arquitetura

O domínio `health` expõe um caso de uso pequeno. O adapter HTTP chama esse caso
de uso e o composition root em `cmd/api` monta router, configuração e servidor.
O domínio não importa `net/http` ou `chi`.

## Stack

- Go 1.26.5
- `net/http` e `github.com/go-chi/chi/v5`
- `log/slog`
- `go test`

## Critérios de aceitação

1. `go test ./...` passa no diretório `api/`.
2. `go run ./cmd/api` inicia um servidor em `:8080` por padrão.
3. `GET /health` responde `200 OK` com JSON contendo `status: "ok"`.
4. `PORT` aceita somente uma porta TCP válida e altera o endereço de escuta.
5. O handler HTTP não contém regra de domínio.

## Estrutura de arquivos

```text
api/
  go.mod
  cmd/api/main.go
  internal/config/config.go
  internal/config/config_test.go
  internal/health/service.go
  internal/health/service_test.go
  internal/transport/httpapi/router.go
  internal/transport/httpapi/router_test.go
```

## Etapas

### 1. Configuração da aplicação

- Criar `api/go.mod` com módulo `github.com/eduardoaugustolb/versum/api`, Go
  1.26.5 e a dependência `github.com/go-chi/chi/v5`.
- Escrever `config_test.go` antes de `config.go` para exigir:

  ```go
  type Config struct {
      Environment string
      Address     string
  }

  func Load(lookup func(string) string) (Config, error)
  ```

- Cobrir valores padrão (`development`, `:8080`), porta válida (`PORT=9090`) e
  porta inválida (`PORT=abc` ou `PORT=70000`).
- Implementar o mínimo para os testes passarem e executar `go test ./internal/config`.
- Commit sugerido: `feat(api): add runtime configuration`.

### 2. Caso de uso de saúde

- Escrever `service_test.go` antes de `service.go` para exigir:

  ```go
  type Status struct {
      State string
  }

  type Service struct{}

  func (Service) Check() Status
  ```

- Verificar que `Check()` devolve `Status{State: "ok"}`.
- Implementar a struct e o método sem dependências externas.
- Executar `go test ./internal/health`.
- Commit sugerido: `feat(api): add health use case`.

### 3. Adapter HTTP e router

- Escrever `router_test.go` antes de criar os handlers. O teste chama
  `NewRouter(health.Service{})` com `httptest.NewRecorder()` e exige:

  ```text
  GET /health
  status: 200
  Content-Type: application/json
  body: {"status":"ok"}
  ```

- Implementar `NewRouter(service health.Service) http.Handler` em
  `internal/transport/httpapi/router.go` usando `chi.NewRouter()`.
- Codificar a resposta apenas no adapter:

  ```go
  type healthResponse struct {
      Status string `json:"status"`
  }
  ```

- Executar `go test ./internal/transport/httpapi` e depois `go test ./...`.
- Commit sugerido: `feat(api): expose health endpoint`.

### 4. Composition root e execução local

- Criar `cmd/api/main.go` depois de os testes anteriores estarem verdes.
- Carregar `config.Load(os.Getenv)`, registrar falha de configuração com `slog`
  e encerrar com código 1.
- Montar `http.Server{Addr: cfg.Address, Handler: httpapi.NewRouter(health.Service{})}`.
- Tratar `http.ErrServerClosed` como encerramento esperado e registrar o endereço
  de escuta no início.
- Executar `go run ./cmd/api` e verificar `curl -i http://localhost:8080/health`.
- Executar `go test ./...` e `go vet ./...`.
- Commit sugerido: `feat(api): add server entrypoint`.

## Fora deste plano

- Banco de dados, Redis, S3, seed bíblico e Discord.
- Magic link e sessões.
- Docker, CI e qualquer cliente web ou mobile.

Essas partes dependem de decisões e contratos que merecem entregas próprias.

---

◀ [[Plans/Active/_Index|Planos Ativos]] · próxima: [[Plans/Archive/_Index|Arquivo]] ▶
