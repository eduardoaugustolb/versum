import { envSchema } from "./schema.ts";


export const env = envSchema.parse( Bun.env);
