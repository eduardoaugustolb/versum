import { relations } from "drizzle-orm";
import { users } from "@/infrastructure/db/schema";
import { consentLogs } from "./consent-logs.table";

export const consentLogsRelations = relations(consentLogs, ({ one }) => ({
  user: one(users, { fields: [consentLogs.userId], references: [users.id] }),
}));
