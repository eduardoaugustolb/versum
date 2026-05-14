import { useCallback, useReducer } from "react";
import { onboardingReducer, type OnboardingReducerState } from "../reducers/onboarding-state";
import type { OnboardingValues } from "../types";

const INITIAL_STATE: OnboardingReducerState = {
  currentIndex: 0,
  direction: 1,
  collectedValues: {},
  isComplete: false,
  error: null,
};

export function useOnboardingState() {
  const [state, dispatch] = useReducer(onboardingReducer, INITIAL_STATE);

  const goNext = useCallback((patch?: Partial<OnboardingValues>) => {
    dispatch({ type: "NEXT", patch });
  }, []);

  const goBack = useCallback(() => {
    dispatch({ type: "BACK" });
  }, []);

  const complete = useCallback(() => {
    dispatch({ type: "COMPLETE" });
  }, []);

  const setError = useCallback((error: string) => {
    dispatch({ type: "SET_ERROR", error });
  }, []);

  const clearError = useCallback(() => {
    dispatch({ type: "CLEAR_ERROR" });
  }, []);

  return { state, goNext, goBack, complete, setError, clearError } as const;
}
