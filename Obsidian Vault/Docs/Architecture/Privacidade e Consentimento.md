---
title: "Privacidade e Consentimento"
section: Docs
subsection: Architecture
type: architecture
status: draft
tags: [versum, docs, architecture, privacy, lgpd, consentimento]
up: "[[Docs/Architecture/_Index|Arquitetura]]"
prev: "[[Docs/Architecture/Sincronização Offline]]"
next: "[[Rules/_Index]]"
related: ["[[PRD]]", "[[Rules/02 - Segurança]]", "[[Docs/Architecture/Autenticação e Sessões]]", "[[Plans/Active/04 - Privacidade e Consentimento]]"]
---

# Privacidade e Consentimento

🏠 [[_Index|Home]] › 📚 [[Docs/_Index|Documentação]] › 📐 [[Docs/Architecture/_Index|Arquitetura]] › **Privacidade e consentimento**

O Versum está em pré-alpha, sem produção e sem política publicada. Este documento
separa o comportamento alvo dos controles observados na API; não autoriza coleta
externa nem comprova condições de infraestrutura. A execução está em
[[Plans/Active/04 - Privacidade e Consentimento|Plano 04]].

## Estado atual e dados pessoais

- A API expõe solicitação de magic link; consumo, sessão HTTP autenticada e
  entrega por worker ainda precisam ser concluídos.
- Conta: identificador, e-mail cifrado AES-GCM, índice cego HMAC versionado e
  timestamps. O token é aleatório, persistido como hash e expira em 15 minutos.
- Outbox: payload cifrado com identificadores de usuário e token. O token bruto
  necessário à entrega ainda não está no payload atual; o desenho alvo deve
  preservá-lo somente cifrado até a entrega, sem logs ou colunas em texto puro.
- Operação: rate limit por IP/e-mail, com janela de um minuto. A chave de e-mail
  usa hoje SHA-256 determinístico; deve passar a HMAC secreto versionado. IP,
  hashes vinculáveis, timestamps e payload cifrado continuam dados pessoais.
- Estruturas de sessão existem, mas não demonstram sessão/cookie implementados.
  Consentimento, DSAR e rotinas de purge não estão integrados à API.
- Progresso, dispositivo, push e métricas são futuros. O mapa deve cobrir também
  SQLite/outbox local, Redis, filas, réplicas, backups e fornecedores utilizados.

O catálogo bíblico é público e não exige conta. Necessidade, finalidade e
retenção devem ser registradas para cada atividade (LGPD arts. 6º e 37).

## Finalidades e bases legais

Base definida para os dados comuns necessários ao acesso: **art. 7º, V**,
execução de acordo de uso gratuito ou pedido preliminar do titular. O app é
aberto ao público, sem venda ou captação para outro produto. Leitura dispensa
conta; conta opcional é confirmada pelo acesso ao e-mail. Termos versionados
definem uso gratuito, sem criar consentimento geral ou autorizar marketing.
Responsável/contato, publicação e revisão jurídica ainda pendentes. Ser aberto
ao público não caracteriza, por si, tratamento pela Administração Pública.

O uso de um leitor bíblico não comprova a religião de toda pessoa. Identidade,
progresso e hábitos associados ao contexto podem revelar ou permitir inferência
de convicção religiosa: devem ser avaliados e protegidos como potencialmente
sensíveis (arts. 5º, II e 11). Contrato e legítimo interesse do art. 7º não
substituem uma hipótese do art. 11 quando houver tratamento sensível.

## Consentimento quando esta for a base

- Finalidade determinada, linguagem acessível e manifestação livre, informada e
  inequívoca; para sensíveis, específica e destacada (arts. 8º e 11, I).
- Progresso, lembretes e analytics têm decisões próprias. Recusa de uma
  finalidade opcional não bloqueia leitura ou conta; nenhuma caixa pré-marcada.
- O ledger futuro registra concessão/revogação e versão recuperável do texto.
  Hash isolado não substitui o texto nem prova que a informação foi apresentada.
- Prova deve ser mínima: IP, user-agent e snapshot só entram com necessidade,
  proteção e retenção justificadas. Append-only não significa guardar para sempre.
- Revogação gratuita/facilitada cessa tratamentos consentidos e filas pendentes,
  com efeito futuro; pedidos de eliminação seguem arts. 15–18. A eventual
  conservação limitada de prova exige hipótese concreta, não retenção eterna.

