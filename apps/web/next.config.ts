import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  allowedDevOrigins: ["192.168.1.79", "painelcomal.duckdns.org"],
  basePath: "/winthor-ia",
  assetPrefix: "/winthor-ia",
  output: "standalone",

  // basePath: "/dev-test",
  // assetPrefix: "/dev-test",
};

export default nextConfig;
