import type { Context } from "hono";
import { NotFoundError } from "../../../utils/app/errors/index";
import type { Session } from "../../auth/repositories/auth.types.repository";
import { ProfileServiceV1 } from "../services/profile.v1.service";

export class ProfileControllerV1 {
  private readonly service: ProfileServiceV1;

  constructor({ service }: { service?: ProfileServiceV1 } = {}) {
    this.service = service ?? new ProfileServiceV1();
  }

  createProfile = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = await c.req.json();

    const profile = await this.service.createProfile({
      userId: session.userId,
      ...body,
    });

    return c.json({ profile }, 201);
  };

  getAuthenticatedProfile = async (c: Context) => {
    const session = c.get("session") as Session;

    const profile = await this.service.getProfileByUserId({
      userId: session.userId,
    });

    if (!profile) {
      throw new NotFoundError("Profile not found");
    }

    return c.json({ profile }, 200);
  };

  updateAuthenticatedProfile = async (c: Context) => {
    const session = c.get("session") as Session;
    const body = await c.req.json();

    const profile = await this.service.updateProfile({
      userId: session.userId,
      ...body,
    });

    return c.json({ profile }, 200);
  };

  getProfileByUsername = async (c: Context) => {
    const rawUsername = c.req.param("username");
    const username = rawUsername.startsWith("@")
      ? rawUsername.slice(1)
      : rawUsername;

    const profile = await this.service.getProfileByUsername({ username });

    if (!profile) {
      throw new NotFoundError("Profile not found");
    }

    return c.json({ profile }, 200);
  };
}
