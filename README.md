# Versum

Versum é um leitor bíblico para quem quer ler a Bíblia inteira no próprio ritmo.
O web oferece leitura e links compartilháveis. O aplicativo Android mantém
conteúdo e progresso disponíveis sem conexão.

O primeiro lançamento cobre autenticação por magic link, progresso sincronizado
entre dispositivos, leitura offline no Android, lembretes diários e
compartilhamento de capítulos e versículos. Não haverá feed, likes, ranking ou
perfil social.

## Estado do projeto

O repositório está em pré-alpha. A arquitetura e o produto foram definidos; a
implementação começou pelo corpus bíblico canônico e pelas ferramentas de
qualidade do repositório. A API do servidor ainda não existe.

## Estrutura

```text
api/             módulo Go do servidor (ainda vazio)
bible/           corpus bíblico canônico versionado e sua ferramenta de geração
tools/vaultlint  verificador de metadados e links do vault Obsidian
web/             leitor online em Next.js (planejado)
mobile/          aplicativo Android em React Native/Expo (planejado)
infra/           ambiente local e deploy (planejado)
Obsidian Vault/  documentação do projeto
```

Os clientes serão projetos independentes. A API Go publicará o contrato que web
e mobile consomem.

## Documentação

O vault Obsidian concentra a documentação humana do projeto:

- [Produto](Obsidian%20Vault/PRD.md)
- [Arquitetura](Obsidian%20Vault/Docs/Architecture/_Index.md)
- [Decisões](Obsidian%20Vault/Docs/Decisions/_Index.md)
- [Regras](Obsidian%20Vault/Rules/_Index.md)
- [Planos](Obsidian%20Vault/Plans/_Index.md)

## Princípios técnicos

- Go organiza o backend por funcionalidade, com casos de uso e adapters.
- PostgreSQL é a fonte de verdade; Redis atende cache, rate limit e locks curtos.
- O mobile registra eventos de leitura localmente e os sincroniza de forma
  idempotente.
- Sessões, tokens e dados pessoais recebem tratamento restritivo desde o início.

## Contribuição

Leia [CONTRIBUTING.md](CONTRIBUTING.md) antes de abrir uma issue ou pull request.
Mudanças de produto e arquitetura devem atualizar o vault quando alterarem uma
decisão que continuará relevante.

Rode `npm install` na raiz uma vez para instalar os hooks de commit do Husky
(bloqueiam coautoria de agentes de IA e rodam as verificações do repositório
antes de cada commit).

## Licença

O código é disponibilizado sob a [licença MIT](LICENSE).
