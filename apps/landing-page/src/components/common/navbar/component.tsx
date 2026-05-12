"use client";

import gsap from "gsap";
import { ListIcon } from "@phosphor-icons/react";
import { useChangeNavColors } from "./use-change-nav-colors";
import { useRef, useEffect } from "react";

export function Navbar() {
  const navbarRef = useRef<HTMLElement>(null);

  useChangeNavColors((newColor) => {
    // ✅ Usa o ref em vez de getElementById
    if (!navbarRef.current) return;

    gsap.to(navbarRef.current, {
      color: newColor,
      duration: 0.3,
    });
  });

  return (
    <nav
      ref={navbarRef}
      id="landing-page-navbar"
      className="flex justify-between items-center w-screen fixed top-0 left-0 z-99 px-8 py-4 text-amber-50"
    >
      <p className="font-instrument-serif text-2xl">Versum</p>

      <ListIcon size={26} weight="duotone" />
    </nav>
  );
}
