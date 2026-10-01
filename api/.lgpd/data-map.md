# Mapa de dados — Versum

Versão diagnóstica v1 — 2026-10-01. Cenário A (greenfield/pré-alpha), confirmado pelo usuário. Owner global: pendente. Inventário técnico não é ROPA final nem aprova bases legais. Fontes atuais e limites estão em [code-evidence.md](./code-evidence.md); propostas de bases estão em [legal-basis.md](./legal-basis.md). Produção, titulares reais, fornecedores contratados e regiões não foram verificados.

## Atividades presentes na API

### A001 — Conta e solicitação de acesso

| Campo | Valor |
|---|---|
| Finalidade | Identificar conta pelo e-mail e solicitar magic link |
| Estado | Código executável de solicitação; consumo/entrega ainda ausentes |
| Dados / origem | E-mail fornecido pelo solicitante, UUID gerado, timestamps, hash do token; não há garantia de que solicitante seja dono do e-mail antes de verificação |
| Titulares / menores | Solicitantes/usuários; idade não coletada, menores não podem ser descartados |
| Sensíveis | E-mail não é sensível automaticamente; vinculação a app religioso pode permitir inferências, a avaliar |
| Sistemas | PostgreSQL users e login_tokens; memória do processo |
| Base | Definida para dados comuns: art. 7º, V, conta opcional/acordo gratuito; ver [access-contract.md](access-contract.md). Revisão externa/identificação do controlador pendentes |
| Compartilhamentos / transferência | DB acessado pela API; operador, região e acessos humanos pendentes |
| Retenção | Prazo/critério de conta desconhecido. Token TTL=15 min; parâmetro retenção=24h sem purge implementado |
| Segurança | AES-GCM de e-mail, lookup HMAC, token aleatório de 32 bytes e hash, transação |
| Risco / RIPD | Potencial de inferência religiosa; exposição, escala e público pendentes; não concluir alto risco sem critérios/contexto |
| Owner | Pendente |

### A002 — Proteção antiabuso

| Campo | Valor |
|---|---|
| Finalidade | Limitar frequência de solicitações por IP/e-mail |
| Estado | Implementado |
| Dados / origem | IP observado em RemoteAddr, SHA256 determinístico do e-mail, contadores e TTL |
| Titulares / menores | Visitantes e solicitantes; presença de menores desconhecida |
| Sensíveis | Não identificados isoladamente; chave de e-mail é pseudonimizada e vinculável |
| Sistemas | Redis, memória da requisição |
| Base | Proposta para revisão: [legal-basis.md](./legal-basis.md), A002 |
| Compartilhamentos / transferência | Operador Redis e região desconhecidos; não inferir fornecedor do SDK |
| Retenção | TTL normal 1 minuto; persistência, dumps/backups e tratamento de falhas desconhecidos |
| Segurança | Auth fail closed; limite 10/min IP e e-mail; limite global 50/min pode falhar aberto; SHA256 do e-mail não é HMAC |
| Risco / RIPD | Avaliação contextual pendente; não demonstrado alto risco |
| Owner | Pendente |

### A003 — Publicação transacional de evento de acesso

| Campo | Valor |
|---|---|
| Finalidade | Registrar pedido para futura entrega assíncrona de magic link |
| Estado | Persistência implementada; processamento/entrega não implementados |
| Dados / origem | user_id, login_token_id; event_id/tipo/timestamps gerados |
| Titulares / menores | Mesmos solicitantes de A001; idade desconhecida |
| Sensíveis | Identificadores vinculáveis, não anonimizados; contexto de A001 |
| Sistemas | PostgreSQL outbox_events |
| Base | Vinculada à finalidade de A001, [legal-basis.md](./legal-basis.md), A003 |
| Compartilhamentos / transferência | Apenas integração API/DB encontrada; futuro provedor de entrega ainda não identificado |
| Retenção | Desconhecida; sem purge encontrado; FK/cascata de usuário não cobre outbox |
| Segurança | Payload AES-GCM versionado; commit transacional junto de usuário/token |
| Risco / RIPD | Avaliação contextual pendente |
| Owner | Pendente |

### A004 — Diagnóstico operacional

