import { relations } from "drizzle-orm";
import { consentLogs, sessions } from "../../../infrastructure/db/schema";
import { likes } from "../../interactions/db/likes.table";
import { marks } from "../../interactions/db/marks.table";
import { journeyReadings } from "../../interactions/db/readings.table";
import { discoveryReadings } from "../../interactions/db/verse-readings.table";
import { profiles } from "./profiles.table";
import { users } from "./users.table";

export const userRelations = relations(users, ({ one, many }) => ({
  profile: one(profiles, {
    fields: [users.id],
    references: [profiles.userId],
  }),
  readings: many(journeyReadings),
  discoveryReadings: many(discoveryReadings),
  likes: many(likes),
  marks: many(marks),
  sessions: many(sessions),
  consentLogs: many(consentLogs),
}));
