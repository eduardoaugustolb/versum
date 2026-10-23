# Worker da outbox

`cmd/worker/main.go` compõe PostgreSQL, criptografia, sender Resend, handler de
magic link e processador. O loop serial fica em `internal/outboxevent/adapters/polling`;
SQL fica no DAO e a conversão de payload cifrado para evento no repositório de entrega.

Além da configuração compartilhada de `.env.example`, configure `RESEND_API_KEY`,
`EMAIL_FROM` (remetente verificado) e `MAGIC_LINK_ACCESS_URL` (página HTTPS de
confirmação, sem query ou fragmento). `.env` é opcional quando o ambiente já fornece
as variáveis. O carregador compartilhado ainda exige `REDIS_URL`, embora este
processo não abra uma conexão Redis. A URL de confirmação e o consumo HTTP precisam
existir no cliente/API antes de testar acesso real; o worker não cria essas rotas.

Execute `go run ./cmd/worker`. As migrations existentes precisam estar aplicadas.
O processo não altera schema automaticamente. SIGINT/SIGTERM cancelam polling e
operações em andamento; o pool é fechado e entregas sem confirmação são recuperadas
após a expiração da lease.

A aquisição usa uma única instrução com `FOR UPDATE SKIP LOCKED` e `UPDATE RETURNING`,
em autocommit. Não passe um executor transacional ao DAO de entrega: nenhum lock
deve ficar aberto durante envio. Cada claim emite UUID de lease e incrementa attempts.
Confirmação, retry e falha exigem a lease correspondente ainda válida.

Configuração inicial: lease de 60 segundos, timeout do handler de 20 segundos,
retry após 30 segundos, cinco tentativas, polling ocioso de um segundo. Cada
instância processa um evento por vez; múltiplas instâncias compartilham as leases.
O retry tem intervalo fixo. Confirmação e falha terminal apagam o payload cifrado,
preservando os metadados. Esses metadados ainda precisam de política/job de retenção.

Erros de processamento são registrados por código fixo, sem erro bruto do provedor,
e-mail, token ou URL. Payload indecifrável mantém ciphertext/lease para recuperação
operacional das chaves, sem envio; pode voltar a ser adquirido após a lease vencer.
Após reparar as chaves, o limite de tentativas do processador continua aplicável.

Entrega é pelo menos uma vez: mantenha a chave de idempotência e o conteúdo estáveis.
O envio pode ter chegado ao provedor quando o processo é interrompido. Não há uma
transação distribuída entre banco e Resend nem bloqueio de envio concorrente à
exclusão do usuário. As limitações de segurança do fluxo estão em
`auth-known-limitations.md`.
