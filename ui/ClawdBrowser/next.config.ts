import type { NextConfig } from "next";
import path from "node:path";

const optionalNodeAlias = path.resolve(
  process.cwd(),
  "src/lib/solana/noop-module.ts",
);

const nextConfig: NextConfig = {
  outputFileTracingRoot: process.cwd(),
  serverExternalPackages: ["@solana/web3.js", "@solana/spl-token"],
  webpack: (config) => {
    config.resolve = config.resolve ?? {};
    config.resolve.alias = {
      ...(config.resolve.alias ?? {}),
      "pino-pretty": optionalNodeAlias,
      encoding: optionalNodeAlias,
    };
    return config;
  },
  turbopack: {
    // Tell Turbopack not to try to bundle Node-only builtins that Solana
    // wallet libs reference via optional requires.
    resolveAlias: {
      "pino-pretty": "@/lib/solana/noop-module",
      encoding: "@/lib/solana/noop-module",
    },
  },
};

export default nextConfig;
