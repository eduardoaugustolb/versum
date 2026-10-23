# Limitações conhecidas do fluxo de magic link

Atualizado em 2026-10-22. Solicitação, persistência, publicação transacional, handler de entrega e worker PostgreSQL existem; integração HTTP do consumo ainda não está concluída. Não houve envio real validado no Gmail nem comprovação de prontidão para produção. O roteiro de implementação está em [feat-auth-implementation.md](../.lgpd/feat-auth-implementation.md).

## Entrega e consumo ainda incompletos

O payload da solicitação contém `login_token_id` e o token bruto, cifrados pelo adaptador da outbox; `user_id` fica na coluna de vínculo. `login_tokens` continua persistindo somente o hash. O handler consulta pelo segredo, confere ID e titular, ignora tokens consumidos/expirados e registros excluídos, monta HTML/texto e propaga falhas do sender. O [worker](outbox-worker.md) adquire e confirma por lease no PostgreSQL, com retry limitado e cancelamento por sinal. Conclusão e falha terminal apagam o payload; consumo/expiração sem processamento pelo worker e retenção de metadados ainda precisam de limpeza própria.

O construtor do handler exige um `clock.Clock` e a URL HTTPS da página web de confirmação, sem query, fragmento ou credenciais. O handler adiciona `?token=...` com encoding; o destino deve evitar logs/captura dessa URL, remover o segredo do histórico e exigir ação explícita para consumir o link, evitando consumo por scanners de e-mail. Essa página ainda não é implementada pela API.

A chave de envio é `magic-link/<event_id>`. HTML/texto usam expiração absoluta e permanecem iguais nos retries com as mesmas dependências/configuração. Conforme a [documentação do Resend](https://resend.com/docs/dashboard/emails/idempotency-keys), a deduplicação dura 24 horas e exige o mesmo payload. Não alterar template, URL, remetente ou destinatário de uma entrega pendente sob a mesma chave; um rollout precisa drenar ou preservar a versão do material de entrega. A chave não substitui confirmação da outbox nem garante entrega exatamente uma vez.

As verificações anteriores ao envio não são uma barreira atômica contra exclusão/consumo concorrentes. Um provedor já pode ter recebido o e-mail quando a conta for excluída; o consumo deve continuar rejeitando o link. Falta demonstrar esses interleavings em integração e sanitizar erros do provedor antes de registrar logs.

Os casos de uso `ConsumeMagicLink` e `RotateSession` existem, com emissão transacional de sessão, rotação e revogação da família ao detectar reutilização. Ainda não há integração HTTP desses casos de uso nem ciclo de cookies de sessão completo. IP e user-agent são metadados de risco, não uma prova da identidade do portador do token.

## Cobertura de concorrência e de bordas HTTP

A suíte não executa duas solicitações simultâneas para o mesmo e-mail contra o PostgreSQL. O fluxo usa busca e criação condicional do usuário dentro de uma transação, portanto esse cenário deve ser coberto para garantir que ambas as solicitações criem tokens para o mesmo usuário e que nenhum erro de conflito ou de transação escape para o cliente.

Também faltam testes para corpos JSON concatenados, `Content-Type: application/json; charset=utf-8`, corpos excessivamente grandes e falhas da dependência do caso de uso. Esses casos devem confirmar a resposta HTTP adequada e que e-mail e token não sejam incluídos em respostas ou logs.

## Rotação da chave de criptografia de e-mail

Quando encontra um usuário por uma chave de lookup antiga, a aplicação atualiza o HMAC e a versão de lookup. Ela não volta a cifrar o e-mail com a chave AES-GCM atual. Assim, a chave de criptografia antiga continua necessária enquanto existirem esses registros.

Antes de remover uma chave de `ENCRYPTION_SECRET_PREVIOUS_KEYS`, é necessário executar uma migração que leia cada e-mail com sua versão atual, recifre com a chave corrente e atualize `email_ciphertext` e `email_encryption_key_version` de forma transacional. A migração precisa ser idempotente, observável e coberta por testes com chaves em versões distintas.

## Proteção contra abuso no endpoint

`POST /auth/magic-link` já limita solicitações por IP e e-mail via Redis: 10 por minuto por chave, com falha fechada quando o cache de autenticação está indisponível. O limite global é 50 por minuto por IP e pode degradar aberto. O identificador de e-mail usa SHA-256 sem segredo; substituir por HMAC sob chave separada para reduzir ataques por dicionário. IP e identificadores derivados continuam dados pessoais vinculáveis.

Ainda falta limitar o corpo antes de decodificar JSON, rejeitar documentos concatenados e tratar media type JSON com charset. Confirmar endereço do cliente atrás de proxy confiável sem aceitar cabeçalhos forjados. A resposta deve continuar uniforme para usuários existentes e inexistentes, evitando enumeração de contas. TTL operacional de um minuto não demonstra descarte de backups/dumps de Redis.

## Retenção, direitos e operação

A política de emissão expõe somente `TTL`. Retenção, invalidação de links anteriores e tentativas de entrega não são configurações ativas; dependem de casos de uso e workers próprios. Expiração lógica não elimina registros. `DeleteUser` revalida a sessão após bloquear o usuário na mesma transação da exclusão; as FKs eliminam tokens, sessões e eventos vinculados. A integração HTTP da exclusão, limpeza por estado, descarte de cache, barreira contra entrega concorrente, exportação/correção e canal para demais direitos continuam pendentes. Logs de erro precisam de sanitização demonstrada. Recifragem também deve cobrir outbox e dependências de backups antes de retirar chave antiga.

Base do acesso comum definida: art. 7º, V para acordo gratuito de conta opcional; ver [access-contract.md](../.lgpd/access-contract.md). Ajustar criação atual para confirmar controle do e-mail antes de ativar conta. Controlador/contato, retenção, fornecedores/regiões e controles produtivos ainda exigem definição. Ciência do aviso não substitui consentimento; não adicionar gate obrigatório de `conta_acesso` por padrão. Fundamentos e pendências estão em [.lgpd/STATUS.md](../.lgpd/STATUS.md).
