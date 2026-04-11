import { Resend } from "resend";
import { env } from "../../utils/env/index.ts";

export const resend = new Resend(env.RESEND_API_KEY);
