import { relations } from "drizzle-orm";
import { sessions } from "../../../infrastructure/db/schema";
import { likes } from "../../interactions/db/likes.table";
import { marks } from "../../interactions/db/marks.table";
import { journeyReadings } from "../../interactions/db/readings.table";
import { profiles } from "./profiles.table";
import { users } from "./users.table";

export const userRelations = relations(users, ({ one, many }) => ({
  profile: one(profiles, {
    fields: [users.id],
    references: [profiles.userId],
  }),
  readings: many(journeyReadings),
  likes: many(likes),
  marks: many(marks),
  sessions: many(sessions),
}));
