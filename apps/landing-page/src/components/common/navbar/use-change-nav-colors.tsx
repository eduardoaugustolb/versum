"use client";
import { useEffect } from "react";

const EVENT_NAME = "change-nav-colors";
let lastColor = "#fffbeb";

export const defaultColors = {
  PRIMARY_NAVBAR_COLOR: "#fffbeb",
  SECONDARY_NAVBAR_COLOR: "#302f2c",
} as const;

type DeafultColorType = keyof typeof defaultColors;

export function useChangeNavColors(onListen?: (color: string) => void) {
  const changeColor = (color: string) => {
    const usingDefualtColor = Object.keys(defaultColors).includes(color);
    const selectedColor = usingDefualtColor
      ? defaultColors[color as DeafultColorType]
      : color;

    // ✅ Removi a checagem que impedia a mudança
    // Agora sempre emite o evento
    const event = new CustomEvent(EVENT_NAME, {
      detail: { color: selectedColor },
    });

    console.log(`Emitindo: de ${lastColor} para ${selectedColor}!`);
    window.dispatchEvent(event);
  };

  useEffect(() => {
    if (!onListen) return;

    const handler = (e: Event) => {
      const evt = e as CustomEvent<{ color: string }>;
      const next = evt.detail.color;

      // ✅ Só atualiza se for diferente
      if (next !== lastColor) {
        console.log(`Mudando de ${lastColor} para ${next}!`);
        lastColor = next;
        onListen(next);
      }
    };

    // ✅ Chama onListen com a cor atual ao montar
    onListen(lastColor);

    window.addEventListener(EVENT_NAME, handler);
    return () => window.removeEventListener(EVENT_NAME, handler);
  }, [onListen]);

  return { changeColor };
}
