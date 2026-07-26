import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@nexus/zod", "@nexus/openapi"],
};

export default nextConfig;
