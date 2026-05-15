import { beforeEach, describe, expect, it, vi } from "vitest";
import { PurgeService } from "./purge.service";

describe("PurgeService", () => {
  let service: PurgeService;

  const createMockAuthRepository = () => ({
    db: {},
    deleteExpiredMagicLinks: vi.fn<() => Promise<number>>(),
    deleteExpiredSessions: vi.fn<() => Promise<number>>(),
  });

  const createService = ({
    mockAuthRepository,
  }: {
    mockAuthRepository: ReturnType<typeof createMockAuthRepository>;
  }) =>
    new PurgeService({
      authRepository:
        mockAuthRepository as unknown as import("../repositories/auth.repository").AuthRepository,
    });

  beforeEach(() => {
    const mockAuthRepository = createMockAuthRepository();
    service = createService({ mockAuthRepository });
    vi.clearAllMocks();
  });

  describe("purgeExpiredMagicLinks", () => {
    it("should delete expired magic links and return count", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredMagicLinks.mockResolvedValue(5);
      service = createService({ mockAuthRepository });

      const result = await service.purgeExpiredMagicLinks();

      expect(result).toEqual({ deleted: 5 });
      expect(mockAuthRepository.deleteExpiredMagicLinks).toHaveBeenCalled();
    });

    it("should return 0 when no expired magic links", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredMagicLinks.mockResolvedValue(0);
      service = createService({ mockAuthRepository });

      const result = await service.purgeExpiredMagicLinks();

      expect(result).toEqual({ deleted: 0 });
    });
  });

  describe("purgeExpiredSessions", () => {
    it("should delete expired sessions and return count", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredSessions.mockResolvedValue(3);
      service = createService({ mockAuthRepository });

      const result = await service.purgeExpiredSessions();

      expect(result).toEqual({ deleted: 3 });
      expect(mockAuthRepository.deleteExpiredSessions).toHaveBeenCalled();
    });

    it("should return 0 when no expired sessions", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredSessions.mockResolvedValue(0);
      service = createService({ mockAuthRepository });

      const result = await service.purgeExpiredSessions();

      expect(result).toEqual({ deleted: 0 });
    });
  });

  describe("runDailyPurge", () => {
    it("should run both purges and return combined totals", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredMagicLinks.mockResolvedValue(5);
      mockAuthRepository.deleteExpiredSessions.mockResolvedValue(3);
      service = createService({ mockAuthRepository });

      const result = await service.runDailyPurge();

      expect(result).toEqual({ magicLinks: 5, sessions: 3 });
      expect(mockAuthRepository.deleteExpiredMagicLinks).toHaveBeenCalled();
      expect(mockAuthRepository.deleteExpiredSessions).toHaveBeenCalled();
    });

    it("should handle both returning 0", async () => {
      const mockAuthRepository = createMockAuthRepository();
      mockAuthRepository.deleteExpiredMagicLinks.mockResolvedValue(0);
      mockAuthRepository.deleteExpiredSessions.mockResolvedValue(0);
      service = createService({ mockAuthRepository });

      const result = await service.runDailyPurge();

      expect(result).toEqual({ magicLinks: 0, sessions: 0 });
    });
  });
});
