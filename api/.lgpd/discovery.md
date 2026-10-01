# Descoberta — 2026-10-01

Escopo solicitado: API atual e `../Obsidian Vault`. Cenário A confirmado. Leitura de código, migrations, testes, planos, PRD, regras e decisões; nenhuma consulta a dados reais, ambientes remotos, segredos ou contratos privados. `prowl` não está disponível nesta sessão; inspeção feita com leitura de arquivos e `rg`.

Referência Git: `5afbb622430bc7d3539c4d87c8a83516a24ece7e`. O working tree continha alterações prévias em sincronização, índices e plano de autenticação, além das novas docs de privacidade e do teste de consumo. Elas integram o retrato auditado e foram preservadas.

Stack: Go, PostgreSQL, Redis, roteador HTTP. Dependências locais não comprovam contratação de operadores. Catálogo público; solicitação de magic link implementada; consumo, envio por provedor e sessões HTTP ainda incompletos. O PRD prevê Android/web, progresso sincronizado e lembretes; esses fluxos futuros não foram tratados como coleta atual.

As próprias docs declaram pré-alpha e sem produção. Isso não confirma que todo dado de desenvolvimento seja sintético; tal confirmação cabe ao controlador. Frontend/mobile executáveis e infraestrutura produtiva estão fora do escopo disponibilizado.

## Evidências e limites

- [Código](code-evidence.md): controles e ausências no escopo.
- [Vault](docs-evidence.md): divergências, decisões e pendências.
- [Inventário](data-map.md): atual e futuro separados.
- [Bases propostas](legal-basis.md): decisões ainda não aprovadas.
- [Gaps](gaps.md): prioridade, norma e critério de aceite.

Pendente de comprovação externa: identidade do controlador, responsáveis, porte/receita/grupo e enquadramento ATPP, usuários/faixas etárias, dados reais em testes, operadores contratados, localização dos dados/backups, TLS e controle de acesso, publicação de informações ao titular e histórico de incidentes. Ausência no repositório não comprova inexistência fora dele.
