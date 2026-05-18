"use client";

import { useRef } from "react";
import { cn } from "@/lib/utils";
import type { InStep, StepDirection } from "../../types";
import type { StepTransitionHandle } from "../step-transition";
import { StepTransition } from "../step-transition";
import { ActionButton } from "../ui/action-button";
import { useGSAP } from "@gsap/react";
import { SplitText, gsap } from "gsap/src/all";

interface StepAnimationProps {
  direction: StepDirection;
  onExitDone: () => void;
}

interface InStepViewProps extends StepAnimationProps {
  step: InStep;
}

export function InStepView({
  step,
  direction,
  onExitDone,
}: InStepViewProps) {
  const transitionRef = useRef<StepTransitionHandle>(null);
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const labelRef = useRef<HTMLParagraphElement | null>(null);
  const subtitleRef = useRef<HTMLParagraphElement | null>(null);

  function handleNext() {
    transitionRef.current?.triggerExit(direction, onExitDone);
  }

  useGSAP(() => {
    if (labelRef.current && subtitleRef.current) {
      const splitedLabel = new SplitText(labelRef.current, { type: "words" });
      const splitedSubtitle = new SplitText(subtitleRef.current, { type: "words" });
      gsap.set(splitedLabel.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });
      gsap.set(splitedSubtitle.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });

      wrapperRef.current?.classList.remove("invisible");

      //animation
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
      <div className="flex flex-col items-center">
        <div ref={wrapperRef} className="invisible flex flex-col items-center gap-2">
          <p ref={labelRef} className={cn("tracking-tight text-foreground overflow-hidden font-instrument-serif text-5xl")}>
              {step.label}
            </p>
            <p className="text-xl text-foreground/50 font-instrument-sans" ref={subtitleRef}>
              Vamos configurar sua conta em alguns passos rápidos.
            </p>
        </div>
        <ActionButton label="Começar" onClick={handleNext} />
      </div>
    </StepTransition>
  );
};
