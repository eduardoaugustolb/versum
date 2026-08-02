# Contribuindo

Obrigado pelo interesse no Versum.

## Antes de começar

Leia o [PRD](Obsidian%20Vault/PRD.md), a
[arquitetura](Obsidian%20Vault/Docs/Architecture/_Index.md) e as
[regras do projeto](Obsidian%20Vault/Rules/_Index.md). Se a mudança alterar um
limite arquitetural ou uma regra de produto, proponha a decisão antes de
implementar.

## Fluxo de trabalho

1. Abra uma issue para trabalho que não seja uma correção pequena.
2. Crie uma branch com um nome descritivo: `feat/`, `fix/`, `refactor/` ou
   `docs/`.
3. Mantenha cada pull request focado em uma entrega verificável.
4. Inclua testes para regras de domínio, integração com dados ou comportamento
   visível que a mudança introduzir.
5. Atualize o vault quando a mudança criar ou substituir conhecimento estável.
6. Descreva na pull request o problema, a solução, a forma de testar e qualquer
   risco de rollback.

## Padrões

- Código do backend deve preservar limites claros entre casos de uso e adapters.
- O banco de dados protege invariantes por transações, índices e restrições; o
  cliente não é fonte de verdade para dados sincronizados.
- Tokens, segredos, e-mails e URLs assinadas não entram em logs, testes ou
  commits.
- Não adicione dependências para resolver uma tarefa que a biblioteca padrão ou
  a estrutura já adotada resolve bem.
- Evite reformatar arquivos sem relação com a mudança.

## Commits e pull requests

Use mensagens no formato Conventional Commits:

```text
feat(api): add reading event sync
fix(mobile): retry pending outbox events
docs: clarify magic link flow
```

Pull requests devem explicar decisões fora do óbvio. Se houver uma alteração de
produto ou arquitetura duradoura, inclua o link para a nota correspondente no
vault.

## Código de conduta

Ao participar deste projeto, siga o [Código de Conduta](CODE_OF_CONDUCT.md).
