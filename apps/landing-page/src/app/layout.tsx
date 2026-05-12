import type { Metadata } from "next";
import "./globals.css";

import { Instrument_Sans, Instrument_Serif } from "next/font/google";
import { Navbar } from "@/components/common/navbar/index";

const instrumentSans = Instrument_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-instrument-sans",
});

const instrumentSerif = Instrument_Serif({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-instrument-serif",
});

export const metadata: Metadata = {
  title: "Versum",
  description: "Leia a Bíblia de forma simples e prática.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body
        className={`${instrumentSans.variable} ${instrumentSerif.variable} antialiased`}
      >
        <Navbar />
        {children}
      </body>
    </html>
  );
}
