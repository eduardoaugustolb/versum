"use client";

import { useState, useRef, Dispatch, SetStateAction } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { CounterDigit } from "./counter-digit";
import "./style.css";

export function LoadingScreen({
  setIsLoading,
  isLoading,
}: {
  setIsLoading: Dispatch<SetStateAction<boolean>>;
  isLoading: boolean;
}) {
  const container = useRef<HTMLDivElement>(null);

  useGSAP(
    () => {
      // Atualiza o contador para um valor específico
      const updateCounter = (value: number) => {
        const hundreds = Math.floor(value / 100);
        const tens = Math.floor((value % 100) / 10);
        const ones = value % 10;

        const duration = Math.random() * 0.4 + 0.2;
        const ease = "power2.out";

        gsap.to("#digit-1", { y: `${-hundreds * 100}px`, duration, ease });
        gsap.to("#digit-2", { y: `${-tens * 100}px`, duration, ease });
        gsap.to("#digit-3", { y: `${-ones * 100}px`, duration, ease });
      };

      // Gera incrementos aleatórios que somam 100
      const generateSteps = () => {
        const steps: number[] = [];
        let total = 0;
        const numSteps = Math.floor(Math.random() * 8) + 5;

        for (let i = 0; i < numSteps - 1; i++) {
          const step = Math.floor(Math.random() * 18) + 3;
          steps.push(step);
          total += step;
        }

        const last = 100 - total;
        steps.push(Math.max(last, 0));
        return steps;
      };

      // Executa a animação
      const animate = async () => {
        const steps = generateSteps();
        let currentValue = 0;

        for (const step of steps) {
          currentValue = Math.min(currentValue + step, 100);
          updateCounter(currentValue);

          // Aguarda duração da animação + pausa
          const animDuration = (Math.random() * 0.4 + 0.2) * 1000;
          const pauseDuration = (Math.random() * 0.3 + 0.1) * 1000;
          await new Promise((resolve) =>
            setTimeout(resolve, animDuration + pauseDuration)
          );
        }

        // Pausa no 100% e sai da tela
        await new Promise((resolve) => setTimeout(resolve, 300));

        gsap.to(container.current, {
          yPercent: -100,
          duration: 0.8,
          ease: "power4.out",
          onComplete: () => setIsLoading(false),
        });
      };

      animate();
    },
    { scope: container }
  );

  if (!isLoading) return null;

  return (
    <div ref={container} id="loading-screen">
      <p>Loading</p>
      <div id="counter">
        <CounterDigit id="digit-1" />
        <CounterDigit id="digit-2" />
        <CounterDigit id="digit-3" />
        <div className="num">%</div>
      </div>
    </div>
  );
}
