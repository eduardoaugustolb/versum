# Gap analysis — Versum — 2026-10-01

**Nota de acompanhamento:** este é o retrato da auditoria antes das revisões documentais solicitadas no mesmo dia. As notas do Vault foram ajustadas; referências de linha abaixo são históricas. Ver [STATUS.md](STATUS.md) e [checklist feat/auth](feat-auth-implementation.md) para a sequência atual. Nenhum controle de código foi implementado nesta revisão.

**Cenário A: projeto novo**, confirmado pelo solicitante. A documentação declara pré-alpha, sem produção. As prioridades abaixo são critérios de prontidão antes da coleta real/lançamento, não constatação de infrações em produção. Auditoria do working tree, inclusive arquivos ainda não commitados, na base `5afbb622430bc7d3539c4d87c8a83516a24ece7e`.

Legenda: RED = controle não encontrado no escopo; YELLOW = parcial, planejado ou depende de comprovação externa; GREEN = evidência técnica delimitada. Prioridade não corresponde à dosimetria da ANPD. Responsáveis são papéis propostos, ainda não designados. Não foram fixados prazos arbitrários.

## Resultado

Base técnica de proteção já presente, governança ainda em desenho. Nenhuma conclusão de conformidade integral é sustentada pelo material. Há **12 gaps: 3 críticos, 7 altos e 2 médios**, além de uma falha de validação técnica registrada separadamente.

| ID | Estado / prioridade | Evidência e impacto | Ação / critério de aceite | Responsável proposto / fundamento |
|---|---|---|---|---|
| G01 | YELLOW / crítica | Plano 04:55–57 deixa bases pendentes, mas Architecture/Privacidade:47–55 já exige consentimento obrigatório de conta. | Aprovar finalidade/base de cada atividade e alinhar o fluxo ao resultado; não confundir ciência da política com consentimento. | Controlador + jurídico; LGPD arts. 7º–11. |
| G02 | RED / crítica | Política publicada e identificação do controlador não demonstradas; Vault Privacidade:19–20 confirma ausência da política. | Antes de convites/coleta real, publicar informação correspondente ao funcionamento efetivo e canal de contato; validar eventual regime ATPP. | Controlador; arts. 9º e 41; Res. 2/2022 arts. 3º e 11. |
| G03 | YELLOW / crítica | Histórico religioso previsto no PRD:49–52; menores e consentimento ainda sem desenho operacional no Plano 04:71–77. | Antes de progresso/push, avaliar inferência religiosa, melhor interesse, público provável, vínculo responsável e riscos; documentar decisão sobre RIPD. | Produto + jurídico; arts. 5º, II, 11, 14 e 38; ECA Digital arts. 1º–2º, 8º e 24; Res. 2/2022 art. 4º. |
| G04 | RED / alta | Não há ledger/rotas de consentimento; migrations auth:1–35 e routes:58–60. Plano 04:58–65 está pendente. | Implementar apenas finalidades que dependem dessa base, com versão do texto, evidência proporcional e revogação efetiva; ver consent/review.md. | API + apps; arts. 8º e 18, IX. |
| G05 | RED / alta | Rotas atuais não incluem `/me`, exportação, exclusão ou alternativa operacional comprovada. | Definir canal autenticado, resposta, correção e eliminação; cobrir cache, outbox e operadores; testar isolamento entre titulares. | API + atendimento; arts. 18–19. |
| G06 | YELLOW / alta | Expiração de tokens/sessões não elimina linhas. Outbox sem FK ao usuário (migration 000004:1–27); prazos vagos no Vault Privacidade:65–72. | Definir retenção por finalidade/estado e implementar purge; vincular outbox à eliminação; validar backup/restauração e prova de consentimento. | API + infraestrutura; arts. 6º, III, e 15–16. |
| G07 | YELLOW / alta | `routes.go:95–97` usa SHA-256 puro de e-mail, contrariando Plano 04:38–39. Redis também recebe IP na chave (`middleware/rate_limit.go:84`). | Usar HMAC separado para identificador previsível, inventariar IP e pseudônimos, restringir acesso e testar descarte. Hash não equivale a anonimização. | API; arts. 6º, III/VII/VIII, 12 e 46. |
| G08 | YELLOW / alta | Repositório pode atualizar lookup, sem recifrar e-mail; `docs/auth-known-limitations.md:13–15`. Outbox também referencia versão criptográfica. | Planejar recifragem idempotente de dados e backups, verificar todas as dependências antes de retirar chaves e demonstrar custódia produtiva. | API + infraestrutura; art. 46. |
| G09 | RED / alta | Runbook aparece somente como checklist futuro no Plano 04:73–74. | Designar equipe/fluxo, critérios de notificabilidade, relógio e registro; testar tabletop antes de tratar dados reais. | Controlador + segurança; LGPD art. 48; Res. 15/2024 arts. 5º, 6º, 9º e 10. |
| G10 | YELLOW / alta | Operadores/regiões/contratos não comprovados; Privacidade:81–85 diz mapa zerado até deploy. TLS do servidor depende de infraestrutura (`cmd/api/main.go:158`). | Inventariar fornecedores reais e backups, instruções contratuais, acessos/regiões, mecanismo válido de transferência e evidência de TLS/custódia. Não tratar bibliotecas locais como operadores por definição. | Infraestrutura + controlador; arts. 33–39 e 46; Res. 19/2024 Anexo I art. 9º. |
| G11 | YELLOW / média | `docs/auth-known-limitations.md:19` diz não haver rate limit, mas `routes.go:32–37` o implementa. Vault apresenta planos como comportamento estável. | Sincronizar docs com o código, distinguir implementado/planejado e vincular evidências de execução. | Engenharia; arts. 6º, VI/X, e 37. |
| G12 | YELLOW / média | Decoder sem limite de corpo (`routes.go:56`); erros de dependências enviados a slog (`routes.go:74,88`; middleware:92), sem garantia demonstrada de sanitização. | Limitar body e validar erros com entradas sintéticas de segredo/PII. É risco potencial: não foi demonstrado vazamento real. | API; arts. 6º, VII/VIII, e 46. |

