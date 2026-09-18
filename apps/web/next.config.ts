import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  allowedDevOrigins: ["192.168.1.79", "painelcomal.duckdns.org"],
  // basePath: "/winthor-ia",
  // assetPrefix: "/winthor-ia",
  output: "standalone",

  basePath: "/dev-test",
  assetPrefix: "/dev-test",
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Frame-Options", value: "DENY" }, // Bloqueia iframe
          { key: "X-Content-Type-Options", value: "nosniff" }, // Impede MIME-sniffing
          { key: "Referrer-Policy", value: "origin-when-cross-origin" },
          {
            key: "Permissions-Policy",
            value: "camera=(), microphone=(), geolocation=()",
          },
        ],
      },
    ];
  },
};

export default nextConfig;
