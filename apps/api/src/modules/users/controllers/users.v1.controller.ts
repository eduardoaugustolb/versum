import { Context } from "hono";
import { UserServiceV1 } from "../services/user.v1.service.ts";
import { Session } from "../../auth/repositories/auth.types.repository.ts";
import { User } from "../repositories/user.types.repository.ts";
import { NotFoundError } from "../../../utils/app/errors/index.ts";

export class UsersControllerV1 {
  private readonly service: UserServiceV1;

  constructor({ service }: { service?: UserServiceV1 } = {}) {
    this.service = service ?? new UserServiceV1();
  }

  getAutheticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;

    const user = await this.service.getUserById({
      id: session.user_id,
    });

    if (!user) {
      throw new NotFoundError("User not found");
    }

    return c.json({ user }, 200);
  };

  updateAuthenticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = (await c.req.json()) as Partial<
      Omit<User, "id" | "email" | "createdAt">
    >;

    const updatedUser = await this.service.updateUser({
      id: session.user_id,
      ...body,
    });

    return c.json({ user: updatedUser }, 200);
  };

  getUserByUsername = async (c: Context) => {
    const rawUsername = c.req.param("username");
    const username = rawUsername.startsWith("@")
      ? rawUsername.slice(1)
      : rawUsername;

    const user = await this.service.getUserByUsername({
      username,
    });

    if (!user) {
      throw new NotFoundError("User not found");
    }

    return c.json({ user }, 200);
  };
}
