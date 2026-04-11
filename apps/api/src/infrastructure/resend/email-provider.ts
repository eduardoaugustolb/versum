import { type CreateEmailOptions, Resend } from "resend";
import { env } from "../../utils/env/index.ts";

export class EmailProvider {
  private readonly resend: Resend;

  constructor({ resend }: { resend?: Resend } = {}) {
    this.resend = resend ?? new Resend(env.RESEND_API_KEY);
  }

  async sendEmail(payload: Omit<CreateEmailOptions, "from">) {
    await this.resend.emails.send({
      ...payload,
      from: "no-reply@send.eduardoaugusto.is-a.dev",
    } as CreateEmailOptions);
  }
}
