# Limitações conhecidas do fluxo de magic link

Este documento registra pontos conhecidos do fluxo de solicitação de magic link. Eles não impedem a emissão atual de links, mas precisam ser resolvidos antes de ampliar o uso do mecanismo em produção.

## Cobertura de concorrência e de bordas HTTP

A suíte não executa duas solicitações simultâneas para o mesmo e-mail contra o PostgreSQL. O fluxo usa busca e criação condicional do usuário dentro de uma transação, portanto esse cenário deve ser coberto para garantir que ambas as solicitações criem tokens para o mesmo usuário e que nenhum erro de conflito ou de transação escape para o cliente.

Também faltam testes para corpos JSON concatenados, `Content-Type: application/json; charset=utf-8`, corpos excessivamente grandes e falhas da dependência do caso de uso. Esses casos devem confirmar a resposta HTTP adequada e que e-mail e token não sejam incluídos em respostas ou logs.

## Rotação da chave de criptografia de e-mail

Quando encontra um usuário por uma chave de lookup antiga, a aplicação atualiza o HMAC e a versão de lookup. Ela não volta a cifrar o e-mail com a chave AES-GCM atual. Assim, a chave de criptografia antiga continua necessária enquanto existirem esses registros.

Antes de remover uma chave de `ENCRYPTION_SECRET_PREVIOUS_KEYS`, é necessário executar uma migração que leia cada e-mail com sua versão atual, recifre com a chave corrente e atualize `email_ciphertext` e `email_encryption_key_version` de forma transacional. A migração precisa ser idempotente, observável e coberta por testes com chaves em versões distintas.

## Proteção contra abuso no endpoint

`POST /auth/magic-link` não estabelece limite para o corpo da requisição e não limita a frequência por IP ou por e-mail. Um cliente pode criar muitos tokens e eventos de outbox, consumindo banco e capacidade do provedor de e-mail.

O endpoint deve limitar o corpo antes de decodificar JSON e aplicar rate limiting persistente ou distribuído, com chaves separadas para IP e para o identificador de e-mail normalizado. A resposta deve continuar uniforme para usuários existentes e inexistentes, evitando enumeração de contas.
