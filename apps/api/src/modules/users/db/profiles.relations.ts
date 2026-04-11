import { relations } from "drizzle-orm";
import { profiles } from "./profiles.table";
import { users } from "./users.table";

export const profileRelations = relations(profiles, ({ one }) => ({
  user: one(users, {
    fields: [profiles.userId],
    references: [users.id],
  }),
}));
