"use client";

import { useRef } from "react";
import type { StepTransitionHandle } from "../step-transition";
import { StepTransition } from "../step-transition";
import { ActionButton } from "../ui/action-button";
import { useGSAP } from "@gsap/react";
import { SplitText, gsap } from "gsap/src/all";

import type { StepDirection } from "../../types";

interface ErrorStepViewProps {
  error: string;
  direction: StepDirection;
  onRetry: () => void;
}

export function ErrorStepView({ error, direction, onRetry }: ErrorStepViewProps) {
  const transitionRef = useRef<StepTransitionHandle>(null);
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const titleRef = useRef<HTMLParagraphElement | null>(null);
  const subtitleRef = useRef<HTMLParagraphElement | null>(null);

  function handleRetry() {
    transitionRef.current?.triggerExit(direction, onRetry);
  }

  useGSAP(() => {
    if (titleRef.current && subtitleRef.current) {
      const splitedTitle = new SplitText(titleRef.current, { type: "words" });
      const splitedSubtitle = new SplitText(subtitleRef.current, { type: "words" });
      gsap.set(splitedTitle.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });
      gsap.set(splitedSubtitle.words, { opacity: 0, transform: "translateY(20px)", filter: "blur(10px)" });

      wrapperRef.current?.classList.remove("invisible");

      gsap.to(splitedTitle.words, { opacity: 1, transform: "translateY(0)", filter: "blur(0px)", stagger: 0.1, duration: 1, ease: "power4.out" });
      gsap.to(splitedSubtitle.words, { opacity: 1, transform: "translateY(0)", filter: "blur(0px)", stagger: 0.1, duration: 0.6, ease: "power4.out" });

      return () => {
        if (splitedTitle.words) {
          gsap.killTweensOf(splitedTitle.words);
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
          <span className="text-5xl">❌</span>
          <p ref={titleRef} className="tracking-tight text-foreground overflow-hidden font-instrument-serif text-5xl">
            Algo deu errado
          </p>
          <p className="text-xl text-foreground/50 font-instrument-sans" ref={subtitleRef}>
            {error}
          </p>
        </div>
        <ActionButton label="Tentar novamente" onClick={handleRetry} />
      </div>
    </StepTransition>
  );
}