## Direitos do titular

O desenho deve oferecer canal acessível para confirmação, acesso, correção,
informação de compartilhamento, portabilidade conforme regulamentação,
eliminação/anonimização quando cabíveis, oposição e revogação (art. 18).
`GET /me`, `GET /me/export` e `DELETE /me` são suporte técnico proposto;
endpoints específicos não substituem o canal nem cobrem sozinhos todos os direitos.

O art. 19 prevê confirmação/acesso simplificados imediatamente ou declaração
completa em até 15 dias; esse prazo não deve ser apresentado como prazo legal
universal de todos os pedidos. O atendimento precisa de autenticação proporcional,
responsável, registro mínimo e critérios para resposta e exceções.

## Retenção e eliminação

Expiração lógica não é eliminação. Os prazos abaixo são **propostas do projeto**,
pendentes de justificativa e aprovação, não prazos prescritos pela LGPD:

- Token consumido/expirado: avaliar purge em 24 horas e precisar o gatilho.
- Sessão expirada/revogada: duração e rotina ainda a definir.
- Outbox: separar remoção do segredo após sua necessidade de metadados de
  entrega; substituir a faixa anterior de 30–90 dias por critério definido.
- Conta e prova de consentimento: definir ciclo de vida, hipóteses de conservação,
  segregação e descarte; compatibilizar o ledger com exclusão.

Eliminação deve alcançar Redis, filas/outbox e dados derivados. Backups precisam
de prazo de expiração e restauração que reaplique exclusões; não basta dizer que
“herdam o prazo”. Progresso offline futuro deve impedir que evento antigo recrie
dados após exclusão/revogação (arts. 6º, III, 15–16, 18 e 46).

Rotação de chaves deve cobrir ciphertext de e-mail e outbox antes de aposentar
chaves antigas; hoje há migração do índice de lookup, sem re-cifragem completa.
Custódia, acesso, réplicas, backups e TLS precisam de verificação operacional.

## Restrições de produto e métricas

[[PRD]] mantém progresso privado, sem feed, ranking, venda de dados ou
perfilamento publicitário. Compartilhamento é explícito e granular. Tokens,
e-mails, URLs de autenticação e progresso não entram em logs/erros públicos.

Métricas 7/30 dias e abertura de lembretes podem correlacionar pessoas e hábitos
religiosos antes da agregação. Definir base, coleta mínima e retenção por métrica;
agregação sozinha não garante anonimização (arts. 5º, III e XI, 11 e 12).

## Crianças e adolescentes

Acesso provável por menores deve ser avaliado antes do lançamento, à luz do
melhor interesse (art. 14) e da Lei 15.211/2025. Definir público, riscos, controles
aplicáveis e minimização dos dados de idade/responsável. O art. 9º, §1º trata de
verificação para conteúdo impróprio: não fundamenta sozinho a obrigação universal
de verificar idade a cada acesso a todo leitor bíblico. Avaliar também o vínculo
parental do art. 24 e a regulamentação aplicável. Não lançar sem essa decisão.

## Operadores, transferências e governança

Registrar fornecedores **planejados** e **efetivamente usados**, seus papéis,
dados, regiões, suboperadores e instruções. Postgres/Redis locais não dispensam
mapear tratamentos de desenvolvimento. Não há fornecedor de e-mail integrado
comprovado na API atual; e-mail, hosting, S3, push e analytics futuros exigem
avaliação conforme arts. 5º, 9º, 37, 39 e 46.

Fluxos internacionais/acesso remoto exigem mecanismo válido dos arts. 33–36.
Cláusulas-padrão da Res. 19/2024 são uma alternativa, não a única hipótese legal.

Antes de convites externos: aprovar bases/finalidades, publicar informações do
art. 9º, definir controlador e canal dos titulares, avaliar encarregado/regime
aplicável, aprovar retenção e preparar resposta a incidentes. RIPD deve ser
avaliado conforme risco e exigência da autoridade (art. 38); dados sensíveis
isolados não tornam o relatório automaticamente obrigatório em todo caso.

---

◀ [[Docs/Architecture/Sincronização Offline|Sincronização offline]] · próxima: [[Rules/_Index|Regras]] ▶
