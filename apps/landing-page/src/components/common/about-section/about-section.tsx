"use client";

import { useRef } from "react";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { ScrollTrigger, SplitText } from "gsap/all";
import { useChangeNavColors } from "../navbar";

gsap.registerPlugin(ScrollTrigger, SplitText);

export function AboutSection() {
  const wrapperRef = useRef<HTMLDivElement>(null);
  const sectionRef = useRef<HTMLDivElement>(null);
  const textRef = useRef<HTMLHeadingElement>(null);
  const { changeColor } = useChangeNavColors();

  useGSAP(
    () => {
      if (!wrapperRef.current || !sectionRef.current || !textRef.current)
        return;

      const section = sectionRef.current;
      const text = textRef.current;

      const split = new SplitText(text, { type: "chars" });
      const chars = split.chars;

      gsap.set(chars, {
        y: 60,
        opacity: 0,
        willChange: "transform, opacity",
      });

      const scrollDistance =
        section.scrollWidth -
        window.innerWidth +
        (innerWidth >= 1024 ? 128 : 30);

      if (scrollDistance <= 0) {
        gsap.to(chars, {
          y: 0,
          opacity: 1,
          stagger: 0.01,
          duration: 0.6,
          ease: "power2.out",
        });
        return () => split.revert();
      }

      // ✅ CORREÇÃO AQUI - arrow function simples
      const changeNavColor = () => changeColor("PRIMARY_NAVBAR_COLOR");

      gsap
        .timeline({
          scrollTrigger: {
            trigger: wrapperRef.current,
            start: "top top",
            end: `+=${scrollDistance}`,
            scrub: 1,
            pin: true,
            anticipatePin: 1,
            invalidateOnRefresh: true,
          },
        })
        .to(section, {
          x: -scrollDistance,
          ease: "none",
        });

      const revealTl = gsap.timeline({
        scrollTrigger: {
          trigger: wrapperRef.current,
          start: `top ${innerWidth >= 769 ? "top" : "5%"}`,
          end: `+=${scrollDistance} ${innerWidth < 769 ? "-15%" : ""}`,
          scrub: 1,
          onEnter: changeNavColor,
          onEnterBack: changeNavColor,
          onLeaveBack: () => changeColor("SECONDARY_NAVBAR_COLOR"),
          // markers: true,
        },
      });

      const totalChars = chars.length;

      chars.forEach((char, index) => {
        const progress = index / totalChars;
        revealTl.to(
          char,
          {
            y: 0,
            opacity: 1,
            duration: 0.1,
            ease: "power2.out",
          },
          progress * 0.8,
        );
      });

      return () => {
        split.revert();
      };
    },
    { scope: wrapperRef, dependencies: [] },
  );

  return (
    <div
      id="about-wrapper"
      ref={wrapperRef}
      className="w-full h-screen bg-[#151513] overflow-hidden relative"
    >
      <div
        id="about"
        ref={sectionRef}
        className="h-screen flex items-center text-[#F4EAD8] whitespace-nowrap px-20 lg:px-32 w-max min-w-screen"
      >
        <h2
          ref={textRef}
          className="text-8xl font-instrument-serif leading-tight overflow-hidden"
        >
          Versum promove constância na leitura da Bíblia, respeitando o ritmo e
          a reflexão.
        </h2>
      </div>
    </div>
  );
}
