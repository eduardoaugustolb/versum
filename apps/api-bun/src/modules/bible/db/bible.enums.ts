import { pgEnum } from "drizzle-orm/pg-core";

export const testamentEnum = pgEnum("testament", ["OLD", "NEW"]);
