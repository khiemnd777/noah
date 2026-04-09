import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { registerFormDialog } from "@core/form/form-dialog.registry";
import { importLanguageXml } from "@features/languages/api/language.api";

function buildLanguageImportFields(): FieldDef[] {
  return [
    {
      name: "languageId",
      label: "Language ID",
      kind: "number",
      showIf: () => false,
    },
    {
      name: "file",
      label: "Chọn file XML",
      kind: "fileupload",
      accept: ".xml,text/xml,application/xml",
      maxFiles: 1,
      multipleFiles: false,
      rules: {
        required: "Yêu cầu chọn file XML",
      },
      helperText: "Hỗ trợ file XML chứa các resource keys dạng admin.{module}.*",
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
            throw new Error("Thiếu mã ngôn ngữ để nhập XML.");
          }
          if (!file) {
            throw new Error("Yêu cầu chọn file XML.");
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
            throw new Error("Thiếu mã ngôn ngữ để nhập XML.");
          }
          if (!file) {
            throw new Error("Yêu cầu chọn file XML.");
          }

          await importLanguageXml(languageId, file);
          return { languageId, fileName: file.name };
        },
      },
    },
    toasts: {
      saved: ({ result }) => `Nhập XML "${result?.fileName ?? ""}" thành công!`,
      failed: () => "Nhập XML thất bại.",
    },
  };
}

registerFormDialog("language-import", buildLanguageImportSchema, {
  title: "Nhập XML resources",
  confirmText: "Nhập",
  cancelText: "Thoát",
  maxWidth: "sm",
});
