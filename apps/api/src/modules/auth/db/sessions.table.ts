import { pgTable, uuid, text, timestamp } from "drizzle-orm/pg-core";

export const sessions = pgTable("sessions", {
  id: uuid("id").primaryKey().defaultRandom(),
  public_id: uuid("public_id").notNull().defaultRandom().unique(),
  user_id: uuid("user_id").notNull(),
  ip: text("ip").notNull(),
  created_at: timestamp("created_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  updated_at: timestamp("updated_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  expires_at: timestamp("expires_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  revoked_at: timestamp("revoked_at", {
    precision: 3,
    withTimezone: true,
  }),
  user_agent: text("user_agent").notNull(),
  token_hash: text("token_hash").notNull().unique(),
});
