# Bases legais por atividade — acesso definido e demais avaliações

Última atualização: 2026-10-01. Cenário confirmado pelo solicitante: **A — projeto novo**. Documento técnico de auditoria; as propostas abaixo não substituem decisão do controlador nem autorizam coleta. Em 2026-10-01, o solicitante esclareceu que o Versum será um app aberto ao público, sem finalidade comercial ou captação para outro produto, e pediu definição da base e do acordo de acesso. Decisão de desenho registrada em [access-contract.md](access-contract.md); controlador/contato e revisão jurídica externa permanecem pendentes.

Fontes: [LGPD, texto oficial](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709.htm); evidências documentais em [docs-evidence.md](docs-evidence.md), especialmente DOC-01–12. Revisar cada proposta ao consolidar o mapa de dados e as operações reais.

## 1. Conta, magic link e sessão

- **Finalidade:** atender pedido de acesso e manter sessão revogável do usuário.
- **Dados:** e-mail cifrado, índice cego HMAC, identificador, hashes de token/sessão, timestamps; demais metadados de dispositivo a confirmar no código.
- **Sensibilidade:** campos ordinários identificadores; avaliar associação contextual ao serviço religioso. Hash/HMAC/cifra não retiram automaticamente dado do escopo (arts. 5º, I–II; 12).
- **Base definida para dados comuns:** art. 7º, V, para pedido de acesso e execução do acordo de uso gratuito da conta opcional. Leitura pública dispensa conta. A base cobre apenas e-mail, identificação interna, token/sessão e datas necessários ao acesso solicitado; não histórico religioso, marketing ou analytics.
- **Contrato:** termos de uso gratuito, versão recuperável e pedido informado; confirmação de controle do e-mail antes de ativar conta. Não existe consentimento obrigatório de `conta_acesso`. Se a operação concreta revelar convicção religiosa, não ativá-la sob art. 7º, V: requer hipótese válida do art. 11. Ver [minuta dos termos](policies/access-terms-v1-draft.md).
- **Retenção:** validade de token não equivale a prazo de conservação; purge 24h é proposta das docs, sessões/conta precisam de critério e duração.
- **Revogação:** se base contratual, encerramento da conta e exercício de direitos; se consentimento validamente escolhido, revogação facilitada (art. 8º, §5º).

## 2. Entrega do e-mail de autenticação pela outbox

- **Finalidade:** entregar o link solicitado e recuperar falhas de envio.
- **Dados atuais:** payload cifrado com `user_id` e `login_token_id`, identificador/tipo/timestamps do evento. Não contém token bruto (`request_magic_link.go:106–109`). Destino de e-mail, material para entrega e metadados de processamento são desenho futuro: worker não implementado. Tentativas/erro redigido existem como colunas, sem processamento demonstrado.
- **Base definida:** art. 7º, V, mesma finalidade de pedido/acesso ao app gratuito; dados comuns estritamente necessários à entrega. Sem finalidade de marketing.
- **Pendências:** provedor e papel efetivo, regiões/acesso internacional, contrato e transparência (arts. 9º, 33–39). Integração planejada não prova fornecedor contratado.
- **Retenção:** separar segredo de entrega dos metadados necessários. Faixa documental 30–90 dias não constitui prazo normativo ou justificativa; eliminar payload quando não mais necessário.
- **Revogação/encerramento:** cancelamento de entregas pendentes conforme pedido, expiração e exclusão.

## 3. Segurança, rate limit e prevenção de abuso

- **Finalidade:** proteger autenticação e disponibilidade, sem rastreamento comportamental.
- **Dados:** IP/chave derivada, contador, janela, timestamps, eventual e-mail HMAC e metadados estritamente necessários.
- **Base candidata para dados comuns:** art. 7º, IX combinado com art. 10, condicionada a avaliação documentada de legítimo interesse, necessidade, expectativas e salvaguardas. Também avaliar necessidades contratuais concretas, sem multiplicar bases para a mesma finalidade.
- **LIA:** pendente; nenhum LIA foi produzido ou aprovado. A proposta não permite usar legítimo interesse para dados sensíveis.
- **Sensíveis:** art. 11, II, g pode ser avaliado apenas no alcance de prevenção à fraude e segurança em identificação/autenticação; não autoriza analytics religioso nem qualquer monitoramento de segurança.
- **Retenção:** docs propõem janela de rate limit de 1 minuto; confirmar TTL efetivo, réplicas e logs. Dados técnicos ligados a pessoa/IP continuam no mapa.
- **Direitos:** transparência, oposição quando cabível e canal de atendimento (arts. 9º; 18, §2º).

## 4. Prova de consentimento e versões de política — futuro

- **Finalidade:** demonstrar manifestações para finalidades baseadas em consentimento, inclusive revogação.
- **Dados:** titular, finalidade, status, versão recuperável do texto, datas, evidência mínima; IP/UA/snapshot exigem justificação específica.
- **Base candidata para registro necessário:** art. 7º, II, vinculada à obrigação concreta de provar consentimento (art. 8º, §2º), com análise de adequação e conservação. Se a própria prova revelar dados sensíveis, avaliar art. 11, II, a no alcance estrito da obrigação; não presumir aplicabilidade a todos os campos.
- **Pendências:** reduzir evidência, proteger UA/IP, definir retenção e compatibilizar append-only com eliminação. Defesa em processo pode demandar art. 7º, VI/art. 11, II, d em situação concreta; não justifica conservação eterna.
- **Revogação:** registrar nova manifestação e cessar a atividade consentida; avaliar conservação limitada da prova nos arts. 15–16.

