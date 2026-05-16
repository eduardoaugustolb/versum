"use client";

import { useRef } from "react";
import type { OnboardingValues, OutStep } from "../../types";
import { onboardingFormSchema } from "../../types";
import type { StepTransitionHandle } from "../step-transition";
import { StepTransition } from "../step-transition";
import { ActionButton } from "../ui/action-button";
import { useGSAP } from "@gsap/react";
import { SplitText, gsap } from "gsap/src/all";
import { useRouter } from "next/navigation";

interface StepAnimationProps {
  direction: 1 | -1;
  onExitDone: () => void;
}

interface OutStepViewProps extends StepAnimationProps {
  step: OutStep;
  values: Partial<OnboardingValues>;
  onComplete: () => void;
}

export function OutStepView({
  step,
  values,
  direction,
  onExitDone,
  onComplete,
}: OutStepViewProps) {
  const transitionRef = useRef<StepTransitionHandle>(null);
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const parsed = onboardingFormSchema.safeParse(values);
  const label = parsed.success ? step.getLabel(parsed.data) : "Seja bem-vindo!";
  const labelRef = useRef<HTMLParagraphElement | null>(null);
  const subtitleRef = useRef<HTMLParagraphElement | null>(null);
  const router = useRouter();

  function handleComplete() {
    transitionRef.current?.triggerExit(direction, () => {
      onExitDone();
      onComplete();
      router.push("/");
    });
  }

  useGSAP(() => {
    if (labelRef.current && subtitleRef.current) {
      const splitedLabel = new SplitText(labelRef.current, { type: "words" });
      const splitedSubtitle = new SplitText(subtitleRef.current, { type: "words" });
      gsap.set(splitedLabel.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });
      gsap.set(splitedSubtitle.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });

      wrapperRef.current?.classList.remove("invisible");

      gsap.to(splitedLabel.words, { opacity: 1, transform: "translateY(0)", filter: "blur(0px)", stagger: 0.1, duration: 1, ease: "power4.out" });
      gsap.to(splitedSubtitle.words, { opacity: 1, transform: "translateY(0)", filter: "blur(0px)", stagger: 0.1, duration: 0.6, ease: "power4.out" });

      return () => {
        if (splitedLabel.words) {
          gsap.killTweensOf(splitedLabel.words);
        }
        if (splitedSubtitle.words) {
          gsap.killTweensOf(splitedSubtitle.words);
        }
      }
    }
  }, [])

  return (
    <StepTransition ref={transitionRef}>
      <div className="flex h-full flex-col items-center justify-center px-8 text-center">
        <div ref={wrapperRef} className="invisible flex flex-col items-center gap-2">
          <span className="text-5xl">🎉</span>
          <p ref={labelRef} className="tracking-tight text-foreground overflow-hidden font-instrument-serif text-5xl">
            {label}
          </p>
          <p className="text-xl text-foreground/50 font-instrument-sans" ref={subtitleRef}>
            Tudo certo! Sua conta está pronta.
          </p>
        </div>
        <ActionButton label="Entrar" onClick={handleComplete} />
      </div>
    </StepTransition>
  );
};
