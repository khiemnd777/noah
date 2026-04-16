import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

function resolveFrontendServerOrigin(originValue: string) {
  let parsed: URL;

  try {
    parsed = new URL(originValue);
  } catch {
    throw new Error(`Invalid APP_FE_ORIGIN: ${originValue}`);
  }

  const port = parsed.port
    ? Number(parsed.port)
    : parsed.protocol === "https:"
      ? 443
      : 80;

  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    throw new Error(`Invalid APP_FE_ORIGIN port: ${originValue}`);
  }

  return {
    host: parsed.hostname,
    port,
  };
}

export default defineConfig(({ mode }) => {
  const repoRoot = path.resolve(__dirname, "..");
  const sharedEnv = loadEnv(mode, repoRoot, "APP_");
  const env = loadEnv(mode, __dirname, "VITE_");
  const frontendOrigin = sharedEnv.APP_FE_ORIGIN || "http://localhost:5173";
  const frontendServer = resolveFrontendServerOrigin(frontendOrigin);
  const baseAddress = env.VITE_BASE_ADDRESS || "127.0.0.1:7999";
  const httpProto = env.VITE_HTTP_PROTOCOL || "http";
  const target = `${httpProto}://${baseAddress}`;
  const apiBasePath = "/api";

  return {
    plugins: [react()],
    resolve: {
      alias: {
        "@root": path.resolve(__dirname, "src"),
        "@core": path.resolve(__dirname, "src/core"),
        "@store": path.resolve(__dirname, "src/store"),
        "@routes": path.resolve(__dirname, "src/routes"),
        "@pages": path.resolve(__dirname, "src/pages"),
        "@features": path.resolve(__dirname, "src/features"),
        "@shared": path.resolve(__dirname, "src/shared"),
      },
    },
    server: {
      host: frontendServer.host,
      port: frontendServer.port,
      strictPort: true,
      proxy: {
        // /api –> localhost:7999
        [apiBasePath]: {
          target,
          changeOrigin: true,
        },
      },
    },
  };
});