## 5. Direitos do titular e exclusão — futuro

- **Finalidade:** receber, autenticar proporcionalmente e executar pedidos de direitos.
- **Dados:** identidade, conteúdo/pedido, resposta, datas e prova mínima de atendimento; exportação contém dados da conta.
- **Base candidata:** art. 7º, II em cumprimento dos arts. 18–19; para sensíveis, art. 11, II, a quando necessário a essa obrigação.
- **Pendências:** canal, verificação proporcional, correção, confirmação/acesso, compartilhamentos, revogação, portabilidade conforme regulamentação, trilha mínima e prazo por tipo de pedido. Art. 19, II é prazo de declaração completa de confirmação/acesso, não regra universal para qualquer DSAR.
- **Retenção:** prazo e necessidade da prova de atendimento pendentes; excluir dependências locais/remotas e impedir recriação por sync.

## 6. Progresso e sincronização offline — futuro

- **Finalidade:** salvar e reconciliar leitura privada entre dispositivos, a pedido da pessoa.
- **Dados:** usuário/dispositivo/evento, sequência, referência bíblica, posição/avanço/conclusão e timestamps; cópia SQLite, outbox local e projeção remota.
- **Sensibilidade:** proteção reforçada por risco de revelar/inferir convicção religiosa (art. 5º, II); contexto e inferência precisam de descrição.
- **Base definida no desenho futuro sensível:** consentimento específico e destacado para finalidade determinada, art. 11, I; não ativar progresso associado à conta antes da escolha demonstrável. Consentimento não pode ser genérico nem embutido na conta. Nenhuma hipótese alternativa de art. 11 foi demonstrada no escopo; controles e revisão do contexto continuam pendentes.
- **Pendências:** necessidade de histórico vs cursor, consentimento offline, revogação entre dispositivos, prevenção de ressurreição de eventos, exportação/exclusão e avaliação de impacto. RIPD é condição proposta do projeto; sensibilidade isolada não significa obrigação legal automática.
- **Retenção/revogação:** duração não definida; revogação deve cessar coleta/sync consentida e permitir eliminação, ressalvadas hipóteses legais concretas.

## 7. Push e lembretes — futuro

- **Finalidade:** enviar lembrete escolhido, na frequência e horário definidos pelo usuário.
- **Dados:** token push, usuário/dispositivo, dias/horário, fuso e registros mínimos de entrega; conteúdo da mensagem a definir.
- **Base candidata:** art. 7º, I para preferência comum opcional; art. 11, I quando a associação/conteúdo revelar religião. Permissão do sistema operacional não substitui automaticamente consentimento LGPD válido.
- **Pendências:** finalidade granular, não obrigatória para leitura/conta; minimizar conteúdo na tela bloqueada, mapear FCM/provedor e fluxo internacional, definir retenção e separar lembrete de analytics de abertura.
- **Revogação:** gratuita/facilitada; interromper envio, remover token/vínculo quando desnecessário e tratar fila pendente.

## 8. Métricas e analytics — futuro

- **Finalidade:** medir retomada 7/30d, capítulos, sync, downloads e entrega/abertura de lembretes conforme PRD.
- **Dados:** eventos e correlações ainda não definidos. Não considerar anônima a coleta individual que será agregada posteriormente.
- **Base candidata:** definir por métrica. Consentimento específico/destacado do art. 11, I para analytics que revele hábitos religiosos; art. 7º, I para analytics opcional comum. Art. 7º, IX/art. 10 só pode ser considerado em atividade comum com LIA concreto; não serve para sensíveis.
- **Anonimização:** avaliar art. 12 e testes de reidentificação; dado comprovadamente anonimizado pode ficar fora do regime, mas transformação prévia de dado pessoal exige hipótese válida. Agregado com grupos pequenos ou eventos raros pode continuar identificável.
- **Pendências:** desenho dos eventos, finalidades distintas, fornecedor, intervalos, limiares, privacidade de menores e retenção. Nenhum LIA foi aprovado para analytics.
- **Revogação:** cessar tratamentos consentidos; distinguir dados pessoais identificáveis de resultados efetivamente anonimizados.

## Condições comuns de aprovação

Validar controlador, público/menores, cada finalidade e campo, base efetiva, necessidade, transparência, retenção, compartilhamento e mecanismos de direitos antes de convites externos. Para crianças/adolescentes, melhor interesse deve prevalecer e o desenho deve observar art. 14 e Lei 15.211/2025; não concluir que apenas consentimento parental serve para toda operação. Não preencher lacunas com bases genéricas ou selecionar simultaneamente bases intercambiáveis para escapar da revogação.

Resumo: acesso, entrega e sessão de dados comuns têm base definida no acordo gratuito (art. 7º, V). As demais atividades mantêm suas decisões específicas e pendências; LIA não aprovado. A definição de desenho não declara implantação, publicação ou revisão jurídica externa.

## Correspondência com o inventário técnico

| ID em data-map.md | Seção de base proposta |
|---|---|
| A001 conta e A005 sessão futura | 1 |
| A002 antiabuso e A004 logs operacionais | 3; logs exigem finalidade/necessidade e retenção próprias na validação |
| A003 outbox | 2 |
| A006 progresso | 6 |
| A007 consentimentos | 4 |
| A008 lembretes | 7 |
| A009 métricas | 8 |
| A010 DSAR | 5 |

As numerações de seção não são IDs de atividade; essa correspondência evita atribuir a base da outbox ao rate limit.
