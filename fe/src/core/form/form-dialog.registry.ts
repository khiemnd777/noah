import type { FormSchema } from "@core/form/form.types";
import type { LocalizedText } from "@root/core/i18n/localized-text";

export type ModeText = LocalizedText | { create: LocalizedText; update: LocalizedText };
export type TitleProp = React.ReactNode | ModeText;

export type FormDialogDefaults = {
  title?: TitleProp;
  confirmText?: ModeText;
  cancelText?: LocalizedText;
  maxWidth?: "xs" | "sm" | "md" | "lg";
};

type Builder = () => FormSchema;

const dialogRegistry = new Map<string, { build?: Builder; defaults?: FormDialogDefaults }>();
const dialogCache = new Map<string, FormSchema>();

export function registerFormDialog(name: string, build?: Builder, defaults?: FormDialogDefaults) {
  dialogRegistry.set(name, { build, defaults });
  dialogCache.delete(name); // cập nhật thì xoá cache
}

export function getFormDialogBuilder(name: string): FormSchema | undefined {
  const cached = dialogCache.get(name);
  if (cached) return cached;

  const b = dialogRegistry.get(name)?.build;
  if (!b) return undefined;

  const schema = b();
  dialogCache.set(name, schema);
  return schema;
}

export function getFormDialogDefaults(name: string): FormDialogDefaults | undefined {
  return dialogRegistry.get(name)?.defaults;
}

/* Ví dụ:
// form
registerForm("account", buildAccountSchema);

// dialog defaults (KHÔNG cần build lại nếu đã registerForm ở trên)
registerFormDialog("account", undefined, { title: "Tạo mới tài khoản", confirmText: "Lưu" });

// Nếu bạn muốn pass builder (không bắt buộc):
registerFormDialog("account", buildAccountSchema, { title: "Tạo mới tài khoản" });
*/
