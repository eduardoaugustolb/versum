import {
  BadRequestError,
  ConflictError,
  NotFoundError,
  ForbiddenError,
} from "../../../utils/app/errors/index";
import { ProfileRepository } from "../repositories/profile.repository";
import type {
  CreateProfileParams,
  Profile,
  UpdateProfileParams,
} from "../repositories/profile.types.repository";
import { ConsentLogsRepository } from "../../consent-logs/repositories/consent-logs.repository";

const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/;
const MAX_NAME_LENGTH = 100;
const MAX_BIO_LENGTH = 500;
const MAX_PICTURE_URL_LENGTH = 500;
const MIN_USERNAME_LENGTH = 3;
const MAX_USERNAME_LENGTH = 50;

export class ProfileServiceV1 {
  private readonly repository: ProfileRepository;
  private readonly consentLogsRepository: ConsentLogsRepository;

  constructor({
    repository,
    consentLogsRepository,
  }: {
    repository?: ProfileRepository;
    consentLogsRepository?: ConsentLogsRepository;
  } = {}) {
    this.repository = repository ?? new ProfileRepository();
    this.consentLogsRepository = consentLogsRepository ?? new ConsentLogsRepository();
  }

  private sanitizeHtml(input: string): string {
    return input
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#x27;")
      .replace(/\//g, "&#x2F;");
  }

  private validateUsername(username: string): void {
    const trimmed = username.trim().toLowerCase();

    if (trimmed.length < MIN_USERNAME_LENGTH) {
      throw new BadRequestError(
        `Username must be at least ${MIN_USERNAME_LENGTH} characters`,
      );
    }

    if (trimmed.length > MAX_USERNAME_LENGTH) {
      throw new BadRequestError(
        `Username must not exceed ${MAX_USERNAME_LENGTH} characters`,
      );
    }

    if (!USERNAME_REGEX.test(trimmed)) {
      throw new BadRequestError(
        "Username can only contain letters, numbers, and underscores",
      );
    }
  }

  private validateName(name: string): void {
    const trimmed = name.trim();

    if (trimmed.length === 0) {
      throw new BadRequestError("Name is required");
    }

    if (trimmed.length > MAX_NAME_LENGTH) {
      throw new BadRequestError(
        `Name must not exceed ${MAX_NAME_LENGTH} characters`,
      );
    }
  }

  private validateBio(bio: string | null | undefined): void {
    if (!bio) return;

    const trimmed = bio.trim();

    if (trimmed.length > MAX_BIO_LENGTH) {
      throw new BadRequestError(
        `Bio must not exceed ${MAX_BIO_LENGTH} characters`,
      );
    }
  }

  private validatePictureUrl(url: string | null | undefined): void {
    if (!url) return;

    if (url.length > MAX_PICTURE_URL_LENGTH) {
      throw new BadRequestError(
        `Picture URL must not exceed ${MAX_PICTURE_URL_LENGTH} characters`,
      );
    }

    let parsed: URL;
    try {
      parsed = new URL(url);
    } catch {
      throw new BadRequestError("Picture URL must be a valid URL");
    }

    if (parsed.protocol !== "https:") {
      throw new BadRequestError("Picture URL must use HTTPS");
    }
  }

  private sanitizeAndValidate(
    params: UpdateProfileParams,
  ): UpdateProfileParams {
    const sanitized: UpdateProfileParams = {};

    if (params.username !== undefined) {
      this.validateUsername(params.username);
      sanitized.username = params.username.trim().toLowerCase();
    }

    if (params.name !== undefined) {
      this.validateName(params.name);
      sanitized.name = this.sanitizeHtml(params.name.trim());
    }

    if (params.bio !== undefined) {
      this.validateBio(params.bio);
      sanitized.bio = params.bio ? this.sanitizeHtml(params.bio.trim()) : null;
    }

    if (params.pictureUrl !== undefined) {
      this.validatePictureUrl(params.pictureUrl);
      sanitized.pictureUrl = params.pictureUrl?.trim() || null;
    }

    return sanitized;
  }

  async createProfile(params: CreateProfileParams): Promise<Profile> {
    const hasConsent = await this.consentLogsRepository.hasConsent({
      userId: params.userId,
      purpose: "profile_content",
    });

    if (!hasConsent) {
      throw new ForbiddenError(
        "Consentimento para armazenar conteúdo do perfil não foi concedido",
      );
    }

    const existingProfile = await this.repository.findByUserId({
      userId: params.userId,
    });

    if (existingProfile) {
      throw new ConflictError("Profile already exists");
    }

    const sanitizedParams = {
      ...params,
      username: params.username.trim().toLowerCase(),
      name: this.sanitizeHtml(params.name.trim()),
      bio: params.bio ? this.sanitizeHtml(params.bio.trim()) : null,
      pictureUrl: params.pictureUrl?.trim() || null,
    };

    this.validateUsername(sanitizedParams.username);
    this.validateName(sanitizedParams.name);
    this.validateBio(sanitizedParams.bio);
    this.validatePictureUrl(sanitizedParams.pictureUrl);

    const profileWithSelectedUsername = await this.repository.existsByUsername({
      username: sanitizedParams.username,
    });

    if (profileWithSelectedUsername.exists) {
      throw new ConflictError("Username already in use");
    }

    return await this.repository.create(sanitizedParams);
  }

  async getProfileByUserId({
    userId,
  }: {
    userId: string;
  }): Promise<Profile | null> {
    return await this.repository.findByUserId({ userId });
  }

  async getProfileByUsername({
    username,
  }: {
    username: string;
  }): Promise<Profile | null> {
    return await this.repository.findByUsername({ username });
  }

  async getProfileById({ id }: { id: string }): Promise<Profile> {
    const profile = await this.repository.findById({ id });

    if (!profile) {
      throw new NotFoundError("Profile not found");
    }

    return profile;
  }

  async updateProfile(
    params: UpdateProfileParams & { userId: string },
  ): Promise<Profile> {
    const profile = await this.repository.findByUserId({
      userId: params.userId,
    });

    if (!profile) {
      throw new NotFoundError("Profile not found");
    }

    const sanitized = this.sanitizeAndValidate(params);

    if (sanitized.username) {
      const profileWithSelectedUsername = await this.repository.existsByUsername({
        username: sanitized.username,
      });

      if (profileWithSelectedUsername.exists && profileWithSelectedUsername.profileId !== profile.id) {
        throw new ConflictError("Username already in use");
      }
    }

    return await this.repository.update({
      ...sanitized,
      id: profile.id,
    });
  }

}
