"use client";

import { useRef, useState } from "react";
import { cn } from "@/lib/utils";
import type { ConsentStep, StepDirection } from "../../types";
import { CONSENT_PURPOSES } from "../../types";
import type { StepTransitionHandle } from "../step-transition";
import { StepTransition } from "../step-transition";
import { ActionButton } from "../ui/action-button";
import { useGSAP } from "@gsap/react";
import { SplitText, gsap } from "gsap/src/all";
import { postApiV1Consent } from "@/dal/orval/fetch/consent-logs/consent-logs";
import { toast } from "sonner";

interface StepAnimationProps {
  direction: StepDirection;
  onExitDone: () => void;
}

interface ConsentStepViewProps extends StepAnimationProps {
  step: ConsentStep;
  onNext: () => void;
  onError: (message: string) => void;
}

export function ConsentStepView({
  step,
  direction,
  onExitDone,
  onNext,
  onError,
}: ConsentStepViewProps) {
  const transitionRef = useRef<StepTransitionHandle>(null);
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const titleRef = useRef<HTMLParagraphElement | null>(null);
  const subtitleRef = useRef<HTMLParagraphElement | null>(null);
  const [selected, setSelected] = useState<Set<string>>(
    new Set(CONSENT_PURPOSES.filter((p) => {
      const opt = step.options.find((o) => o.purpose === p);
      return opt?.required ?? false;
    })),
  );
  const [isSubmitting, setIsSubmitting] = useState(false);

  function toggle(purpose: string) {
    const opt = step.options.find((o) => o.purpose === purpose);
    if (opt?.required) return;

    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(purpose)) {
        next.delete(purpose);
      } else {
        next.add(purpose);
      }
      return next;
    });
  }

  async function handleConfirm() {
    if (selected.size === 0) {
      toast.error("Selecione pelo menos uma opção para continuar.");
      return;
    }

    const requiredAllSelected = step.options
      .filter((o) => o.required)
      .every((o) => selected.has(o.purpose));

    if (!requiredAllSelected) {
      toast.error("As opções obrigatórias devem ser aceitas.");
      return;
    }

    setIsSubmitting(true);

    try {
      await postApiV1Consent({
        consents: CONSENT_PURPOSES.map((purpose) => ({
          purpose,
          granted: selected.has(purpose),
        })),
      });

      transitionRef.current?.triggerExit(direction, () => {
        onExitDone();
        onNext();
      });
    } catch (e) {
      const message = e instanceof Error ? e.message : "Não foi possível registrar seu consentimento.";
      toast.error(message);
      onError(message);
    } finally {
      setIsSubmitting(false);
    }
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
      };
    }
  }, []);

  return (
    <StepTransition ref={transitionRef}>
      <div className="flex h-full flex-col justify-center gap-6 px-8">
        <div ref={wrapperRef} className="invisible flex flex-col gap-2">
          <p ref={titleRef} className="tracking-tight text-foreground overflow-hidden font-instrument-serif text-5xl">
            Privacidade
          </p>
          <p className="text-xl text-foreground/50 font-instrument-sans" ref={subtitleRef}>
            Sua privacidade é importante. Escolha como seus dados serão usados.
          </p>
        </div>

        <div className="flex flex-col gap-3 my-4">
          {step.options.map((option) => (
            <button
              key={option.purpose}
              type="button"
              onClick={() => toggle(option.purpose)}
              className={cn(
                "flex items-start gap-3 rounded-xl border p-4 text-left transition",
                "hover:border-neutral-400 dark:hover:border-neutral-500",
                selected.has(option.purpose)
                  ? "border-neutral-900 bg-neutral-100 dark:border-white dark:bg-neutral-800"
                  : "border-neutral-200 bg-white dark:border-neutral-700 dark:bg-neutral-900",
              )}
            >
              <div
                className={cn(
                  "mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border transition",
                  selected.has(option.purpose)
                    ? "border-neutral-900 bg-neutral-900 dark:border-white dark:bg-white"
                    : "border-neutral-300 dark:border-neutral-600",
                )}
              >
                {selected.has(option.purpose) && (
                  <svg className="h-3 w-3 text-white dark:text-neutral-900" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                )}
              </div>

              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">
                  {option.label}
                  {option.required && (
                    <span className="ml-1 text-xs text-neutral-400">(obrigatório)</span>
                  )}
                </span>
                <span className="text-xs text-neutral-500 dark:text-neutral-400 mt-0.5">
                  {option.description}
                </span>
              </div>
            </button>
          ))}
        </div>

        <ActionButton
          label="Confirmar e continuar"
          onClick={handleConfirm}
          disabled={isSubmitting}
        />
      </div>
    </StepTransition>
  );
}
