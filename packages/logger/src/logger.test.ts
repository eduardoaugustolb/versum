import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import { logger } from "./logger.ts";
import { sensitive } from "./sensitive.ts";
import { resetEnvHandlers } from "./env.ts";

describe("logger", () => {
  let consoleLogSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    consoleLogSpy = vi.spyOn(console, "log").mockImplementation(() => {});
  });

  afterEach(() => {
    consoleLogSpy.mockRestore();
    delete process.env.DEBUG;
    delete process.env.BUN_ENV;
  });

  describe("with string config", () => {
    it("should log with info level", () => {
      logger("info", "Hello", "World");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should log with error level", () => {
      logger("error", "Something went wrong");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should log with success level", () => {
      logger("success", "Operation completed");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should log with warn level", () => {
      logger("warn", "Warning message");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should log with debug level", () => {
      logger("debug", "Debug message");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should log with log level", () => {
      logger("log", "Log message");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe("with object config", () => {
    it("should log with level in object config", () => {
      logger({ level: "info" }, "Test message");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should accept custom icon", () => {
      logger({ level: "info", icon: "🔥" }, "Custom icon test");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });

    it("should accept custom color", () => {
      logger({ level: "info", color: "red" }, "Custom color test");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe("debug level behavior", () => {
    it("should not log when DEBUG env is set and level is debug", () => {
      process.env.DEBUG = "true";
      logger("debug", "This should not appear");

      expect(consoleLogSpy).not.toHaveBeenCalled();
    });

    it("should log when DEBUG is not set with debug level", () => {
      logger("debug", "Debug message");

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe("with multiple arguments", () => {
    it("should pass multiple arguments to console.log", () => {
      logger("info", "Hello", { foo: "bar" }, [1, 2, 3]);

      expect(consoleLogSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe("sensitive value masking", () => {
    const sensitiveValue = `access_token_abc123xyz`;

    afterEach(() => {
      resetEnvHandlers();
    });

    it("should show original value in dev mode with DEBUG enabled", () => {
      process.env.BUN_ENV = "development";
      process.env.DEBUG = "true";
      logger("info", "Token:", sensitive(sensitiveValue));

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        "Token:",
        expect.stringContaining(sensitiveValue),
      );
    });

    it("should mask value in dev mode without DEBUG", () => {
      process.env.BUN_ENV = "development";
      delete process.env.DEBUG;
      logger("info", "Token:", sensitive(sensitiveValue));

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        "Token:",
        expect.stringContaining("****"),
      );
    });

    it("should mask value in production mode with DEBUG", () => {
      process.env.BUN_ENV = "production";
      process.env.DEBUG = "true";
      logger("info", "Token:", sensitive(sensitiveValue));

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        "Token:",
        expect.stringContaining("****"),
      );
    });

    it("should not affect regular strings", () => {
      process.env.BUN_ENV = "development";
      process.env.DEBUG = "true";
      logger("info", "Regular message without sensitive data");

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        expect.stringContaining("Regular message without sensitive data"),
      );
    });
  });
});