import { relations } from "drizzle-orm";
import { users } from "./users.table.ts";
import { readings } from "../../interactions/db/readings.table.ts";
import { likes } from "../../interactions/db/likes.table.ts";
import { marks } from "../../interactions/db/marks.table.ts";
import { sessions } from "../../../infrastructure/db/schema.ts";

export const userRelations = relations(users, ({ many }) => ({
  readings: many(readings),
  likes: many(likes),
  marks: many(marks),
  sessions: many(sessions),
}));
