import { describe, it, expect } from "vitest";
import { onboardingFormSchema } from "./types";

describe("onboardingFormSchema", () => {
  describe("name", () => {
    it("aceita nome válido", () => {
      const result = onboardingFormSchema.shape.name.safeParse("João Silva");
      expect(result.success).toBe(true);
    });

    it("rejeita nome com menos de 2 caracteres", () => {
      const result = onboardingFormSchema.shape.name.safeParse("A");
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("Nome deve ter pelo menos 2 caracteres");
    });

    it("rejeita nome com mais de 60 caracteres", () => {
      const result = onboardingFormSchema.shape.name.safeParse("A".repeat(61));
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("Nome muito longo");
    });
  });

  describe("username", () => {
    it("aceita username válido", () => {
      const result = onboardingFormSchema.shape.username.safeParse("joao_silva");
      expect(result.success).toBe(true);
    });

    it("rejeita username com menos de 3 caracteres", () => {
      const result = onboardingFormSchema.shape.username.safeParse("ab");
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("Username deve ter pelo menos 3 caracteres");
    });

    it("rejeita username com mais de 30 caracteres", () => {
      const result = onboardingFormSchema.shape.username.safeParse("a".repeat(31));
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("Username muito longo");
    });

    it("rejeita username com caracteres especiais", () => {
      const result = onboardingFormSchema.shape.username.safeParse("joão@silva");
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("Apenas letras minúsculas, números e underscores");
    });

    it("rejeita username com letras maiúsculas", () => {
      const result = onboardingFormSchema.shape.username.safeParse("JoaoSilva");
      expect(result.success).toBe(false);
    });
  });

  describe("bio", () => {
    it("aceita bio válida", () => {
      const result = onboardingFormSchema.shape.bio.safeParse("Sou um desenvolvedor.");
      expect(result.success).toBe(true);
    });

    it("rejeita bio vazia", () => {
      const result = onboardingFormSchema.shape.bio.safeParse("");
      expect(result.success).toBe(false);
      expect(result.error?.issues[0]?.message).toBe("A bio é obrigatória");
    });
  });
});
