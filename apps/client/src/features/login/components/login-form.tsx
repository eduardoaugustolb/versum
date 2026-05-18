"use client";

import { useRef } from "react";
import { SignInIcon } from "@phosphor-icons/react";
import { StepTransition } from "@/components/shared/step-transition";
import type { StepTransitionHandle } from "@/components/shared/step-transition";
import { ActionButton } from "@/components/shared/action-button";
import { FieldError } from "@/components/shared/field-error";
import { useLoginForm } from "../hooks/use-login-form";
import { useGSAP } from "@gsap/react";
import { SplitText, gsap } from "gsap/src/all";

export function LoginForm() {
  const transitionRef = useRef<StepTransitionHandle>(null);
  const wrapperRef = useRef<HTMLDivElement | null>(null);
  const titleRef = useRef<HTMLParagraphElement | null>(null);
  const subtitleRef = useRef<HTMLParagraphElement | null>(null);

  const form = useLoginForm();

  useGSAP(() => {
    if (titleRef.current && subtitleRef.current) {
      const splitedTitle = new SplitText(titleRef.current, { type: "words" });
      const splitedSubtitle = new SplitText(subtitleRef.current, {
        type: "words",
      });
      gsap.set(splitedTitle.words, {
        opacity: 0,
        transform: "translateY(20px)",
        filter: "blur(10px)",
      });
      gsap.set(splitedSubtitle.words, {
        opacity: 0,
        transform: "translateY(20px)",
        filter: "blur(10px)",
      });

      wrapperRef.current?.classList.remove("invisible");

      gsap.to(splitedTitle.words, {
        opacity: 1,
        transform: "translateY(0)",
        filter: "blur(0px)",
        stagger: 0.1,
        duration: 1,
        ease: "power4.out",
      });
      gsap.to(splitedSubtitle.words, {
        opacity: 1,
        transform: "translateY(0)",
        filter: "blur(0px)",
        stagger: 0.1,
        duration: 0.6,
        ease: "power4.out",
      });

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
          <p
            ref={titleRef}
            className="tracking-tight text-foreground overflow-hidden font-instrument-serif text-5xl"
          >
            Bem-vindo de volta
          </p>
          <p
            className="text-xl text-foreground/50 font-instrument-sans"
            ref={subtitleRef}
          >
            Enviaremos um link mágico para seu email
          </p>
        </div>

        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            void form.handleSubmit();
          }}
          noValidate
        >
          <form.Field name="email">
            {(field) => (
              <div>
                <input
                  id={field.name}
                  type="email"
                  placeholder="seu@email.com"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  onBlur={field.handleBlur}
                  autoComplete="email"
                  autoFocus
                  className="w-full rounded-xl border border-neutral-200 bg-white px-4 py-3 text-sm text-neutral-900 placeholder-neutral-400 outline-none ring-0 transition focus:border-neutral-900 focus:ring-2 focus:ring-neutral-900/10 dark:border-neutral-700 dark:bg-neutral-800 dark:text-white dark:placeholder-neutral-500 dark:focus:border-white dark:focus:ring-white/10"
                />
                <FieldError field={field} />
              </div>
            )}
          </form.Field>

          <ActionButton
            label={
              form.state.isSubmitting ? "Enviando..." : "Enviar magic link"
            }
            type="submit"
            disabled={form.state.isSubmitting}
          />
        </form>

        <p className="text-xs text-center text-muted-foreground">
          Ao continuar, você concorda com nossos termos de uso e política de
          privacidade
        </p>
      </div>
    </StepTransition>
  );
}
