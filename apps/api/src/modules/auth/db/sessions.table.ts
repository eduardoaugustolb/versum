import { pgTable, text, timestamp, uuid } from "drizzle-orm/pg-core";

export const sessions = pgTable("sessions", {
  id: uuid("id").primaryKey().defaultRandom(),
  publicId: uuid("public_id").notNull().defaultRandom().unique(),
  userId: uuid("user_id").notNull(),
  ip: text("ip").notNull(),
  createdAt: timestamp("created_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  updatedAt: timestamp("updated_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  expiresAt: timestamp("expires_at", {
    precision: 3,
    withTimezone: true,
  })
    .defaultNow()
    .notNull(),
  revokedAt: timestamp("revoked_at", {
    precision: 3,
    withTimezone: true,
  }),
  userAgent: text("user_agent").notNull(),
  tokenHash: text("token_hash").notNull().unique(),
});
