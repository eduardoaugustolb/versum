import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Política de Privacidade — Versum",
  description: "Política de privacidade do Versum — saiba como tratamos seus dados pessoais.",
};

const sectionClass = "mb-10";
const headingClass = "text-2xl font-bold tracking-tight font-instrument-serif mb-4";
const subheadingClass = "text-lg font-semibold font-instrument-serif mt-6 mb-2";
const textClass = "text-base leading-relaxed text-neutral-600 dark:text-neutral-400";
const listClass = "list-disc pl-6 space-y-2 text-base leading-relaxed text-neutral-600 dark:text-neutral-400";
const tableClass = "w-full border-collapse text-sm";
const thClass = "border border-neutral-300 dark:border-neutral-700 px-4 py-2 bg-neutral-100 dark:bg-neutral-800 text-left font-medium";
const tdClass = "border border-neutral-300 dark:border-neutral-700 px-4 py-2 text-neutral-600 dark:text-neutral-400";

export default function PrivacyPage() {
  return (
    <main className="mx-auto max-w-3xl px-6 py-16">
      <h1 className="text-4xl font-bold tracking-tight font-instrument-serif mb-2">
        Política de Privacidade
      </h1>
      <p className="text-sm text-neutral-500 mb-10">
        Última atualização: 15 de maio de 2026
      </p>

      <section className={sectionClass}>
        <p className={textClass}>
          Sua privacidade é importante para nós. Esta Política de Privacidade descreve como o Versum coleta,
          usa, armazena e protege seus dados pessoais, em conformidade com a Lei Geral de Proteção de Dados
          Pessoais (LGPD — Lei nº 13.709/2018).
        </p>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>1. Dados Coletados</h2>
        <p className={textClass}>Coletamos os seguintes dados pessoais:</p>
        <div className="overflow-x-auto mt-4">
          <table className={tableClass}>
            <thead>
              <tr>
                <th className={thClass}>Dado</th>
                <th className={thClass}>Finalidade</th>
                <th className={thClass}>Base Legal (LGPD)</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td className={tdClass}>Email</td>
                <td className={tdClass}>Autenticação e comunicação</td>
                <td className={tdClass}>Art. 7º, V (execução de contrato)</td>
              </tr>
              <tr>
                <td className={tdClass}>Nome</td>
                <td className={tdClass}>Identificação no perfil</td>
                <td className={tdClass}>Art. 7º, V (execução de contrato)</td>
              </tr>
              <tr>
                <td className={tdClass}>Username</td>
                <td className={tdClass}>Identificador público</td>
                <td className={tdClass}>Art. 7º, V (execução de contrato)</td>
              </tr>
              <tr>
                <td className={tdClass}>Bio e foto</td>
                <td className={tdClass}>Conteúdo do perfil</td>
                <td className={tdClass}>Art. 7º, I (consentimento)</td>
              </tr>
              <tr>
                <td className={tdClass}>Endereço IP</td>
                <td className={tdClass}>Auditoria e rate limiting</td>
                <td className={tdClass}>Art. 7º, IX (interesse legítimo)</td>
              </tr>
              <tr>
                <td className={tdClass}>User-Agent</td>
                <td className={tdClass}>Identificação do dispositivo</td>
                <td className={tdClass}>Art. 7º, IX (interesse legítimo)</td>
              </tr>
              <tr>
                <td className={tdClass}>Comportamento de leitura</td>
                <td className={tdClass}>Funcionalidade do aplicativo</td>
                <td className={tdClass}>Art. 7º, V (execução de contrato)</td>
              </tr>
              <tr>
                <td className={tdClass}>Anotações</td>
                <td className={tdClass}>Conteúdo gerado pelo usuário</td>
                <td className={tdClass}>Art. 7º, I (consentimento)</td>
              </tr>
              <tr>
                <td className={tdClass}>Favoritos</td>
                <td className={tdClass}>Preferências do usuário</td>
                <td className={tdClass}>Art. 7º, I (consentimento)</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>2. Terceiros com Acesso a Dados</h2>
        <div className="overflow-x-auto mt-4">
          <table className={tableClass}>
            <thead>
              <tr>
                <th className={thClass}>Processador</th>
                <th className={thClass}>Dados Recebidos</th>
                <th className={thClass}>Finalidade</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td className={tdClass}>Resend</td>
                <td className={tdClass}>Email</td>
                <td className={tdClass}>Envio de emails (links mágicos)</td>
              </tr>
              <tr>
                <td className={tdClass}>PostgreSQL (hospedagem)</td>
                <td className={tdClass}>Todos os dados pessoais</td>
                <td className={tdClass}>Armazenamento em banco de dados</td>
              </tr>
              <tr>
                <td className={tdClass}>Redis (hospedagem)</td>
                <td className={tdClass}>Endereço IP (transitório)</td>
                <td className={tdClass}>Cache de rate limiting</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>3. Períodos de Retenção</h2>
        <div className="overflow-x-auto mt-4">
          <table className={tableClass}>
            <thead>
              <tr>
                <th className={thClass}>Dado</th>
                <th className={thClass}>Período de Retenção</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td className={tdClass}>Conta de usuário</td>
                <td className={tdClass}>Até exclusão da conta</td>
              </tr>
              <tr>
                <td className={tdClass}>Sessões revogadas</td>
                <td className={tdClass}>90 dias após revogação</td>
              </tr>
              <tr>
                <td className={tdClass}>Links mágicos expirados</td>
                <td className={tdClass}>30 dias após uso ou 7 dias se abandonados</td>
              </tr>
              <tr>
                <td className={tdClass}>Registros de consentimento</td>
                <td className={tdClass}>Até exclusão da conta</td>
              </tr>
              <tr>
                <td className={tdClass}>Endereço IP (cache)</td>
                <td className={tdClass}>TTL natural do rate limiter</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>4. Seus Direitos (LGPD Art. 18)</h2>
        <p className={textClass}>Você pode exercer os seguintes direitos a qualquer momento:</p>
        <ul className={`${listClass} mt-4`}>
          <li>
            <strong>Acessar seus dados:</strong> Utilize{" "}
            <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">GET /api/v1/users/@me/export</code>
          </li>
          <li>
            <strong>Corrigir seus dados:</strong> Atualize seu perfil em{" "}
            <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">PATCH /api/v1/profiles/@me</code>
          </li>
          <li>
            <strong>Excluir sua conta:</strong> Utilize{" "}
            <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">DELETE /api/v1/users/@me</code>
          </li>
          <li>
            <strong>Revogar consentimento:</strong> Utilize{" "}
            <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">POST /api/v1/consent</code> com{" "}
            <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">granted: false</code>
          </li>
          <li>
            <strong>Portabilidade:</strong> Solicite a exportação dos seus dados (acima)
          </li>
        </ul>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>5. Política de Cookies</h2>
        <p className={textClass}>
          Utilizamos apenas cookies estritamente necessários para o funcionamento do serviço. O cookie de sessão
          (<code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">session</code> ou{" "}
          <code className="text-sm bg-neutral-100 dark:bg-neutral-800 px-1 rounded">__Host-session</code>) é
          um cookie HttpOnly e seguro, utilizado exclusivamente para manter sua sessão autenticada. Não usamos
          cookies de rastreamento, publicidade ou terceiros.
        </p>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>6. Segurança</h2>
        <p className={textClass}>
          Implementamos medidas técnicas e organizacionais para proteger seus dados, incluindo:
        </p>
        <ul className={listClass}>
          <li>Hash de senhas com Argon2</li>
          <li>Cookies HttpOnly e seguros</li>
          <li>HSTS (HTTP Strict Transport Security)</li>
          <li>Rate limiting para proteção contra ataques</li>
          <li>Criptografia em trânsito (TLS)</li>
        </ul>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>7. Encarregado de Dados (DPO)</h2>
        <p className={textClass}>
          Para questões relacionadas ao tratamento de dados pessoais, entre em contato com nosso Encarregado
          de Dados (DPO):
        </p>
        <p className="mt-4 text-base">
          <strong>Email:</strong>{" "}
          <a href="mailto:dpo@versum.app" className="underline hover:text-neutral-900 dark:hover:text-white">
            dpo@versum.app
          </a>
        </p>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>8. Notificação de Violação</h2>
        <p className={textClass}>
          Em caso de violação de dados pessoais que possa acarretar risco ou dano aos titulares,
          notificaremos a Autoridade Nacional de Proteção de Dados (ANPD) no prazo de até 2 (dois)
          dias úteis, conforme Art. 48 da LGPD, e comunicaremos diretamente os titulares afetados.
        </p>
      </section>

      <section className={sectionClass}>
        <h2 className={headingClass}>9. Contato</h2>
        <p className={textClass}>
          Se você tiver dúvidas sobre esta Política de Privacidade ou sobre o tratamento dos seus dados,
          entre em contato conosco:
        </p>
        <ul className={`${listClass} mt-4`}>
          <li><strong>Email:</strong> dpo@versum.app</li>
        </ul>
      </section>
    </main>
  );
}
