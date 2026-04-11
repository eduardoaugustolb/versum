import { relations } from "drizzle-orm";
import { users } from "../../../infrastructure/db/schema.ts";
import { sessions } from "./sessions.table.ts";

export const authRelations = relations(sessions, ({ one }) => ({
  user: one(users, {
    fields: [sessions.userId],
    references: [users.id],
  }),
}));
