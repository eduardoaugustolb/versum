"use client";

import HeroImg from "@/assets/img/hero-img.webp";
import { fragmentShader, vertexShader } from "@/assets/shader/home-page";
import { LoadingScreen } from "@/components/common/loading-screen/component";
import { useChangeNavColors } from "@/components/common/navbar";
import { hexToRGB } from "@/utils/hexToRGB";
import { useGSAP } from "@gsap/react";
import { ArrowElbowDownRightIcon } from "@phosphor-icons/react";
import gsap from "gsap";
import { ScrollTrigger, SplitText } from "gsap/all";
import Lenis from "lenis";
import Image from "next/image";
import Link from "next/link";
import { useRef, useState } from "react";
import * as THREE from "three";
import { AboutSection } from "../components/common/about-section/about-section";

export default function Home() {
  const { changeColor } = useChangeNavColors();
  const [isLoading, setIsLoading] = useState(true);

  // Refs
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const heroRef = useRef<HTMLElement | null>(null);
  const heroParagraphRef = useRef<HTMLParagraphElement | null>(null);

  const heroHeaderTitleRef = useRef<HTMLHeadingElement | null>(null);
  const heroHeaderParagraphRef = useRef<HTMLParagraphElement | null>(null);
  const heroLinksRef = useRef<HTMLDivElement | null>(null);

  gsap.registerPlugin(useGSAP, ScrollTrigger, SplitText);

  useGSAP(() => {
    if (!heroHeaderTitleRef.current || !heroHeaderParagraphRef.current) return;
    const titleSplited = new SplitText(heroHeaderTitleRef.current, {
      type: "chars",
    });

    const paragraphSplited = new SplitText(heroHeaderParagraphRef.current, {
      type: "words",
    });

    gsap.set(titleSplited.chars, {
      y: 100,
    });

    gsap.set(paragraphSplited.words, {
      y: 50,
    });

    if (heroLinksRef.current) {
      gsap.set(heroLinksRef.current, {
        opacity: 0,
        y: 30,
      });
    }

    if (isLoading) return;

    gsap.to(titleSplited.chars, {
      y: 0,
      stagger: 0.05,
      duration: 1,
      ease: "power4.out",
      delay: 0.5,
    });

    gsap.to(paragraphSplited.words, {
      y: 0,
      stagger: 0.05,
      duration: 1,
      ease: "power4.out",
      delay: 0.5,
    });

    gsap.to(heroLinksRef.current, {
      opacity: 1,
      y: 0,
      duration: 1,
      ease: "power4.out",
      delay: 1,
      stagger: 0.3,
    });
  }, [isLoading]);

  useGSAP(() => {
    if (!canvasRef.current || !heroRef.current) return;

    const lenis = new Lenis({
      syncTouch: true,
    });

    function raf(time: number) {
      lenis.raf(time);
      ScrollTrigger.update();
      requestAnimationFrame(raf);
    }

    requestAnimationFrame(raf);

    const CONFIG = {
      color: "#F4EAD8",
      spread: 0.5,
      speed: 2,
    };

    const scene = new THREE.Scene();
    const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0, 1);
    const renderer = new THREE.WebGLRenderer({
      canvas: canvasRef.current,
      alpha: true,
      antialias: true,
    });

    const rgb = hexToRGB(CONFIG.color);

    // CRIAR O MATERIAL PRIMEIRO
    const geometry = new THREE.PlaneGeometry(2, 2);
    const material = new THREE.ShaderMaterial({
      vertexShader,
      fragmentShader,
      uniforms: {
        uProgress: { value: 0 },
        uResolution: {
          value: new THREE.Vector2(
            heroRef.current.offsetWidth,
            heroRef.current.offsetHeight,
          ),
        },
        uColor: { value: new THREE.Vector3(rgb.r, rgb.g, rgb.b) },
        uSpread: { value: CONFIG.spread },
      },
      transparent: true,
    });

    const mesh = new THREE.Mesh(geometry, material);
    scene.add(mesh);

    // AGORA DEFINIR A FUNÇÃO RESIZE (que usa material)
    function resize() {
      const width = heroRef.current?.offsetWidth;
      const height = heroRef.current?.offsetHeight;

      if (!width || !height) return;

      renderer.setSize(width, height);
      renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

      material.uniforms.uResolution.value.set(width, height);
    }

    // CHAMAR RESIZE DEPOIS
    resize();
    window.addEventListener("resize", resize);

    function animate() {
      renderer.render(scene, camera);
      requestAnimationFrame(animate);
    }

    animate();

    lenis.on("scroll", ({ scroll }) => {
      const heroHeight = heroRef.current?.offsetHeight || 0;
      const windowHeight = window.innerHeight;

      const maxScroll = heroHeight - windowHeight;
      const progress = Math.max(
        0,
        Math.min((scroll / maxScroll) * CONFIG.speed, 1.1),
      );

      material.uniforms.uProgress.value = progress;
    });

    const split = new SplitText(heroParagraphRef.current, {
      type: "words",
    });

    const words = split.words;

    gsap.set(words, { opacity: 0 });

    if (isLoading) return;

    const heroTextScrollTrigger = ScrollTrigger.create({
      trigger: "#hero-content",
      start: "top 25%",
      end: "bottom 100%",
      onUpdate: (self) => {
        const progress = self.progress;

        if (progress >= 0.2) {
          changeColor("SECONDARY_NAVBAR_COLOR");
        }
        console.log(progress);
        if (progress < 0.2) {
          changeColor("PRIMARY_NAVBAR_COLOR");
        }

        const totalWords = words.length;

        words.forEach((word, index) => {
          const wordProgress = index / totalWords;
          const nextWordProgress = (index + 1) / totalWords;

          let opacity = 0;

          if (progress >= nextWordProgress) opacity = 1;
          if (progress >= wordProgress) {
            const fadeProgress =
              (progress - wordProgress) / (nextWordProgress - wordProgress);
            opacity = fadeProgress;
          }

          if (index == 0) {
            const current = Math.max(
              0,
              parseFloat(
                gsap.getProperty("#hero-header", "opacity").toString(),
              ),
            );
            const next = Math.max(0, 1 - opacity);

            if (current !== next) {
              gsap.to("#hero-header", {
                opacity: next,
                overwrite: "auto",
              });
            }
          }

          gsap.to(word, {
            opacity,
          });
        });
      },
    });

    // Cleanup
    return () => {
      window.removeEventListener("resize", resize);
      lenis.destroy();
      geometry.dispose();
      material.dispose();
      renderer.dispose();
      heroTextScrollTrigger.disable();
    };
  }, [isLoading]);

  return (
    <main className="w-screen overflow-x-hidden">
      <LoadingScreen isLoading={isLoading} setIsLoading={setIsLoading} />
      <>
        <section
          id="hero"
          ref={heroRef}
          className="relative w-full h-[175svh] text-[#F4EAD8] overflow-hidden"
        >
          <div id="hero-img">
            <Image
              src={HeroImg}
              alt="Jesus"
              sizes="auto"
              priority
              className="w-full h-full object-cover absolute brightness-70"
            />
          </div>
          <header
            id="hero-header"
            className="absolute w-full h-svh flex flex-col
        justify-center items-center gap-2 text-center"
          >
            <h1
              ref={heroHeaderTitleRef}
              className="uppercase font-medium leading-[0.9] font-instrument-serif overflow-hidden"
              style={{ fontSize: "clamp(4rem, 7.5vw, 10rem)" }}
            >
              Versum
            </h1>
            <p
              ref={heroHeaderParagraphRef}
              className="font-instrument-sans text-lg lg:text-xl font-normal w-75/100 overflow-hidden text-shadow-2xs"
            >
              Leia a Bíblia do seu jeito, sem perder o caminho.
            </p>
            <div
              ref={heroLinksRef}
              className="flex flex-col md:flex-row gap-2 md:gap-10 overflow-hidden"
            >
              <Link
                href="/login"
                className="text-amber-50 text-lg flex items-end
              justify-center text-center gap-4 font-instrument-sans
              z-99 group "
              >
                <ArrowElbowDownRightIcon weight="duotone" size={24} />
                <span className="relative hover:cursor-pointer">
                  Começar agora
                  <div
                    className="h-[0.4] w-0 group-hover:w-full bg-amber-50 -mt-0.5
                transition-all duration-300 ease-out"
                  />
                </span>
              </Link>
              <Link
                href="/login"
                className="text-amber-50 text-lg flex items-end
              justify-center text-center gap-4 font-instrument-sans
              z-99 group"
              >
                <ArrowElbowDownRightIcon
                  weight="duotone"
                  size={24}
                  className="mt-2"
                />
                <span className="relative hover:cursor-pointer">
                  Apoie o projeto
                  <div
                    className="h-[0.4] w-0 group-hover:w-full bg-amber-50 -mt-0.5
                transition-all duration-300 ease-out"
                  />
                </span>
              </Link>
            </div>
          </header>
          <canvas
            ref={canvasRef}
            id="hero-canvas"
            className="absolute bottom-0 w-full h-full pointer-events-none"
          ></canvas>
          <div
            id="hero-content"
            className="absolute bottom-0 w-full h-[125svh] flex flex-col justify-center
        items-center text-center gap-2 text-[#302f2c] pointer-events-none"
          >
            <p
              className="uppercase font-medium leading-[0.98] font-instrument-serif
            w-[100%-4rem] lg:w-75/100"
              ref={heroParagraphRef}
              style={{ fontSize: "clamp(2.5rem, 4.5vw, 5rem)" }}
            >
              "Ainda que eu andasse pelo vale da sombra da morte, não temeria
              mal algum, porque tu estás comigo."
            </p>
            <p>— Salmos 23:4</p>
          </div>
        </section>
        <AboutSection />
      </>
    </main>
  );
}
