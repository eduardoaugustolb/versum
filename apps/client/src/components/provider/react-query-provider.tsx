"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type ReactNode, useEffect, useState } from "react";

export function ReactQueryProvider({ children }: { children: ReactNode }) {
  useEffect(() => {
    import("@kubb/plugin-client/clients/fetch").then(({ setConfig }) => {
      setConfig({
        baseURL: process.env.NEXT_PUBLIC_API_URL,
        credentials: "include",
      });
    });
  }, []);
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
