"use client";

import { useMemo, useCallback } from "react";
import { useOnboardingState } from "../hooks/use-onboarding-state";
import { InStepView, ConsentStepView, FormStepView, OutStepView, ErrorStepView } from "./steps";
import { STEPS, FORM_STEPS } from "../constants";
import { onboardingFormSchema } from "../types";

interface OnboardingFlowProps {
  onComplete: (values: OnboardingValues) => void;
}

export { FORM_STEPS, type OnboardingValues };

export function OnboardingFlow({ onComplete }: OnboardingFlowProps) {
  const { state, goNext, goBack, complete, setError, clearError } = useOnboardingState();
  const { currentIndex, direction, collectedValues, error } = state;

  const currentStep = STEPS[currentIndex];

  const firstFormStepIndex = useMemo(
    () => STEPS.findIndex((s) => s.kind === "form"),
    []
  );

  const lastFormStepIndex = useMemo(() => {
    let last = -1;
    STEPS.forEach((s, i) => { if (s.kind === "form") last = i; });
    return last;
  }, []);

  const noop = useCallback(() => {}, []);

  function handleComplete() {
    complete();
    const parsed = onboardingFormSchema.safeParse(collectedValues);
    if (parsed.success) onComplete(parsed.data);
  }

  return (
    <div
      className={"relative mx-auto flex w-9/10 md:w-[720px] h-svh flex-col items-center justify-center overflow-hidden"}
    >
      <div className="flex flex-col items-center justify-center">
        {error ? (
          <ErrorStepView
            key="error"
            error={error}
            direction={direction}
            onRetry={clearError}
          />
        ) : currentStep.kind === "in" ? (
          <InStepView
            key={currentStep.id}
            step={currentStep}
            direction={direction}
            onExitDone={goNext}
          />
        ) : currentStep.kind === "consent" ? (
          <ConsentStepView
            key={currentStep.id}
            step={currentStep}
            direction={direction}
            onExitDone={goNext}
            onNext={() => goNext({})}
            onError={setError}
          />
        ) : currentStep.kind === "form" ? (
          <FormStepView
            key={currentStep.id}
            step={currentStep}
            defaultValue={
              (collectedValues[currentStep.field] as string | undefined) ?? ""
            }
            collectedValues={collectedValues}
            direction={direction}
            onExitDone={noop}
            onNext={goNext}
            onBack={goBack}
            isFirstFormStep={currentIndex === firstFormStepIndex}
            isLastFormStep={currentIndex === lastFormStepIndex}
            onError={setError}
          />
        ) : currentStep.kind === "out" ? (
          <OutStepView
            key={currentStep.id}
            step={currentStep}
            values={collectedValues}
            direction={direction}
            onExitDone={noop}
            onComplete={handleComplete}
          />
        ) : null}
      </div>

      <p className="pb-5 text-center text-xs text-neutral-400 mt-4">
        {currentIndex + 1} / {STEPS.length}
      </p>
    </div>
  );
}
