"use client";

import { forwardRef, useImperativeHandle, useRef } from "react";
import type { ReactNode } from "react";
import { useGSAP } from "@gsap/react";
import { gsap } from "gsap";

export interface StepTransitionHandle {
  triggerExit: (direction: number, onDone: () => void) => void;
}

interface StepTransitionProps {
  children: ReactNode;
}

const EXIT_DURATION = 0.18;
const ENTER_DURATION = 0.26;

export const StepTransition = forwardRef<
  StepTransitionHandle,
  StepTransitionProps
>(function StepTransition({ children }, ref) {
  const elRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    const el = elRef.current;
    if (!el) return;

    gsap.set(el, { opacity: 0 });
    gsap.to(el, {
      opacity: 1,
      duration: ENTER_DURATION,
      clearProps: "opacity",
    });
  }, []);

  useImperativeHandle(
    ref,
    () => ({
      triggerExit: (_direction: number, onDone: () => void) => {
        const el = elRef.current;
        if (!el) {
          onDone();
          return;
        }

        gsap.to(el, {
          opacity: 0,
          duration: EXIT_DURATION,
          onComplete: onDone,
        });
      },
    }),
    [],
  );

  return (
    <div ref={elRef} style={{ willChange: "opacity" }}>
      {children}
    </div>
  );
});
