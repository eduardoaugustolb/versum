import { pgEnum } from "drizzle-orm/pg-core";

export const readingModeEnum = pgEnum("reading_mode", ["DISCOVERY", "JOURNEY"]);
