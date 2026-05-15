import { AuthRepository } from "../repositories/auth.repository";

export class PurgeService {
  private readonly authRepository: AuthRepository;

  constructor({
    authRepository,
  }: {
    authRepository?: AuthRepository;
  } = {}) {
    this.authRepository = authRepository ?? new AuthRepository();
  }

  async purgeExpiredMagicLinks(): Promise<{ deleted: number }> {
    const deleted = await this.authRepository.deleteExpiredMagicLinks();
    return { deleted };
  }

  async purgeExpiredSessions(): Promise<{ deleted: number }> {
    const deleted = await this.authRepository.deleteExpiredSessions();
    return { deleted };
  }

  async runDailyPurge(): Promise<{ magicLinks: number; sessions: number }> {
    const [magicLinks, sessions] = await Promise.all([
      this.authRepository.deleteExpiredMagicLinks(),
      this.authRepository.deleteExpiredSessions(),
    ]);

    return {
      magicLinks,
      sessions,
    };
  }
}
