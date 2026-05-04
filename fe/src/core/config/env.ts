const baseAddress = import.meta.env.VITE_BASE_ADDRESS?.trim() ?? "";
const httpProto = import.meta.env.VITE_HTTP_PROTOCOL ?? "http";
const wsProto = import.meta.env.VITE_WS_PROTOCOL ?? "ws";
const wsEnabledEnv = import.meta.env.VITE_ENABLE_WEBSOCKET;
const explicitApiOrigin = import.meta.env.VITE_API_ORIGIN?.trim() ?? "";
const explicitWsOrigin = import.meta.env.VITE_WS_ORIGIN?.trim() ?? "";

function resolveSameOriginHttp() {
  if (typeof window === "undefined") return "";
  return `${window.location.protocol}//${window.location.host}`;
}

function resolveSameOriginWs() {
  if (typeof window === "undefined") return "";

  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}`;
}

function resolveApiOrigin() {
  if (explicitApiOrigin) return explicitApiOrigin;
  if (import.meta.env.PROD && !baseAddress) return resolveSameOriginHttp();

  const fallbackAddress = baseAddress || "127.0.0.1:7999";
  return `${httpProto}://${fallbackAddress}`;
}

function resolveWsOrigin() {
  if (explicitWsOrigin) return explicitWsOrigin;
  if (import.meta.env.PROD && !baseAddress) return resolveSameOriginWs();

  const fallbackAddress = baseAddress || "127.0.0.1:7999";
  return `${wsProto}://${fallbackAddress}`;
}

const apiOrigin = resolveApiOrigin();
const wsOrigin = resolveWsOrigin();

function parseBooleanEnv(value: string | undefined, fallback: boolean) {
  if (value == null) return fallback;

  const normalized = value.trim().toLowerCase();
  if (normalized === "true") return true;
  if (normalized === "false") return false;

  return fallback;
}

export const env = {
  mode: import.meta.env.MODE,
  wsEnabled: parseBooleanEnv(
    wsEnabledEnv,
    import.meta.env.MODE !== "development",
  ),
  apiOrigin,
  wsOrigin,
  apiBasePath: "/api",
  apiBaseUrl: `${apiOrigin}/api`,
  wsBaseUrl: `${wsOrigin}/ws`,
};
