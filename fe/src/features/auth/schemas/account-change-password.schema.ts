import type { FieldDef } from "@core/form/types";
import type { FormSchema, SubmitDef } from "@core/form/form.types";
import { changeMyPassword } from "@root/core/network/me.api";
import { l } from "@root/core/i18n/localized-text";
import { registerForm } from "@root/core/form/form-registry";
import { useI18nStore } from "@store/i18n-store";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

export function buildAccountChangePasswordSchema(): FormSchema {
  const fields: FieldDef[] = [
    {
      name: "password",
      kind: "change-password",
      label: l("admin.auth.change_password.fields.password.label"),
      currentLabel: l("admin.auth.change_password.fields.password.current_label"),
      newLabel: l("admin.auth.change_password.fields.password.new_label"),
      confirmLabel: l("admin.auth.change_password.fields.password.confirm_label"),
      rules: {
        required: translate("admin.auth.change_password.validation.required"),
      },
      passwordRules: {
        disallowReuseCurrent: false,
        minLength: 8,
        requireDigit: true,
        requireUpper: false,
        requireLower: false,
      },
    },
  ];

  const submit: SubmitDef = {
    type: "fn",
    run: async (values) => {
      const { current, password } = values.dto.password;
      await changeMyPassword(current, password);
    }
  };

  return {
    fields,
    toasts: {
      saved: l("admin.auth.change_password.messages.save_success"),
      failed: l("admin.auth.change_password.messages.save_failed"),
    },
    submit,
  };
}

registerForm("account-change-password", buildAccountChangePasswordSchema);