## Controles positivos com alcance delimitado

- GREEN: e-mail protegido por AES-GCM e índice HMAC versionado; hashes de tokens e sessões; payload de outbox cifrado. Evidências detalhadas em [code-evidence.md](code-evidence.md).
- GREEN: rate limits de autenticação por IP/e-mail existem e falham fechado quando cache indisponível; resposta bem-sucedida neutra.
- YELLOW: ROPA final, encarregado/designação, treinamento, controles de infraestrutura e contratos não são verificáveis apenas pelo código. [data-map.md](data-map.md) é inventário diagnóstico, não ROPA aprovado.

## Correções jurídicas das docs

1. **ECA Digital:** art. 9º, §1º trata acesso a conteúdo impróprio/proibido para menores; não é regra universal de verificação a cada acesso de qualquer leitor bíblico. Acesso provável e vínculo responsável exigem avaliação própria (arts. 1º–2º e 24). Não acrescentar coleta de documentos/biometria por automatismo.
2. **RIPD:** manter o gate interno conservador para progresso, mas distinguir escolha do produto de obrigação legal automática. A ANPD recomenda avaliação de alto risco e pode exigir relatório; a Res. 2/2022 art. 4º combina critério geral e específico. Não foi adotado o limiar fixo de dois milhões de titulares da referência da skill, pois esse dispositivo não o estabelece.
3. **Transferência:** SCCs brasileiras são um dos mecanismos válidos, não a única alternativa para qualquer dado fora do Brasil (Res. 19/2024, Anexo I art. 9º).
4. **DSAR:** os 15 dias do art. 19, II referem-se à declaração completa de confirmação/acesso. Uma meta interna igual para outros direitos deve ser identificada como meta, não prazo universal desse artigo; ATPP depende de enquadramento comprovado.
5. **Prova:** reconciliar “nunca apaga prova” com cascata de consentimentos. Append-only não dispensa base/prazo de conservação (arts. 15–16).

## Sequência proposta

1. Resolver G01, G02 e G03 com os responsáveis por produto e tratamento, antes de convites reais.
2. Alinhar autenticação, consentimento quando aplicável, DSAR, retenção e cache (G04–G08), com testes de efeitos concretos.
3. Fechar incidentes e infraestrutura/fornecedores (G09–G10); atualizar docs e bordas HTTP (G11–G12).

## Validação técnica do estado atual

`go test ./...` em 2026-10-01 terminou com falha: `internal/identityaccess/application/commands/consume_magic_link_test.go:29: undefined: NewConsumeMagicLink`. Esse teste era arquivo local não rastreado antes da auditoria. Os demais resultados exibiram principalmente cache; não comprovam execução contra PostgreSQL real. Falha funcional de desenvolvimento, não infração LGPD por si só. Nenhum código/teste existente foi alterado.

## Fontes oficiais consultadas

- [LGPD — texto compilado](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709.htm).
- [ECA Digital — texto compilado](https://www.planalto.gov.br/ccivil_03/_ato2023-2026/2025/lei/l15211.htm).
- [Resolução ANPD 2/2022](https://www.gov.br/anpd/pt-br/acesso-a-informacao/institucional/atos-normativos/regulamentacoes_anpd/resolucao-cd-anpd-no-2-de-27-de-janeiro-de-2022).
- [ANPD — orientação sobre RIPD](https://www.gov.br/anpd/pt-br/canais_atendimento/agente-de-tratamento/relatorio-de-impacto-a-protecao-de-dados-pessoais-ripd).
- [ANPD — regulamento de incidentes](https://www.gov.br/anpd/pt-br/assuntos/noticias/anpd-aprova-o-regulamento-de-comunicacao-de-incidente-de-seguranca).
- [Resolução ANPD 19/2024](https://www.gov.br/anpd/pt-br/acesso-a-informacao/institucional/atos-normativos/regulamentacoes_anpd/resolucao-cd-anpd-no-19-de-23-de-agosto-de-2024).
