# Limitações conhecidas do fluxo de magic link

Retrato de `feat/auth` em 2026-10-01. Solicitação, persistência e publicação transacional existem; entrega e integração HTTP do consumo ainda não estão concluídas. Este documento não comprova emissão funcional de links nem prontidão para produção. O roteiro de implementação está em [feat-auth-implementation.md](../.lgpd/feat-auth-implementation.md).

## Entrega e consumo ainda incompletos

O payload cifrado atual contém apenas `user_id` e `login_token_id`. O token bruto gerado é descartado depois do hash; não pode ser reconstruído a partir desse hash. Completar o material cifrado de entrega na mesma transação e implementar o consumidor da outbox antes de afirmar que o link é enviado. Remover o segredo quando consumido, expirado ou cancelado; metadados têm retenção própria.

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

`TokenRetention` é parâmetro, sem purge conectado. Expiração lógica e FK em tokens/sessões não eliminam automaticamente conta ou outbox. Implementar limpeza por estado, exclusão autenticada que alcance filas/cache e barreira contra entrega posterior à exclusão, exportação/correção e canal para demais direitos. Logs de erro precisam de sanitização demonstrada. Recifragem também deve cobrir outbox e dependências de backups antes de retirar chave antiga.

Base do acesso comum definida: art. 7º, V para acordo gratuito de conta opcional; ver [access-contract.md](../.lgpd/access-contract.md). Ajustar criação atual para confirmar controle do e-mail antes de ativar conta. Controlador/contato, retenção, fornecedores/regiões e controles produtivos ainda exigem definição. Ciência do aviso não substitui consentimento; não adicionar gate obrigatório de `conta_acesso` por padrão. Fundamentos e pendências estão em [.lgpd/STATUS.md](../.lgpd/STATUS.md).
