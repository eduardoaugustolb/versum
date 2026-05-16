import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  async rewrites() {
    if (process.env.NODE_ENV === "development") {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL;
      if (apiUrl) {
        return [
          {
            source: "/api/:path*",
            destination: `${apiUrl}/api/:path*`,
          },
        ];
      }
    }
    return [];
  },
};

export default nextConfig;
