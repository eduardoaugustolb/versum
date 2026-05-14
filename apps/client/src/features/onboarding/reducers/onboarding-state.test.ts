import { describe, it, expect } from "vitest";
import { onboardingReducer } from "./onboarding-state";
import type { OnboardingReducerState } from "./onboarding-state";

function createState(overrides: Partial<OnboardingReducerState> = {}): OnboardingReducerState {
  return {
    currentIndex: 0,
    direction: 1,
    collectedValues: {},
    isComplete: false,
    error: null,
    ...overrides,
  };
}

describe("onboardingReducer", () => {
  describe("NEXT", () => {
    it("avança currentIndex em 1", () => {
      const state = createState({ currentIndex: 0 });
      const next = onboardingReducer(state, { type: "NEXT" });
      expect(next.currentIndex).toBe(1);
      expect(next.direction).toBe(1);
    });

    it("não ultrapassa o último índice", () => {
      const state = createState({ currentIndex: 4 });
      const next = onboardingReducer(state, { type: "NEXT" });
      expect(next.currentIndex).toBe(4);
    });

    it("acumula patch nos collectedValues", () => {
      const state = createState({ collectedValues: { name: "João" } });
      const next = onboardingReducer(state, { type: "NEXT", patch: { username: "joao" } });
      expect(next.collectedValues).toEqual({ name: "João", username: "joao" });
    });

    it("mantém collectedValues se patch for undefined", () => {
      const state = createState({ collectedValues: { name: "João" } });
      const next = onboardingReducer(state, { type: "NEXT" });
      expect(next.collectedValues).toEqual({ name: "João" });
    });

    it("limpa error no NEXT", () => {
      const state = createState({ error: "Algo deu errado" });
      const next = onboardingReducer(state, { type: "NEXT" });
      expect(next.error).toBeNull();
    });
  });

  describe("BACK", () => {
    it("volta currentIndex em 1", () => {
      const state = createState({ currentIndex: 2 });
      const back = onboardingReducer(state, { type: "BACK" });
      expect(back.currentIndex).toBe(1);
      expect(back.direction).toBe(-1);
    });

    it("não vai abaixo de 0", () => {
      const state = createState({ currentIndex: 0 });
      const back = onboardingReducer(state, { type: "BACK" });
      expect(back.currentIndex).toBe(0);
    });
  });

  describe("COMPLETE", () => {
    it("marca isComplete como true", () => {
      const state = createState();
      const completed = onboardingReducer(state, { type: "COMPLETE" });
      expect(completed.isComplete).toBe(true);
    });
  });

  describe("SET_ERROR", () => {
    it("define a mensagem de erro", () => {
      const state = createState();
      const errored = onboardingReducer(state, { type: "SET_ERROR", error: "Falhou" });
      expect(errored.error).toBe("Falhou");
    });
  });

  describe("CLEAR_ERROR", () => {
    it("limpa o erro", () => {
      const state = createState({ error: "Falhou" });
      const cleared = onboardingReducer(state, { type: "CLEAR_ERROR" });
      expect(cleared.error).toBeNull();
    });
  });
});
