# Política de privacidade do Versum — minuta v1

**Status: NÃO PUBLICADA; revisão jurídica e factual pendentes.** Elaborada em 2026-10-01. Vigência: PENDENTE de aprovação e publicação. Esta minuta contém campos essenciais em aberto e não é um aviso pronto para uso com titulares.

O Versum é definido como aplicativo gratuito aberto ao público, sem finalidade de vender ou promover outro produto, com leitura sem conta e conta opcional. Está em desenvolvimento. A API examinada recebe e-mail para solicitar acesso, protege o e-mail e registra um pedido de magic link. Também usa controles de frequência e logs de falhas. Entrega do link, autenticação por sessão, progresso entre dispositivos, consentimentos e lembretes ainda não estão integrados na API examinada. Os dados e recursos efetivamente disponíveis nos clientes web/Android precisam ser confirmados antes da publicação.

## Quem responde pelos dados

Controlador: **PENDENTE — nome/razão social, identificação pertinente e endereço**. Contato do controlador: **PENDENTE — canal público para dúvidas e exercício de direitos**. Encarregado ou canal aplicável: **PENDENTE — designação/enquadramento e contato**. A pessoa responsável pelo aplicativo deverá preencher e validar estas informações.

O controlador define finalidades, meios e decisões sobre o tratamento, organiza o atendimento dos direitos e verifica as medidas de proteção. Operadores identificados e contratados deverão tratar dados conforme instruções e responsabilidades aplicáveis. A relação efetiva de agentes e contatos está **PENDENTE**; bibliotecas de software, por si só, não identificam quem opera a infraestrutura.

## O que a API atual trata e para quê

| Finalidade | Dados e modo de tratamento | Duração / base proposta |
|---|---|---|
| Receber pedido de acesso e identificar a conta | E-mail informado pelo solicitante, identificador gerado, datas e hash do token; banco PostgreSQL | Conta: **PENDENTE**. Token válido por 15 minutos; remoção física ainda não implementada. Base definida de desenho: art. 7º, V para os dados comuns necessários ao pedido e ao acordo de uso gratuito. Conta opcional, confirmação do endereço antes de ativação; contrato/controles ainda a implementar. |
| Registrar o pedido para futura entrega | Identificadores da conta/token e metadados de evento em outbox cifrada; não há envio integrado na API examinada | **PENDENTE** prazo e mecanismo de limpeza. Art. 7º, V no alcance necessário à entrega do acesso solicitado; sem marketing. |
| Proteger disponibilidade e evitar solicitações excessivas | IP observado, chave SHA256 derivada do e-mail, contadores no Redis | Janela normal de 1 minuto; persistência e backups: **PENDENTE**. Avaliação de legítimo interesse para dados comuns (arts. 7º, IX e 10) **pendente**, incluindo necessidade e salvaguardas. |
| Diagnosticar falhas | Mensagens e erros das dependências em logs; conteúdo pessoal eventual precisa ser revisado | Destino, acesso e duração: **PENDENTE**. Finalidade, necessidade e base própria devem ser validadas. |

E-mail e identificadores protegidos ainda podem ser associados a uma pessoa. O código usa cifra AES-GCM de e-mail/outbox, HMAC para pesquisa de conta, hashes de tokens e limitação de frequência. Isso não significa anonimização. Segurança de hospedagem, conexões, acesso às chaves, logs e backups está **PENDENTE de verificação**.

A participação em um serviço religioso pode permitir inferências sobre convicção religiosa, conforme o contexto. Se o tratamento revelar dado sensível, uma base do art. 7º para dados comuns não será suficiente: a decisão deverá observar o art. 11 e as salvaguardas correspondentes. Essa avaliação está **PENDENTE**.

## Recursos previstos e escolhas

Progresso/sincronização entre dispositivos, registro de consentimentos, lembretes/push, analytics e personalização estão previstos em documentos de projeto. Não são declarados como coleta ativa pela API auditada. Antes de ativá-los, o Versum deverá informar dados necessários, finalidade, base efetiva, duração, destinatários e escolhas em versão revisada desta política.

O desenho prevê escolhas separadas para tratamentos opcionais. Se consentimento for a base escolhida, deverá ser livre, informado, específico e demonstrável, com revogação facilitada. O aceite de um aviso de privacidade não substitui automaticamente consentimento válido. Para progresso que revele convicção religiosa, a proposta é avaliar consentimento específico e destacado do art. 11, I; decisão **PENDENTE**. Recusa a opções facultativas não deverá impedir recursos que independam delas.

