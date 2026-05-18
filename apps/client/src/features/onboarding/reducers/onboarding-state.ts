import { TOTAL_STEPS } from "../constants";
import type { OnboardingValues, StepDirection } from "../types";

export interface OnboardingReducerState {
  currentIndex: number;
  direction: StepDirection;
  collectedValues: Partial<OnboardingValues>;
  isComplete: boolean;
  error: string | null;
}

type Action =
  | { type: "NEXT"; patch?: Partial<OnboardingValues> }
  | { type: "BACK" }
  | { type: "COMPLETE" }
  | { type: "SET_ERROR"; error: string }
  | { type: "CLEAR_ERROR" };

export function onboardingReducer(
  state: OnboardingReducerState,
  action: Action,
): OnboardingReducerState {
  switch (action.type) {
    case "NEXT": {
      const next = Math.min(state.currentIndex + 1, TOTAL_STEPS - 1);
      return {
        ...state,
        currentIndex: next,
        direction: 1,
        error: null,
        collectedValues: action.patch
          ? { ...state.collectedValues, ...action.patch }
          : state.collectedValues,
      };
    }
    case "BACK": {
      const prev = Math.max(state.currentIndex - 1, 0);
      return {
        ...state,
        currentIndex: prev,
        direction: -1,
      };
    }
    case "COMPLETE":
      return { ...state, isComplete: true };
    case "SET_ERROR":
      return { ...state, error: action.error };
    case "CLEAR_ERROR":
      return { ...state, error: null };
    default:
      return state;
  }
}
