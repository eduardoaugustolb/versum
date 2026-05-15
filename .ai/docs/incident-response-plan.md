# Plano de Resposta a Incidentes de Segurança

> **DPO:** dpo@versum.app
>
> **Prazo de notificação à ANPD:** 2 dias úteis (Art. 48 LGPD)

---

## 1. Definição de Violação de Dados Pessoais

Uma violação de dados pessoais é um incidente de segurança que resulta em:

- **Destruição** acidental ou ilícita de dados pessoais
- **Perda** acidental ou ilícita de dados pessoais
- **Alteração** acidental ou ilícita de dados pessoais
- **Divulgação** não autorizada de dados pessoais
- **Acesso** não autorizado a dados pessoais

## 2. Classificação de Incidentes

| Nível | Descrição | Exemplos | Prazo de Resposta |
|-------|-----------|----------|-------------------|
| **Baixo** | Sem exposição de dados pessoais | Tentativa de força bruta frustrada, scan de porta | 24h para investigar |
| **Médio** | Exposição limitada de dados não sensíveis | Acesso indevido a logs, vazamento de IP | 12h para conter |
| **Alto** | Exposição de dados pessoais identificáveis | Vazamento de emails, acesso a tokens de sessão | 4h para conter |
| **Crítico** | Exposição massiva de dados sensíveis | Acesso ao banco de dados completo, vazamento de token hashes | Imediato |

## 3. Fluxo de Detecção e Resposta

### 3.1 Detecção

- **Automática:**
  - Rate limiting excedido → alerta de múltiplas tentativas
  - Erro 500 no módulo de autenticação → alerta de possível incidente
  - Tentativas de acesso com tokens inválidos → alerta de força bruta

- **Manual:**
  - Reporte de usuários
  - Reporte de parceiros (Resend, hospedagem)
  - Descoberta em auditoria de logs

### 3.2 Triagem (30 min)

1. Confirmar se o incidente envolve dados pessoais
2. Classificar o nível (Baixo/Médio/Alto/Crítico)
3. Notificar o DPO imediatamente se nível ≥ Médio
4. Abrir registro do incidente no repositório

### 3.3 Contenção (2h-24h)

| Ação | Responsável |
|------|-------------|
| Isolar sistemas afetados | Engenharia |
| Rotacionar chaves/tokens comprometidos | Engenharia |
| Bloquear IPs maliciosos | DevOps |
| Fazer snapshot forense (logs, DB) | Engenharia |
| Preservar evidências | Engenharia |

### 3.4 Investigação (48h)

1. Identificar causa raiz
2. Determinar escopo: quais dados foram acessados/expostos
3. Identificar titulares afetados
4. Documentar cronologia do incidente

### 3.5 Notificação (2 dias úteis)

**Para a ANPD (Art. 48):**
- Natureza dos dados pessoais afetados
- Circunstâncias da violação
- Medidas de contenção adotadas
- Riscos potenciais aos titulares
- Contato do DPO

**Para os titulares (quando necessário):**
- Descrição clara do incidente
- Dados pessoais afetados
- Medidas de proteção adotadas
- Recomendações para o titular (ex.: alterar senhas em outros serviços)
- Contato do DPO

## 4. Template de Notificação à ANPD

```
Assunto: Notificação de Violação de Dados Pessoais — Versum

Data do incidente: [DATA]
Data desta notificação: [DATA]

1. Descrição do incidente:
   [Descrever o que aconteceu, como foi detectado]

2. Dados pessoais afetados:
   [Listar tipos de dados: email, nome, IP, etc.]

3. Titulares afetados:
   [Número aproximado de usuários]

4. Circunstâncias:
   - Causa: [causa raiz]
   - Data/hora da violação: [timestamp]
   - Data/hora da detecção: [timestamp]
   - Sistema comprometido: [sistema]

5. Medidas de contenção adotadas:
   [Listar ações já tomadas]

6. Riscos potenciais:
   [Descrever riscos para os titulares]

7. Contato do DPO:
   Nome: [Nome]
   Email: dpo@versum.app
   Telefone: [telefone]
```

## 5. Responsabilidades

| Função | Responsável | Contato |
|--------|-------------|---------|
| DPO | — | dpo@versum.app |
| Engenharia | Time de desenvolvimento | — |
| DevOps | Time de infraestrutura | — |
| Jurídico | — | — |

## 6. Pós-Incidente

1. Análise de causa raiz (post-mortem)
2. Implementação de medidas corretivas
3. Atualização deste plano
4. Relatório final arquivado

## 7. Alertas Automáticos (Discord Webhook)

Em caso de:
- Múltiplas tentativas de login falhas (rate limit excedido)
- Tentativas de acesso a tokens inválidos
- Erro 500 inesperado no módulo de autenticação

O sistema deve enviar alerta para o canal de segurança no Discord indicando:
- Timestamp
- IP do solicitante
- Rota afetada
- Código do erro