| Campo | Valor |
|---|---|
| Finalidade | Diagnosticar falhas e disponibilidade da API |
| Estado | slog de startup/erro implementado |
| Dados / origem | Mensagens, erros encadeados das dependências, ambiente/endereço; PII potencial em erros não excluída por teste |
| Titulares / menores | Titulares envolvidos na falha se erro carregar dados; operadores humanos; idade desconhecida |
| Sensíveis | Não há campo sensível explícito; risco de conteúdo em erros requer revisão |
| Sistemas | Destino padrão de log; coletor/armazenamento externos desconhecidos |
| Base | Proposta para revisão: [legal-basis.md](./legal-basis.md), A004 |
| Compartilhamentos / transferência | Destinos, pessoal autorizado, fornecedor e região desconhecidos |
| Retenção | Desconhecida |
| Segurança | Não imprime explicitamente e-mail/token no handler; sem redator central identificado |
| Risco / RIPD | Depende do conteúdo e destino; não demonstrado alto risco |
| Owner | Pendente |

## Estrutura disponível sem tratamento integrado

`sessions` armazena user_id, secret_hash, expires_at, revoked_at, used_at e created_at (`000003_create-auth.up.sql:23`). Domínio e repositório existem; main/HTTP não criam sessões. Tratar como estrutura preparada, não autenticação de produção comprovada. Futuro A005 (gestão de sessão): base art. 7º, V ligada ao acordo gratuito A001, retenção e integração pendentes; incluir revogação, exportação e exclusão.

## Atividades futuras descritas no Vault

Fontes: `../Obsidian Vault/Docs/Architecture/Privacidade e Consentimento.md:28`, `:48`; `../Obsidian Vault/Docs/Architecture/Sincronização Offline.md:18`; `../Obsidian Vault/Plans/Active/04 - Privacidade e Consentimento.md:98`. Não encontradas implementadas na API atual.

| ID / finalidade | Dados e sistemas previstos | Base proposta e pendências |
|---|---|---|
| A006 — Progresso/sincronização | Histórico/posição de leitura, user_id, device_id, sequência e instante; SQLite Android e futura API/DB | [legal-basis.md](./legal-basis.md); considerar dado sensível quando revelar convicção religiosa (art. 5, II) e base art. 11. Docs adotam consentimento específico. Necessidade de cada campo, operadores, transferência e retenção pendentes. Alto risco/RIPD requer avaliação contextual, não conclusão automática. |
| A007 — Registro de consentimentos | Finalidade, decisão, versão de aviso e timestamps vinculados ao usuário; ledger ainda ausente | [legal-basis.md](./legal-basis.md); prestação de contas e prova, arts. 8/9; retenção proporcional e preservação justificável pendentes. Não impor consentimento para toda operação. |
| A008 — Lembretes/push | Token push/dispositivo e preferências; fornecedor não definido | [legal-basis.md](./legal-basis.md); opcional na documentação. Dados necessários, fornecedor/região, prazo e revogação pendentes. |
| A009 — Analytics/personalização | Eventos e perfil ainda não especificados; ferramentas não identificadas | [legal-basis.md](./legal-basis.md); opcional na documentação. Avaliar se haverá dados pessoais/sensíveis, terceiros, finalidade e técnica de minimização antes de construir. |
| A010 — Atendimento aos titulares | Pedido, identidade estritamente necessária à verificação, andamento, resposta e justificativas | [legal-basis.md](./legal-basis.md); arts. 18/19; canal e workflow não encontrados na API. Retenção proporcional, controles de acesso e responsável pendentes. |

Para todas as atividades futuras: owner pendente; usuários/crianças/adolescentes/idosos não caracterizados; retenção desconhecida; destinatários e transferências não verificados. Backup do sistema operacional/SQLite precisa ser incluído no desenho. A leitura pública de catálogo sem identificação não é automaticamente tratamento sensível; sua associação identificável ao titular exige análise.

## Pendências para validar o inventário

- Identificar controlador, owners e operadores contratados, regiões e destinatários internos.
- Confirmar quais clientes web/mobile realmente coletam dados hoje, cookies/storage, telemetria, permissões e backups; este mapa não auditou esses códigos.
- Aprovar finalidade/base por atividade; definir retenção e exceções sem importar os prazos fictícios do template da skill.
- Caracterizar público e riscos pelos critérios normativos antes de afirmar incidência específica sobre menores ou necessidade de RIPD.

Fundamentos: LGPD arts. 5, II/IX/X/XI; 6, I/III/VII/X; 7/11; 8/9; 12; 15/16; 18/19; 37; 46. O mapa é evidência de descoberta e não substitui verificação operacional.
