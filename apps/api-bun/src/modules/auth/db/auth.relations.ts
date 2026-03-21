import { relations } from "drizzle-orm";
import { sessions } from "./sessions.table.ts";
import { users } from "../../../infrastructure/db/schema.ts";

export const authRelations = relations(sessions, ({ one }) => ({
  user_id: one(users, {
    fields: [sessions.user_id],
    references: [users.id],
    relationName: "user_id",
  }),
}));
