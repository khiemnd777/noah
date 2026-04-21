/* eslint-disable react-refresh/only-export-components */
import * as React from "react";
import type { LocalizedText } from "@core/i18n/localized-text";

export type RouteMeta = {
  key: string;
  label?: LocalizedText;
  title?: LocalizedText;
  subtitle?: LocalizedText;
  path: string;
  hidden?: boolean;
  parentKey?: string;
  parentPath?: string;
  isDetailRoute?: boolean;
  extra?: Record<string, unknown>;
};

const Ctx = React.createContext<RouteMeta | null>(null);

export function RouteMetaProvider({ meta, children }: { meta: RouteMeta; children: React.ReactNode }) {
  return <Ctx.Provider value={meta}>{children}</Ctx.Provider>;
}

export function useRouteMeta() {
  const ctx = React.useContext(Ctx);
  if (!ctx) throw new Error("useRouteMeta must be used within <RouteMetaProvider>");
  return ctx;
}
