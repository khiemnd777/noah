import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { registerFormDialog } from "@core/form/form-dialog.registry";
import { l } from "@root/core/i18n/localized-text";
import { useI18nStore } from "@store/i18n-store";
import { importLanguageXml } from "@features/languages/api/language.api";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

function buildLanguageImportFields(): FieldDef[] {
  return [
    {
      name: "languageId",
      label: l("admin.languages.import.fields.language_id.label"),
      kind: "number",
      showIf: () => false,
    },
    {
      name: "file",
      label: l("admin.languages.import.fields.file.label"),
      kind: "fileupload",
      accept: ".xml,text/xml,application/xml",
      maxFiles: 1,
      multipleFiles: false,
      rules: {
        required: translate("admin.languages.import.validation.file_required"),
      },
      helperText: l("admin.languages.import.fields.file.helper_text"),
    },
  ];
}

export function buildLanguageImportSchema(): FormSchema {
  return {
    fields: buildLanguageImportFields(),
    submit: {
      create: {
        type: "fn",
        run: async (values) => {
          const languageId = Number(values.languageId ?? values.id ?? 0);
          const files = Array.isArray(values.file) ? values.file : [];
          const file = files[0] as File | undefined;

          if (!languageId) {
            throw new Error(translate("admin.languages.import.validation.language_id_missing"));
          }
          if (!file) {
            throw new Error(translate("admin.languages.import.validation.file_required"));
          }

          await importLanguageXml(languageId, file);
          return { languageId, fileName: file.name };
        },
      },
      update: {
        type: "fn",
        run: async (values) => {
          const languageId = Number(values.languageId ?? values.id ?? 0);
          const files = Array.isArray(values.file) ? values.file : [];
          const file = files[0] as File | undefined;

          if (!languageId) {
            throw new Error(translate("admin.languages.import.validation.language_id_missing"));
          }
          if (!file) {
            throw new Error(translate("admin.languages.import.validation.file_required"));
          }

          await importLanguageXml(languageId, file);
          return { languageId, fileName: file.name };
        },
      },
    },
    toasts: {
      saved: ({ result }) =>
        translate("admin.languages.import.messages.success").replace("{fileName}", String(result?.fileName ?? "")),
      failed: () => translate("admin.languages.import.messages.failed"),
    },
  };
}

registerFormDialog("language-import", buildLanguageImportSchema, {
  title: l("admin.languages.import.dialog.title"),
  confirmText: l("admin.general.import_button"),
  cancelText: l("admin.general.cancel_button"),
  maxWidth: "sm",
});
