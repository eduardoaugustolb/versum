"use client";

import { useChangeNavColors } from "@/components/common/navbar";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { SplitText } from "gsap/all";
import { useRef } from "react";

enum ShowingText {
  TEXT_1,
  TEXT_2,
}

export default function NotFound() {
  const { changeColor } = useChangeNavColors();
  gsap.registerPlugin(useGSAP, SplitText);

  const showingText = useRef<ShowingText>(ShowingText.TEXT_1);

  useGSAP(() => {
    changeColor("SECONDARY_NAVBAR_COLOR");
    const text1 = document.querySelector("#text-1");
    const text2 = document.querySelector("#text-2");

    if (!text1 || !text2) return;

    const splitedText1 = new SplitText(text1, {
      type: "chars",
    });

    const splitedText2 = new SplitText(text2, {
      type: "chars",
    });

    gsap.set(text2, {
      opacity: 1,
    });
    gsap.set(splitedText1.chars, { y: 40 });
    gsap.set(splitedText2.chars, { y: 40 });

    const animateText = () => {
      const currentElement =
        showingText.current === ShowingText.TEXT_1
          ? splitedText1.chars
          : splitedText2.chars;

      const nextElement =
        showingText.current === ShowingText.TEXT_1
          ? splitedText2.chars
          : splitedText1.chars;

      gsap.to(currentElement, {
        y: -40,
        stagger: 0.05,
        ease: "power4.in",
        duration: 0.6,
      });

      gsap.fromTo(
        nextElement,
        { y: 40 },
        {
          y: 0,
          stagger: 0.05,
          ease: "power4.out",
          duration: 0.8,
          delay: 0.35,
        }
      );

      showingText.current =
        showingText.current === ShowingText.TEXT_1
          ? ShowingText.TEXT_2
          : ShowingText.TEXT_1;
    };

    gsap.to(splitedText1.chars, {
      y: 0,
      stagger: 0.05,
      ease: "power4.out",
      duration: 0.8,
    });

    const interval = setInterval(animateText, 4000);

    return () => clearInterval(interval);
  }, []);

  return (
    <main className="bg-amber-100 w-screen h-svh flex items-center justify-center">
      <div className="relative overflow-hidden h-12 flex items-center justify-center">
        <p id="text-1" className="text-4xl font-instrument-serif">
          Not Found
        </p>
        <p
          id="text-2"
          className="text-4xl font-instrument-serif absolute"
          style={{
            opacity: 0,
          }}
        >
          404 Error
        </p>
      </div>
    </main>
  );
}
