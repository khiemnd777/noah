import type { FieldDef } from "@core/form/types";
import type { FormSchema, SubmitDef } from "@core/form/form.types";
import { uploadImages } from "@core/form/image-upload-utils";
import { l } from "@root/core/i18n/localized-text";
import { mapper } from "@root/core/mapper/auto-mapper";
import { existsEmail, existsPhone, updateMe } from "@root/core/network/me.api";
import type { MeModel } from "@core/auth/auth.types";
import { registerForm } from "@core/form/form-registry";
import { useAuthStore } from "@root/store/auth-store";
import { useI18nStore } from "@store/i18n-store";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

export function buildAccountSchema(): FormSchema {
  const fields: FieldDef[] = [
    {
      name: "name",
      label: l("admin.auth.account.fields.name.label"),
      kind: "text",
      rules: {
        required: translate("admin.auth.account.validation.name_required"),
        maxLength: 50,
      },
    },
    {
      name: "email",
      label: l("admin.auth.account.fields.email.label"),
      kind: "email",
      rules: {
        required: translate("admin.auth.account.validation.email_required"),
        maxLength: 300,
        async: async (val: string | null) => {
          if (!val) return null;
          if (val) {
          }
          const existed = await existsEmail(val);
          return existed ? translate("admin.auth.account.validation.email_exists").replace("{value}", val) : null;
        }
      },
    },
    {
      name: "phone",
      label: l("admin.auth.account.fields.phone.label"),
      kind: "text",
      placeholder: l("admin.auth.account.fields.phone.placeholder"),
      rules: {
        async: async (val: string | null) => {
          if (!val) return null;
          const ok = /^\+?\d{8,15}$/.test(val);
          if (!ok) {
            return translate("admin.auth.account.validation.phone_invalid");
          }
          const existed = await existsPhone(val);
          if (existed) {
            return translate("admin.auth.account.validation.phone_exists").replace("{value}", val);
          }
          return null;
        },
      },
      helperText: l("admin.auth.account.fields.phone.helper_text"),
    },
    {
      name: "avatar",
      label: l("admin.auth.account.fields.avatar.label"),
      kind: "imageupload",
      accept: "image/*",
      maxFiles: 1,
      multipleFiles: false,
      helperText: l("admin.auth.account.fields.avatar.helper_text"),
      uploader: uploadImages,
    },
  ];

  const submit: SubmitDef = {
    type: "fn",
    run: async (values) => {
      await updateMe(values.dto as MeModel);
    }
  };

  return {
    fields,
    initialResolver() {
      return useAuthStore.getState().user;
    },
    async afterSaved() {
      await useAuthStore.getState().fetchMe();
    },
    toasts: {
      saved: l("admin.auth.account.messages.save_success"),
      failed: l("admin.auth.account.messages.save_failed"),
    },
    submit,
    hooks: {
      mapToDto: (v) => mapper.map("Me", v, "model_to_dto"),
    }
  };
}

registerForm("account", buildAccountSchema);
