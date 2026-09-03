import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  allowedDevOrigins: ["192.168.1.79", "painelcomal.duckdns.org"],
  basePath: "/winthor-ia",
  output: "standalone",
  assetPrefix: "/winthor-ia",
};

export default nextConfig;
