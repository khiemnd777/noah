import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { registerFormDialog } from "@core/form/form-dialog.registry";
import { registerForm } from "@core/form/form-registry";
import { reloadTable } from "@core/table/table-reload";
import { l } from "@root/core/i18n/localized-text";
import { useI18nStore } from "@store/i18n-store";
import { createLanguage, getLanguageById, updateLanguage } from "@features/languages/api/language.api";
import type { LanguageModel } from "@features/languages/model/language.model";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

function buildLanguageFields(): FieldDef[] {
  return [
    {
      name: "code",
      label: l("admin.languages.form.fields.code.label"),
      kind: "text",
      placeholder: l("admin.languages.form.fields.code.placeholder"),
      rules: {
        required: translate("admin.languages.validation.code_required"),
        maxLength: 32,
        pattern: {
          regex: /^[A-Za-z0-9_-]+(?:\.[A-Za-z0-9_-]+)?$/,
          message: translate("admin.languages.validation.code_invalid"),
        },
      },
    },
    {
      name: "name",
      label: l("admin.languages.form.fields.name.label"),
      kind: "text",
      rules: {
        required: translate("admin.languages.validation.name_required"),
        maxLength: 120,
      },
    },
    {
      name: "nativeName",
      label: l("admin.languages.form.fields.native_name.label"),
      kind: "text",
      rules: {
        required: translate("admin.languages.validation.native_name_required"),
        maxLength: 120,
      },
    },
    {
      name: "isDefault",
      label: l("admin.languages.form.fields.is_default.label"),
      kind: "switch",
      defaultValue: false,
    },
    {
      name: "active",
      label: l("admin.languages.form.fields.active.label"),
      kind: "switch",
      defaultValue: true,
    },
  ];
}

function buildLanguageBaseSchema(): FormSchema {
  return {
    idField: "id",
    fields: buildLanguageFields(),
    submit: {
      create: {
        type: "fn",
        run: async (values) => {
          return await createLanguage(values.dto as LanguageModel);
        },
      },
      update: {
        type: "fn",
        run: async (values) => {
          return await updateLanguage(values.dto as LanguageModel);
        },
      },
    },
    async initialResolver(data?: Partial<LanguageModel>) {
      if (data?.resources) return data;
      if (data?.id) return await getLanguageById(data.id);
      return {
        active: true,
        isDefault: false,
        resources: [],
      };
    },
    toasts: {
      saved: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.languages.messages.create_success"
            : "admin.languages.messages.update_success"
        ).replace("{name}", String(values?.name ?? "")),
      failed: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.languages.messages.create_failed"
            : "admin.languages.messages.update_failed"
        ).replace("{name}", String(values?.name ?? "")),
    },
    async afterSaved() {
      reloadTable("languages");
    },
  };
}

export function buildLanguageDialogSchema() {
  return buildLanguageBaseSchema();
}

export function buildLanguageDetailSchema() {
  return buildLanguageBaseSchema();
}

registerFormDialog("language", buildLanguageDialogSchema, {
  title: {
    create: l("admin.languages.dialog.create_title"),
    update: l("admin.languages.dialog.update_title"),
  },
  confirmText: {
    create: l("admin.general.create_button"),
    update: l("admin.general.save_button"),
  },
  cancelText: l("admin.general.cancel_button"),
  maxWidth: "md",
});

registerForm("language-detail", buildLanguageDetailSchema);
