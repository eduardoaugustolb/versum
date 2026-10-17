# Limites de arquitetura da API

O domínio mantém entidades e invariantes sem dependências de transporte ou persistência. A aplicação coordena casos de uso e declara portas. Adaptadores implementam essas portas; `cmd/api` compõe as implementações. Os handlers HTTP recebem interfaces pequenas e traduzem os resultados da aplicação em respostas.

## Catálogo

`PublishBook` valida o livro e a composição dos versículos antes de persistir. O domínio rejeita versículos de outro livro, capítulos acima de `ChapterCount` e identificadores duplicados de versículo. `Chapter` valida sua composição e copia os slices de entrada e saída.

`PublishCatalog` publica todos os livros e a versão em uma única unidade de trabalho. O seed lê arquivos e converte entradas, mas não determina a fronteira de commit. Os gerenciadores transacionais sempre agendam rollback, inclusive para panic ou falha no commit.

## Identidade e acesso

Exclusão de usuário, logout e rotação usam uma transação e o bloqueio do usuário definido por `LockSessionsByUserID`. Após adquirir o bloqueio, os casos de uso releem a sessão. Exclusão exige uma sessão válida no instante fornecido por `clock.Clock`; logout preserva a possibilidade de revogar a família com uma credencial substituída e ainda não expirada.

`RequestMagicLink` exige uma porta `EmailRequestLimiter`. A cota por e-mail é verificada antes da transação, inclusive fora do HTTP. O adaptador de cache implementa a cota; o HTTP mantém o limite por IP e traduz cota esgotada para 429 e indisponibilidade para 503. A política de emissão configura apenas o TTL efetivamente implementado. Entrega, retenção e invalidação de links anteriores continuam pendentes, conforme `auth-known-limitations.md`.

O nome do evento `identityaccess.magic_link_requested` pertence ao contrato de integração da aplicação de identidade. A outbox aceita tipos genéricos e não conhece nomes de eventos dos módulos publicadores.

## Verificação

Execute `go test ./...` e `go vet ./...`. Testes em `integration` exigem `DATABASE_URL` apontando para um banco descartável com as migrations aplicadas. Eles alteram dados e devem executar com `go test -count=1 -p 1 ./integration/...`, pois os pacotes compartilham tabelas. Há testes de rollback de conteúdo e versão e de revogação concorrente durante uma exclusão.
