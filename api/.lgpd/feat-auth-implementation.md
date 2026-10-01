# Implementação necessária na feat/auth

Data: 2026-10-01. Plano de execução, não implementação realizada. Branch conferida: `feat/auth`. Cobre autenticação e seu ciclo de privacidade; progresso, push, analytics e clientes têm planos próprios. Checklist corresponde aos requisitos do projeto; mecanismos técnicos escolhidos não são exigências literais universais da LGPD.

## 1. Fechar o contrato de acesso

- [x] Definir base/contrato de acesso: art. 7º, V para conta gratuita opcional, pedido/entrega/sessão de dados comuns. Decisão em [access-contract.md](access-contract.md). Responsável/contato e revisão dos textos pendentes; sem consentimento obrigatório de conta.
- [ ] Implementar pedido informado com versão dos termos e confirmar controle do e-mail antes de ativar conta; identidade provisória/pedido deve expirar e ser limpo. Leitura pública continua sem login.
- [ ] Definir criação automática vs conta existente, invalidação dos links anteriores, retry e retenção efetivos. `InvalidatePreviousLinks` definido no domínio não demonstra invalidação conectada.
- [ ] Disponibilizar aviso versionado antes da coleta; manter texto recuperável. Ciência da política/aceite contratual distintos dos consentimentos opcionais. Não pedir aceite para módulo ainda desativado.

Aceite: contrato/API/docs concordam; recusa opcional não bloqueia leitura pública ou conta. LGPD arts. 6º–11.

## 2. Tornar o magic link entregável

- [ ] Incluir token bruto apenas dentro do payload cifrado da outbox, junto dos vínculos necessários; persistir somente hash em `login_tokens`. Hoje payload guarda apenas IDs, logo não permite envio.
- [ ] Manter usuário/token/evento na mesma unidade de trabalho já implementada; rollback não deixa entrega solta.
- [ ] Implementar claim/lease, envio fora da transação, confirmação condicional pela lease, retry/backoff/limite, cancelamento e recuperação de crash. Provedor com idempotência quando disponível.
- [ ] Não enviar token expirado/consumido ou pertencente a conta excluída. Testar interleaving entre claim, exclusão e envio; se provedor já recebeu mensagem, registrar a limitação e garantir que o link se torne inválido.
- [ ] Remover segredo após uso, expiração ou cancelamento e separar metadados de auditoria com retenção própria. Não guardar token cifrado em outbox por 30–90 dias indiscriminadamente.

Aceite: teste com transportador sintético monta link válido; evento/token atômicos; retry duplicado não permite segundo consumo; segredo ausente nos logs. LGPD arts. 6º, III/VII/VIII, e 46.

## 3. Consumo e sessão completos

- [ ] Implementar `ConsumeMagicLink` e ajustar o teste em desenvolvimento, sem simplesmente removê-lo para obter suíte verde.
- [ ] Consumo SQL condicional (`consumed_at IS NULL`, `expires_at > now`) e criação de sessão na mesma transação; falha de sessão reverte consumo.
- [ ] Segredo de sessão aleatório, apenas hash persistido, TTL efetivo; definir dado mínimo de dispositivo, sem fingerprint desnecessário.
- [ ] Implementar autenticar sessão, consultar sessão atual, logout individual/global, expiração e revogação por dispositivo. Autorizar sempre pelo servidor, não só Redis.
- [ ] Expor troca de token e cookies web com `HttpOnly`, `Secure` em produção, `SameSite` coerente, escopo restrito, proteção CSRF quando aplicável; allowlist de redirects e `Referrer-Policy: no-referrer`.
- [ ] Remover token de URL/histórico no cliente; não capturar URL em proxy/APM. Catálogo segue público. Cliente Android fica em plano próprio.

Aceite: dois consumidores simultâneos produzem uma sessão; token expirado/reutilizado rejeitado; logout invalida sessão; isolamento entre titulares/dispositivos. LGPD art. 46 e regra de segurança do projeto.

## 4. Entrada, rate limit e logs

- [ ] Limite de body antes do decode, rejeição de JSON concatenado, parser de media type, erros públicos neutros.
- [ ] HMAC do e-mail com segredo separado, versionamento/rotação e namespace; não reutilizar indiscriminadamente o índice cego do banco em Redis. Documentar tratamento do IP e TTL.
- [ ] Resolver proxy confiável; não confiar em `X-Forwarded-For` de qualquer cliente. Manter auth fail closed.
- [ ] Redigir erros das dependências antes de logs; não logar token, cookie, e-mail, URL ou progresso. Trilhas administrativas mínimas sob acesso restrito, sem conteúdo pessoal desnecessário.

