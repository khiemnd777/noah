import type { FieldDef } from "@core/form/types";
import type { FormSchema } from "@core/form/form.types";
import { uploadImages } from "@core/form/image-upload-utils";
import { l } from "@root/core/i18n/localized-text";
import { mapper } from "@core/mapper/auto-mapper";
import type { StaffModel } from "@features/staff/model/staff.model";
import { create, existsEmail, existsPhone, id, update } from "@features/staff/api/staff.api";
import { reloadTable } from "@core/table/table-reload";
import { openFormDialog } from "@core/form/form-dialog.service";
import { fetchRolesByUserId, search as searchRoles } from "@root/features/rbac/api/rbac.api";
import { useI18nStore } from "@store/i18n-store";

type Options = {
  withPassword: boolean;
  passwordRequired?: boolean;
};

function translate(key: string, fallback?: string): string {
  return useI18nStore.getState().resources[key] ?? fallback ?? key;
}

function passwordField(opts: Options): FieldDef {
  return {
    name: "password",
    label: l("admin.staff.form.fields.password.label"),
    kind: "password",
    rules: {
      ...(opts.withPassword && opts.passwordRequired ? {
        required: translate("admin.staff.validation.password_required"),
      } : {}),
      minLength: 6,
      maxLength: 128
    },
  };
}

function commonFields(): FieldDef[] {
  return [
    {
      name: "name",
      label: l("admin.staff.form.fields.name.label"),
      kind: "text",
      rules: { required: translate("admin.staff.validation.name_required"), maxLength: 50 },
    },
    {
      name: "email",
      label: l("admin.staff.form.fields.email.label"),
      kind: "email",
      rules: {
        required: translate("admin.staff.validation.email_required"),
        maxLength: 300,
        async: async (val: string | null, { id }) => {
          if (!val) return null;
          const existed = await existsEmail({ id, email: val });
          return existed ? translate("admin.staff.validation.email_exists").replace("{value}", val) : null;
        },
      },
    },
    {
      name: "phone",
      label: l("admin.staff.form.fields.phone.label"),
      kind: "text",
      placeholder: l("admin.staff.form.fields.phone.placeholder"),
      rules: {
        async: async (val: string | null, { id }) => {
          if (!val) return null;
          const ok = /^\+?\d{8,15}$/.test(val);
          if (!ok) return translate("admin.staff.validation.phone_invalid");
          const existed = await existsPhone({ id, phone: val });
          return existed ? translate("admin.staff.validation.phone_exists").replace("{value}", val) : null;
        },
      },
      helperText: l("admin.staff.form.fields.phone.helper_text"),
    },
    {
      name: "",
      label: "",
      kind: "metadata",
      metadata: {
        collection: "staff",
        mode: "whole",
      }
    },
    {
      name: "avatar",
      label: l("admin.staff.form.fields.avatar.label"),
      kind: "imageupload",
      accept: "image/*",
      maxFiles: 1,
      multipleFiles: false,
      helperText: l("admin.staff.form.fields.avatar.helper_text"),
      uploader: uploadImages,
    },
    {
      name: "active",
      label: l("admin.staff.form.fields.active.label"),
      kind: "switch",
      defaultValue: true,
    },
    // ---- Roles ----
    {
      name: "roleIds",
      label: l("admin.staff.form.fields.role_ids.label"),
      kind: "searchlist",
      placeholder: l("admin.staff.form.fields.role_ids.placeholder"),
      fullWidth: true,

      getOptionLabel: (d: any) => d?.displayName,
      getOptionValue: (d: any) => d?.id,

      async searchPage(kw: string, page: number, limit: number) {
        const searched = await searchRoles({ keyword: kw, limit, page, orderBy: "display_name" });
        return searched.items;
      },
      pageLimit: 20,

      async hydrateByIds(ids: Array<number | string>, values: Record<string, any>) {
        if (!ids || ids.length === 0) return [];
        const table = await fetchRolesByUserId(values.id, { limit: 20, page: 1, orderBy: "display_name" });
        const set = new Set(ids.map(String));
        return (table.items ?? []).filter((d: any) => set.has(String(d.id)));
      },

      async fetchList(values: Record<string, any>) {
        const table = await fetchRolesByUserId(values.id, { limit: 20, page: 1, orderBy: "display_name" });
        return table.items;
      },

      onDragEnd(items) {
        console.log(items);
      },

      renderItem: (d: any) => <> {d.displayName} </>,
      disableDelete: (d: any) => d.locked === true,
      onOpenCreate: () => openFormDialog("role"),
      autoLoadAllOnMount: true,
    },
  ];
}

export function buildStaffSchemaShared(opts: Options): FormSchema {
  const fields = [...commonFields()];
  if (opts.withPassword) {
    // chèn password ngay sau phone (index 2 là phone, vậy password ở 3)
    fields.splice(3, 0, passwordField(opts));
  }

  return {
    idField: "id",
    fields,
    submit: {
      create: {
        type: "fn",
        run: async (values) => {
          await create(values.dto as StaffModel);
          return values.dto;
        },
      },
      update: {
        type: "fn",
        run: async (values) => {
          await update(values.dto as StaffModel);
          return values.dto;
        },
      },
    },
    toasts: {
      saved: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.staff.messages.create_success"
            : "admin.staff.messages.update_success"
        ).replace("{name}", String(values?.name ?? "")),
      failed: ({ mode, values }) =>
        translate(
          mode === "create"
            ? "admin.staff.messages.create_failed"
            : "admin.staff.messages.update_failed"
        ).replace("{name}", String(values?.name ?? "")),
    },
    async initialResolver(data: any) {
      if (data) {
        return await id(data.id);
      }
      return {};
    },
    async afterSaved() {
      reloadTable("staffs");
    },
    hooks: {
      mapToDto: (v) => mapper.map("Staff", v, "model_to_dto"),
    },
  };
}
