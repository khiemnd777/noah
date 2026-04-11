import type { FieldDef } from "@core/form/types";
import type { FormSchema, SubmitDef } from "@core/form/form.types";
import { uploadImages } from "@root/core/form/image-upload-utils";
import { mapper } from "@root/core/mapper/auto-mapper";
import { updateDepartment } from "@features/settings/api/department.api";
import { registerForm } from "@root/core/form/form-registry";
import { l } from "@root/core/i18n/localized-text";
import { useAuthStore } from "@root/store/auth-store";
import { useI18nStore } from "@store/i18n-store";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

export function buildDepartmentSettingsSchema(): FormSchema {
  const fields: FieldDef[] = [
    {
      name: "name",
      label: l("admin.settings.department.fields.name.label"),
      kind: "text",
      rules: {
        required: translate("admin.settings.department.validation.name_required"),
        minLength: 2,
        maxLength: 120,
      },
    },
    {
      name: "address",
      label: l("admin.settings.department.fields.address.label"),
      kind: "text",
      rules: { maxLength: 300 },
    },
    {
      name: "phoneNumber",
      label: l("admin.settings.department.fields.phone_number.label"),
      kind: "text",
      placeholder: l("admin.settings.department.fields.phone_number.placeholder"),
      rules: {
        async: async (val: string | null) => {
          if (!val) return null;
          const ok = /^\+?\d{8,15}$/.test(val);
          return ok ? null : translate("admin.settings.department.validation.phone_invalid");
        },
      },
      helperText: l("admin.settings.department.fields.phone_number.helper_text"),
    },
    {
      name: "logo",
      label: l("admin.settings.department.fields.logo.label"),
      kind: "imageupload",
      accept: "image/*",
      maxFiles: 1,
      multipleFiles: false,
      helperText: l("admin.settings.department.fields.logo.helper_text"),
      uploader: uploadImages,
    },
    {
      name: "active",
      label: l("admin.settings.department.fields.active.label"),
      kind: "switch",
    },
  ];

  const submit: SubmitDef = {
    type: "fn",
    run: async (values) => {
      return updateDepartment(values.dto);
    },
  }

  return {
    fields,
    initialResolver() {
      return useAuthStore.getState().department;
    },
    async afterSaved() {
      await useAuthStore.getState().fetchDepartment();
    },
    toasts: {
      saved: l("admin.settings.department.messages.save_success"),
      failed: l("admin.settings.department.messages.save_failed"),
    },
    submit,
    hooks: {
      mapToDto: (v) => mapper.map("MyDepartment", v, "model_to_dto"),
    }
  };
}

registerForm("department-settings", buildDepartmentSettingsSchema);
