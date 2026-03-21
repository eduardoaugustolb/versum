import { index, pgTable, timestamp, uuid, varchar } from "drizzle-orm/pg-core";

export const magicLinks = pgTable(
  "magic_links",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    public_id: uuid("public_id").notNull().unique().defaultRandom(),
    email: varchar("email", { length: 255 }).notNull(),
    token_hash: varchar("token_hash", { length: 500 }).notNull().unique(),
    used_at: timestamp("used_at", { precision: 3, withTimezone: true }),
    invalidated_at: timestamp("invalidated_at", {
      precision: 3,
      withTimezone: true,
    }),
    expires_at: timestamp("expires_at", {
      precision: 3,
      withTimezone: true,
    }).notNull(),
    created_at: timestamp("created_at", { precision: 3, withTimezone: true })
      .defaultNow()
      .notNull(),
  },
  (table) => [
    index("magic_links_email_public_id_idx").on(table.email, table.public_id),
    index("magic_links_expires_at_idx").on(table.expires_at),
  ],
);
