# Decisão de base legal e contrato de acesso — Versum

Versão 1, 2026-10-01. Definição de desenho solicitada pelo responsável pelo projeto: app aberto ao público, sem finalidade de comercializar produto ou captar interessados para outro produto. Adota-se acesso gratuito, sem cobrança/assinatura ou uso da conta para promoção comercial. Identidade do controlador/contato e publicação dos termos ainda pendentes. Esta decisão não declara revisão jurídica externa, contrato já celebrado com usuários ou funcionalidades implantadas.

## Base definida para autenticação

**LGPD art. 7º, V — execução de contrato ou procedimentos preliminares a pedido do titular.** O contrato é um acordo de uso gratuito das funções de conta. Gratuidade e disponibilização ao público não tornam o acordo uma venda nem afastam por si a possibilidade de contrato. Trata-se de interpretação aplicada ao desenho, apoiada nos arts. 107, 422 e 425 do Código Civil; o fundamento LGPD é restrito à necessidade real de cada dado.

Abrange A001 (conta/pedido), A003 (entrega) e A005 (sessão): e-mail para enviar/verificar o acesso, identificador interno, token transitório, hashes e instantes necessários à validade/revogação. Não exige nome, CPF, endereço, filiação religiosa, perfil público ou fingerprint. Não criar conta persistente apenas porque qualquer pessoa informou o e-mail de um terceiro: o pedido inicia procedimento preliminar, e a conta é confirmada após a prova de acesso ao endereço.

O simples fato de um app ser público não o torna Administração Pública nem cria base de execução de política pública. Também não significa que os dados da conta sejam públicos. Se o responsável efetivo for órgão público, reavaliar o regime antes do lançamento.

## Limite de dados sensíveis

A conta identifica acesso; não registra religião ou usa o vínculo para inferir crença. Isso não elimina o risco contextual. Se a operação concreta revelar/tratar convicção religiosa, art. 7º, V não basta: não ativar essa operação sob a base contratual de dados comuns. O desenho de progresso associado à conta adotará consentimento específico e destacado, art. 11, I, antes da coleta, separado da conta. Sem consentimento de progresso, leitura e conta continuam disponíveis; progresso sincronizado permanece desativado. Lembretes/analytics têm finalidade e escolha próprias. Não se presume filiação religiosa de quem lê a Bíblia.

Segurança antiabuso (A002) e diagnóstico (A004) mantêm avaliação própria de legítimo interesse; a presente decisão não aprova um LIA inexistente. Direitos do titular usam obrigação legal no alcance necessário. Não recorrer a bases intercambiáveis para evitar revogação.

## Contrato de acesso e regras implementáveis

1. **Leitura pública:** catálogo e leitura disponíveis sem cadastro; acesso pode sofrer limite técnico antiabuso proporcional. Sem checkbox contratual ou login obrigatório para consultar conteúdo público.
2. **Conta opcional:** necessária somente às funções pessoais habilitadas. Formulário pede e-mail e apresenta links para termos e aviso. Nenhum aceite de analytics, progresso ou marketing é embutido.
3. **Pedido informado:** ação explícita “Enviar link de acesso” sob texto claro inicia o pedido. Servidor recebe versão válida dos termos apresentados; não aceita referência inexistente/desatualizada segundo regra de versão publicada. Guarda prova mínima do pedido, não IP/UA indiscriminadamente.
4. **Confirmação:** consumo atômico do link confirma controle do endereço e ativa a conta/sessão; a unidade de trabalho deve suportar identidade provisória ou pedido sem usuário definitivo. Eventos não verificados expiram e são limpos; sem conta ativa/telemetria pessoal permanente antes da prova.
5. **Aceite contratual:** conservar evento mínimo de pedido/aceite com referência ao texto recuperável, sujeito à retenção definida. Não armazenar esse aceite como `GRANTED(conta_acesso)` no ledger LGPD; é acordo de uso, não consentimento geral de dados.
6. **Encerramento:** logout revoga sessão; excluir conta encerra funções pessoais, cancela entregas e inicia eliminação conforme regras justificadas. A pessoa continua podendo ler conteúdo público.
7. **Termos alterados:** versionar, informar mudanças relevantes e obter nova manifestação quando necessário à alteração do acordo. Não converter mudança de termos em autorização automática de novas finalidades.
8. **Menores:** celebrar acordo válido requer capacidade/representação cabível e controles do ECA Digital; não tratar um clique como solução para todas as idades. Fluxo aplicável precisa estar definido antes de aceitar essas contas.

Texto proposto no formulário: “Informe seu e-mail para receber um link e acessar sua conta gratuita. Ao solicitar o acesso, você concorda com os Termos de uso. Veja no Aviso de privacidade como usamos os dados necessários à conta. O catálogo pode ser lido sem conta.” Não adicionar checkbox “consinto com todos os meus dados”.

## Mudanças na feat/auth

- Contrato HTTP de pedido inclui `terms_version_id`; título/texto/versão recuperáveis no servidor. A versão do aviso pode ser referenciada separadamente, sem representar consentimento.
- Refatorar criação automática atual: conta provisória com ativação após verificação ou pedido sem conta definitiva; limitar repetição e limpar pedidos/identidades não confirmados.
- Consumo cria/ativa conta e sessão atomicamente e finaliza o aceite atribuído ao titular verificado. Não confundir prova de controle do e-mail com identidade civil ou idade comprovada.
- Testar leitura sem conta, envio sem consentimentos opcionais, termos inexistentes, expiração de pedido, confirmação concorrente, exclusão e ausência de conta ativa para e-mail não verificado.

## Fontes e revisão

- [LGPD](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709.htm): arts. 6º–11, 14 e 18–19.
- [Código Civil](https://www.planalto.gov.br/ccivil_03/leis/2002/l10406compilada.htm): arts. 104, 107, 422 e 425.

Consultar advogado especializado para revisar o acordo e sua adequação ao responsável real, ao público/menores e ao risco de inferência religiosa antes de publicação. A escolha de arquitetura/base está definida neste documento; revisão de fatos e de legalidade da execução não está dispensada.
