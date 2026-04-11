import * as React from "react";
import type { FormSchema, FormMode, ModeText } from "@core/form/form.types";
import { isI18nText } from "@root/core/i18n/localized-text";

export function resolveMode(schema: FormSchema, initialVals: any): FormMode {
  const idField = schema.idField ?? "id";
  if (schema.modeResolver) return schema.modeResolver(initialVals ?? {});
  const id = initialVals?.[idField];
  return id !== null && id !== undefined && id !== "" ? "update" : "create";
}

type Ctx = { mode: FormMode; values: any; result?: any; initial?: any };

export function pickModeText(t: ModeText | undefined, ctx: Ctx) {
  if (!t) return undefined;
  if (typeof t === "string" || isI18nText(t)) return t;
  if (typeof t === "function") return t({ mode: ctx.mode, values: ctx.values, result: ctx.result });
  return t[ctx.mode];
}

// Title có thể là ReactNode tĩnh *hoặc* ModeText; nếu là ModeText → render theo ctx
export function resolveTitle(
  input: React.ReactNode | ModeText | undefined | null,
  ctx: Ctx
): unknown {
  if (input == null) return undefined;
  if (isI18nText(input)) return input;
  if (React.isValidElement(input) || typeof input !== "object") return input;
  return pickModeText(input as ModeText, ctx);
}
