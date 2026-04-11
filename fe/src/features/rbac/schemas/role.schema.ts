import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { mapper } from "@core/mapper/auto-mapper";
import { registerFormDialog } from "@core/form/form-dialog.registry";
import { l } from "@root/core/i18n/localized-text";
import { createRole, fetchRoleByID, updateRole } from "@features/rbac/api/rbac.api";
import type { RoleModel } from "@features/rbac/model/role.model";
import { reloadTable } from "@core/table/table-reload";
import { EV_RBAC_MATRIX_INVALIDATE } from "@features/rbac/model/rbac.events";
import { invalidate } from "@core/module/event-invalidation";
import { useI18nStore } from "@store/i18n-store";

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

export function buildRoleSchema(): FormSchema {
  const fields: FieldDef[] = [
    {
      name: "displayName",
      label: l("admin.rbac.roles.form.fields.display_name.label"),
      kind: "text",
      rules: {
        required: translate("admin.rbac.roles.validation.display_name_required"),
        maxLength: 50,
      },
    },
    {
      name: "roleName",
      label: l("admin.rbac.roles.form.fields.role_name.label"),
      kind: "text",
      rules: {
        required: translate("admin.rbac.roles.validation.role_name_required"),
        maxLength: 20,
      },
      // derive: {
      //   field: "displayName",
      //   mode: "whenEmpty",
      //   map: (srcVal) => slugify(String(srcVal ?? "")),
      // },
    },
    {
      name: "brief",
      label: l("admin.rbac.roles.form.fields.brief.label"),
      kind: "textarea",
      rules: { maxLength: 300 },
    },
  ];

  return {
    idField: "id",
    fields,
    submit: {
      create: {
        type: "fn",
        run: async (values) => {
          await createRole(values.dto as RoleModel);
          return values.dto;
        },
      },
      update: {
        type: "fn",
        run: async (values) => {
          await updateRole(values.dto as RoleModel);
          return values.dto;
        },
      },
    },

    toasts: {
      saved: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.rbac.roles.messages.create_success"
            : "admin.rbac.roles.messages.update_success"
        ).replace("{name}", String(values?.displayName ?? "")),
      failed: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.rbac.roles.messages.create_failed"
            : "admin.rbac.roles.messages.update_failed"
        ).replace("{name}", String(values?.displayName ?? "")),
    },

    async initialResolver(data: any) {
      if (data) {
        return await fetchRoleByID(data.id);
      }
      return {};
    },

    async afterSaved() {
      // reload table
      reloadTable("roles");
      // invalidate
      invalidate(EV_RBAC_MATRIX_INVALIDATE, { reason: "role:save" });
    },

    hooks: {
      mapToDto: (v) => mapper.map("Role", v, "model_to_dto"),
    },
  };
}

registerFormDialog("role", buildRoleSchema, {
  title: {
    create: l("admin.rbac.roles.dialog.create_title"),
    update: l("admin.rbac.roles.dialog.update_title"),
  },
  confirmText: {
    create: l("admin.general.create_button"),
    update: l("admin.general.save_button"),
  },
  cancelText: l("admin.general.cancel_button"),
});
