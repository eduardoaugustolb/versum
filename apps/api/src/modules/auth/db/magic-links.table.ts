import { index, pgTable, timestamp, uuid, varchar } from "drizzle-orm/pg-core";

export const magicLinks = pgTable(
  "magic_links",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    publicId: uuid("public_id").notNull().unique().defaultRandom(),
    email: varchar("email", { length: 255 }).notNull(),
    tokenHash: varchar("token_hash", { length: 500 }).notNull().unique(),
    usedAt: timestamp("used_at", { precision: 3, withTimezone: true }),
    invalidatedAt: timestamp("invalidated_at", {
      precision: 3,
      withTimezone: true,
    }),
    expiresAt: timestamp("expires_at", {
      precision: 3,
      withTimezone: true,
    }).notNull(),
    createdAt: timestamp("created_at", { precision: 3, withTimezone: true })
      .defaultNow()
      .notNull(),
  },
  (table) => [
    index("magic_links_email_public_id_idx").on(table.email, table.publicId),
    index("magic_links_expires_at_idx").on(table.expiresAt),
  ],
);
