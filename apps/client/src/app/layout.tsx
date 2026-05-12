import type { Metadata } from "next";
import { IBM_Plex_Mono, IBM_Plex_Sans, IBM_Plex_Serif, Instrument_Sans, Instrument_Serif } from "next/font/google";
import { ReactQueryProvider } from "@/components/provider/react-query-provider";
import { ThemeProvider } from "@/components/provider/theme-provider";
import { Toaster } from "@/components/ui/sonner";
import { cn } from "@/lib/utils";
import "./globals.css";

const fontSans = IBM_Plex_Sans({
  subsets: ["latin"],
  variable: "--font-sans",
});

const fontSerif = IBM_Plex_Serif({
  subsets: ["latin"],
  variable: "--font-serif",
  weight: ["100", "200", "300", "400", "500", "600", "700"],
});

const fontMono = IBM_Plex_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  weight: ["100", "200", "300", "400", "500", "600", "700"],
});


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
  description: "Leitura Bíblica",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="pt-BR"
      className={cn(
        "h-full",
        "antialiased",
        fontMono.variable,
        fontSerif.variable,
        fontSans.variable,
        instrumentSans.variable,
        instrumentSerif.variable,
      )}
      suppressHydrationWarning
    >
      <body className="min-h-full flex flex-col">
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <ReactQueryProvider>
            {children}
          </ReactQueryProvider>
          <Toaster />
        </ThemeProvider>
      </body>
    </html>
  );
}