Aceite: testes de body grande/duplicado, charset, limites Redis, falha de cache e entradas sintéticas com segredos; sem PII em logs/erros. LGPD arts. 6º, III/VII/VIII/X, 12 e 46.

## 5. Direitos do titular

- [ ] `GET /me` e exportação autenticados: dados da conta, sessões/metadados, evidências de consentimento quando existirem e informação compreensível de finalidades/compartilhamentos. Nunca exportar tokens, cookies, hashes de autenticação ou material de chaves.
- [ ] Correção de e-mail com prova do novo endereço e prevenção de tomada de conta; atualizar ciphertext/HMAC atomicamente e decidir revogação de sessões/links antigos.
- [ ] Exclusão autenticada e resistente a replay: revogar acesso, cancelar entregas, remover tokens/sessões/conta/outbox e derivados; não permitir que worker ou retry recrie registro.
- [ ] Canal operacional para confirmação/acesso, correção, oposição quando cabível, portabilidade conforme regulamentação, informação, revogação e eliminação; protocolo, resposta e motivo de eventual conservação.

Aceite: titular A não acessa dados de B; exclusão alcança outbox sem depender de FK inexistente e não entrega link utilizável. Prazo de declaração completa de confirmação/acesso: art. 19, II; demais fluxos têm prazo aplicável definido, sem alegar regra universal de 15 dias. LGPD arts. 18–19.

## 6. Retenção, purge e chaves

- [ ] Aprovar/configurar regras por entidade e estado, gatilho e exceção; ver [retention.md](retention.md). Não confundir TTL com retenção física.
- [ ] Job em lotes, idempotente, com monitoramento e recuperação; integra tokens/sessões expirados/revogados e outbox terminada/cancelada.
- [ ] Definir restauração de backups que reaplique exclusões antes de reabrir tráfego. Chave global não permite crypto-shredding individual por conta.
- [ ] Recifragem idempotente de e-mail/outbox e rotação do lookup; não retirar chave enquanto dados, backups ou material de entrega dependem dela.

Aceite: relógio controlado + integração demonstram purge no limite aprovado, exceções e recuperação após falha; rotação mantém acesso sem segredo em logs. LGPD arts. 15–16 e 46.

## 7. Consentimento, se ativado nesta branch

- [ ] `policy_versions` conserva texto/hash/versão; ledger aponta para texto específico da finalidade e inclui `occurred_at`, ação e desempate determinístico.
- [ ] Concessão/revogação autenticadas, idempotentes e isoladas; servidor determina titular. Não ordenar revogações somente por `granted_at`.
- [ ] Estado consentido impede novos tratamentos após revogação; descarte limitado por regra aprovada de prova. Hash de IP e UA não são coleta obrigatória.

Se nenhuma finalidade consentida for ativada, documentar esse fato e manter feature desativada; infraestrutura de consentimento não transforma base contratual em consentimento. LGPD arts. 8º, 11 e 18, IX.

## 8. Fechamento da branch e condições de lançamento

- [ ] `go test ./...`, `go vet ./...`, integração realmente executada contra PostgreSQL/Redis de teste, concorrência e crash recovery. Resultados cached/skipped não bastam.
- [ ] Docs alinhadas ao comportamento entregue; aviso corresponde às funções realmente ativadas; atualizar mapa/ROPA por mudança.
- [ ] Controlador/contato, bases e retenção confirmados; política preenchida e revisada antes de publicação. Fornecedores e regiões, instruções contratuais e mecanismo válido internacional documentados.
- [ ] TLS, segredos, acesso mínimo, backups/restauração e runbook de incidentes demonstrados no deploy.
- [ ] Triagem de menores concluída antes de convites reais. O art. 9º do ECA Digital não impõe verificação a cada acesso de qualquer app; avaliar o alcance real, vínculo responsável e minimização.

Deploy e publicação jurídica são condições de lançamento, não evidências produzidas por fechar a branch. Progresso/push/analytics e mecanismos de menores ficam bloqueados até seus desenhos específicos, quando necessários. LGPD arts. 9º, 14, 33–39, 41, 46, 48; Lei 15.211/2025 arts. 1º–2º, 8º–9º e 24.
