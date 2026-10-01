# Retenção e eliminação — especificação para feat/auth

Versão 0.1, 2026-10-01. Proposta a operacionalizar; não declara purge existente. Base: LGPD arts. 6º, III, 15–16, 18 e 46. Prazos de produto não são prazos impostos pela LGPD.

| Dado | Gatilho e ação | Prazo e decisão pendente |
|---|---|---|
| Conta/e-mail | Enquanto necessário ao acesso; no encerramento, eliminar e cancelar dependentes; conservar exceção mínima apenas se fundamentada | Critério de inatividade e eventual conservação probatória pendentes do controlador; não reter conta por prazo genérico de prescrição |
| Login token | Expiração, consumo ou invalidação impede uso; remover linha por purge | Proposta: 24h após o primeiro estado terminal; o parâmetro atual deve ser conectado e seu gatilho explicitado |
| Segredo cifrado de entrega | Remover após consumo/expiração/cancelamento e, se já não necessário, após entrega confirmada | Imediatamente ao estado terminal, por transação ou job com janela máxima aprovada; sem retry depois do prazo de validade |
| Sessão | Revogação/expiração impede acesso; remover metadados não necessários | Janela pós-terminal ainda precisa ser fixada/configurada; retenção ilimitada não é default aceitável |
| Outbox/metadados | Excluir segredo separado do eventual registro de entrega/erro mínimo | Fixar duração e finalidade de observabilidade; eliminar faixa arbitrária 30–90 dias como regra universal |
| Redis rate limit | Chave expira ao fim da janela | 1 minuto é configuração atual; definir acessos e persistência/dumps/backups para não prolongar sem justificativa |
| Logs operacionais | Não armazenar segredos ou PII desnecessária; descarte automático no destino | Duração/destino a confirmar antes de deploy; logs de acesso eventualmente exigidos por outra norma precisam de avaliação separada |
| Prova de consentimento/DSAR | Restringir acesso, cessar uso original, conservar apenas prova necessária com base própria e descartar ao término | Prazo fundamentado pendente; append-only durante retenção não significa eternidade |
| Backups | Expirar segundo política; reaplicar exclusões ao restaurar antes de acesso/entrega | Janela de backup e procedimento de restauração pendentes de infraestrutura; não prometer remoção seletiva que não existe |
| Registro de incidente | Manter informações exigidas, com acesso restrito | Mínimo de 5 anos a partir do registro, Res. 15/2024 art. 10; evidências técnicas anexas requerem minimização própria |

Implementação: política por entidade/estado, relógio UTC, lotes idempotentes e retry, métricas sem PII, trilha mínima de ação e teste de restauração. Eliminação cobre outbox sem FK e entregas concorrentes. Identificador de bloqueio/tombstone, se necessário, também recebe finalidade/base/prazo; não conservar a conta inteira apenas para evitar recriação.

Antes de publicar o aviso, substituir cada pendência por decisão do controlador e alinhar ROPA/configuração. Se outra lei exigir guarda, indicar dispositivo, dado, finalidade, acesso e data de descarte; não importar tabelas genéricas fiscais/CDC para este app sem atividade correspondente.
