import { index, pgTable, timestamp, uuid, varchar } from "drizzle-orm/pg-core";

export const profiles = pgTable(
  "profiles",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    userId: uuid("user_id").notNull().unique(),
    username: varchar("username", { length: 50 }).notNull().unique(),
    name: varchar("name", { length: 100 }).notNull(),
    bio: varchar("bio", { length: 500 }),
    pictureUrl: varchar("picture_url", { length: 500 }),
    createdAt: timestamp("created_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
    updatedAt: timestamp("updated_at", {
      precision: 3,
      withTimezone: true,
    })
      .notNull()
      .defaultNow(),
  },
  (table) => [
    index("profiles_username_idx").on(table.username),
    index("profiles_user_id_idx").on(table.userId),
  ],
);
