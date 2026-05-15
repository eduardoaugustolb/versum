import { beforeEach, describe, expect, it, vi } from "vitest";
import type { ConsentLog } from "../repositories/consent-logs.types.repository";
import { ConsentLogServiceV1 } from "./consent-log.v1.service";

describe("ConsentLogServiceV1", () => {
  let service: ConsentLogServiceV1;

  const mockConsentLog: ConsentLog = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    userId: "123e4567-e89b-12d3-a456-426614174001",
    purpose: "profile_content",
    granted: true,
    ip: "127.0.0.1",
    userAgent: "Mozilla/5.0",
    createdAt: new Date("2024-01-01T00:00:00Z"),
  };

  const createMockRepository = () => ({
    config: { dbInstance: {} },
    createConsentLog: vi.fn<() => Promise<ConsentLog>>(),
    createConsentLogs: vi.fn<() => Promise<ConsentLog[]>>(),
    getConsentLogsByUserId: vi.fn<() => Promise<ConsentLog[]>>(),
    hasConsent: vi.fn<() => Promise<boolean>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new ConsentLogServiceV1({
      repository:
        mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
    });
    vi.clearAllMocks();
  });

  describe("recordConsents", () => {
    it("should create multiple consent logs", async () => {
      const mockRepository = createMockRepository();
      const logs = [
        { ...mockConsentLog, purpose: "profile_content" },
        { ...mockConsentLog, purpose: "annotations" },
      ];
      mockRepository.createConsentLogs.mockResolvedValue(logs);
      service = new ConsentLogServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
      });

      const result = await service.recordConsents({
        userId: mockConsentLog.userId,
        consents: [
          { purpose: "profile_content", granted: true },
          { purpose: "annotations", granted: true },
        ],
        ip: "127.0.0.1",
        userAgent: "Mozilla/5.0",
      });

      expect(result).toEqual(logs);
      expect(result).toHaveLength(2);
      expect(mockRepository.createConsentLogs).toHaveBeenCalledWith(
        expect.arrayContaining([
          expect.objectContaining({ purpose: "profile_content", granted: true }),
          expect.objectContaining({ purpose: "annotations", granted: true }),
        ]),
      );
    });

    it("should throw error for invalid purpose", async () => {
      await expect(
        service.recordConsents({
          userId: mockConsentLog.userId,
          consents: [{ purpose: "invalid_purpose", granted: true }],
          ip: "127.0.0.1",
          userAgent: "Mozilla/5.0",
        }),
      ).rejects.toThrow("Invalid consent purpose: invalid_purpose");
    });
  });

  describe("getUserConsents", () => {
    it("should return all consents for a user", async () => {
      const mockRepository = createMockRepository();
      const logs = [mockConsentLog];
      mockRepository.getConsentLogsByUserId.mockResolvedValue(logs);
      service = new ConsentLogServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
      });

      const result = await service.getUserConsents({ userId: mockConsentLog.userId });

      expect(result).toEqual(logs);
      expect(mockRepository.getConsentLogsByUserId).toHaveBeenCalledWith({
        userId: mockConsentLog.userId,
      });
    });

    it("should return empty array when no consents", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getConsentLogsByUserId.mockResolvedValue([]);
      service = new ConsentLogServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
      });

      const result = await service.getUserConsents({ userId: "nonexistent-id" });

      expect(result).toEqual([]);
    });
  });

  describe("hasConsent", () => {
    it("should return true when consent exists and is granted", async () => {
      const mockRepository = createMockRepository();
      mockRepository.hasConsent.mockResolvedValue(true);
      service = new ConsentLogServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
      });

      const result = await service.hasConsent({
        userId: mockConsentLog.userId,
        purpose: "profile_content",
      });

      expect(result).toBe(true);
    });

    it("should return false when consent is not granted", async () => {
      const mockRepository = createMockRepository();
      mockRepository.hasConsent.mockResolvedValue(false);
      service = new ConsentLogServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/consent-logs.repository").ConsentLogsRepository,
      });

      const result = await service.hasConsent({
        userId: mockConsentLog.userId,
        purpose: "annotations",
      });

      expect(result).toBe(false);
    });
  });
});
