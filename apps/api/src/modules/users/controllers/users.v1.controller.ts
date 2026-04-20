import type { Context } from "hono";
import {
  BadRequestError,
  ConflictError,
  NotFoundError,
} from "../../../utils/app/errors/index";
import type { Session } from "../../auth/repositories/auth.types.repository";
import { UserServiceV1 } from "../services/user.v1.service";
import type { User } from "../repositories/user.types.repository";
import type { Profile } from "../repositories/profile.types.repository";

const userService = new UserServiceV1();

export class UsersControllerV1 {
  private readonly service: UserServiceV1;

  constructor({ service }: { service?: UserServiceV1 } = {}) {
    this.service = service ?? userService;
  }

  getAutheticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;

    const user = await this.service.getUserByIdWithProfile({
      id: session.userId,
    });

    const profile = user.profile;

    const hydratedUser: User & { profile?: Profile } = Object.create(user);

    delete hydratedUser.profile;

    if (!user || !hydratedUser) {
      throw new NotFoundError("User not found");
    }

    const onboardingIsCompleted = !!profile;

    return c.json(
      {
        user: {
          email: user.email,
        },
        onboardingIsCompleted,
      },
      200,
    );
  };

  updateAuthenticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = (await c.req.json()) as { email?: string };

    // TODO: Add new email validation and ip adress validation

    if (!body.email) {
      throw new BadRequestError("Email is required");
    }

    try {
      const user = await this.service.updateUserEmail({
        id: session.userId,
        email: body.email,
      });

      return c.json(
        {
          user: {
            email: user.email,
          },
        },
        200,
      );
    } catch (error) {
      if (error instanceof Error && error.message === "Email already in use") {
        throw new ConflictError("Email already in use");
      }
      throw error;
    }
  };
}
