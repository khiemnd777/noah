import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { registerFormDialog } from "@core/form/form-dialog.registry";
import { registerForm } from "@core/form/form-registry";
import { reloadTable } from "@core/table/table-reload";
import { createLanguage, getLanguageById, updateLanguage } from "@features/languages/api/language.api";
import type { LanguageModel } from "@features/languages/model/language.model";

function buildLanguageFields(): FieldDef[] {
  return [
    {
      name: "code",
      label: "Mã ngôn ngữ",
      kind: "text",
      placeholder: "vi, en, ja, kr, ...",
      rules: {
        required: "Yêu cầu nhập mã ngôn ngữ",
        maxLength: 32,
        pattern: {
          regex: /^[A-Za-z0-9_-]+(?:\.[A-Za-z0-9_-]+)?$/,
          message: "Mã ngôn ngữ chỉ gồm chữ, số, gạch ngang, gạch dưới hoặc dấu chấm.",
        },
      },
    },
    {
      name: "name",
      label: "Tên hiển thị",
      kind: "text",
      rules: {
        required: "Yêu cầu nhập tên hiển thị",
        maxLength: 120,
      },
    },
    {
      name: "nativeName",
      label: "Tên bản địa",
      kind: "text",
      rules: {
        required: "Yêu cầu nhập tên bản địa",
        maxLength: 120,
      },
    },
    {
      name: "isDefault",
      label: "Ngôn ngữ mặc định",
      kind: "switch",
      defaultValue: false,
    },
    {
      name: "active",
      label: "Kích hoạt",
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
        mode === "create"
          ? `Tạo ngôn ngữ "${values?.name ?? ""}" thành công!`
          : `Cập nhật ngôn ngữ "${values?.name ?? ""}" thành công!`,
      failed: ({ mode, values }) =>
        mode === "create"
          ? `Tạo ngôn ngữ "${values?.name ?? ""}" thất bại, xin thử lại!`
          : `Cập nhật ngôn ngữ "${values?.name ?? ""}" thất bại, xin thử lại!`,
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
  title: { create: "Thêm ngôn ngữ", update: "Cập nhật ngôn ngữ" },
  confirmText: { create: "Thêm", update: "Lưu" },
  cancelText: "Thoát",
  maxWidth: "md",
});

registerForm("language-detail", buildLanguageDetailSchema);
