import type { Context } from "hono";
import { setCookie } from "hono/cookie";
import {
  BadRequestError,
  ConflictError,
} from "../../../utils/app/errors/index";
import type { Session } from "../../auth/repositories/auth.types.repository";
import { UserServiceV1 } from "../services/user.v1.service";

const isSecure = Bun.env.COOKIE_SECURE === "true";
const cookieOptions = {
  httpOnly: true,
  secure: isSecure,
  sameSite: (isSecure ? "strict" : "lax") as "strict" | "lax",
  path: "/",
};

const userService = new UserServiceV1();

export class UsersControllerV1 {
  private readonly service: UserServiceV1;

  constructor({ service }: { service?: UserServiceV1 } = {}) {
    this.service = service ?? userService;
  }

  getAuthenticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;

    const user = await this.service.getUserByIdWithProfile({
      id: session.userId,
    });

    const onboardingIsCompleted = !!user.profile;

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

  exportUserData = async (c: Context) => {
    const session = c.get("session") as Session;

    const data = await this.service.exportUserData({ id: session.userId });

    return c.json(data, 200);
  };

  deleteAuthenticatedUser = async (c: Context) => {
    const session = c.get("session") as Session;

    await this.service.deleteUser({ id: session.userId });

    const cookieName = isSecure ? "__Host-session" : "session";
    setCookie(c, cookieName, "", {
      ...cookieOptions,
      expires: new Date(0),
    });

    return c.body(null, 204);
  };
}