## Compartilhamento e transferências

A API se comunica com PostgreSQL e Redis. **PENDENTE**: identificar hospedagem, operadores contratados, destinatários internos/externos, finalidade de cada acesso, regiões e eventual acesso a partir de outro país. Nenhum provedor de e-mail, push ou analytics foi identificado como integração ativa na API. Não há confirmação de que os dados permaneçam exclusivamente no Brasil.

Se houver transferência internacional de dados pessoais, deverá ser identificada a hipótese válida do art. 33 e o mecanismo aplicável, com transparência e garantias. Cláusulas-padrão brasileiras são um mecanismo possível, não uma conclusão automática para todo fornecedor.

## Conservação e exclusão

Os prazos e critérios de conta, eventos, sessões futuras, logs e backups estão **PENDENTES**. O parâmetro de retenção de token de 24 horas é uma proposta no código, sem rotina de eliminação; não é prazo legal nem garantia de limpeza. Expiração de um token impede seu uso pelo desenho do domínio, mas não elimina o registro.

Antes de publicar, definir duração necessária por finalidade, procedimento de eliminação, efeitos nos dispositivos/filas/backups e exceções legais concretas dos arts. 15–16. Uma solicitação de exclusão não implica conservar tudo para sempre nem apagar registros sujeitos a conservação legal sem análise.

## Seus direitos e como pedir

Canal: **PENDENTE — inserir contato público funcional**. O atendimento deverá ser gratuito e verificar a identidade de forma proporcional, sem coletar documentos desnecessários. Não há endpoint de atendimento aos titulares na API atual; um canal humano funcional poderá atender os direitos enquanto o suporte técnico é construído.

Nos termos do art. 18, você pode solicitar confirmação do tratamento; acesso; correção de dados incompletos, inexatos ou desatualizados; anonimização, bloqueio ou eliminação de dados desnecessários, excessivos ou tratados irregularmente; portabilidade nos termos da regulamentação; eliminação de dados tratados com consentimento, ressalvado o art. 16; informação sobre compartilhamentos e sobre recusar consentimento e suas consequências; e revogação do consentimento. Também pode exercer oposição em caso de descumprimento da LGPD e peticionar à ANPD ou aos órgãos competentes.

Para confirmação e acesso, o art. 19 prevê resposta simplificada imediatamente ou declaração clara e completa em até 15 dias. Esse prazo não deve ser apresentado como prazo universal para todos os pedidos. Se não for possível tomar providências imediatamente, deverá haver resposta com razões de fato ou de direito, ou indicação do agente competente quando aplicável (art. 18, §4º). Pedidos relacionados a decisões unicamente automatizadas que afetem seus interesses serão tratados conforme o art. 20, caso tais decisões venham a existir. Não foi identificado perfil automatizado desse tipo; limitação técnica de frequência existe e deve ser explicada quando pertinente.

## Cookies, clientes e pessoas menores de idade

Não foi encontrado cookie de sessão ativo na API examinada. Os clientes web/Android, cookies, armazenamento local, permissões e backups ainda precisam ser inventariados; esta afirmação não garante ausência de cookies no produto inteiro.

O público e a presença de crianças/adolescentes estão **PENDENTES de caracterização**. O projeto deverá observar o melhor interesse e as condições do art. 14; incidência e medidas da Lei 15.211/2025 precisam de avaliação do produto. Não anunciar bloqueio etário, verificação ou autorização parental que não estejam efetivamente implantados.

## Alterações desta política

Local de publicação e comunicação de alterações: **PENDENTE**. Mudanças de finalidade, base, compartilhamento ou duração deverão atualizar a política e as escolhas cabíveis antes de novo tratamento. Manter versões acessíveis conforme procedimento a definir.

Histórico: v1 minuta, 2026-10-01 — primeira consolidação do diagnóstico. Nenhuma aprovação, publicação ou vigência registrada.

Referências para revisão: [LGPD oficial](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709.htm), arts. 5º, 6º, 7º–11, 14–20, 33, 37, 39, 41 e 46; [mapa](../data-map.md), [bases propostas](../legal-basis.md), [evidências técnicas](../code-evidence.md).
