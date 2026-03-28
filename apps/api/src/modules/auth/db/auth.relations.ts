import { relations } from "drizzle-orm";
import { sessions } from "./sessions.table.ts";
import { users } from "../../../infrastructure/db/schema.ts";

export const authRelations = relations(sessions, ({ one }) => ({
  user: one(users, {
    fields: [sessions.userId],
    references: [users.id],
  }),
}));
