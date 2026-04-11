import { registerFormDialog } from "@core/form/form-dialog.registry";
import { l } from "@root/core/i18n/localized-text";
import { buildStaffSchemaShared } from "./staff.schema.shared";

export function buildStaffEditDialogSchema() {
  return buildStaffSchemaShared({ withPassword: false });
}

registerFormDialog("staff-edit-dialog", buildStaffEditDialogSchema, {
  title: { create: l("admin.staff.dialog.create_title"), update: l("admin.staff.dialog.update_title") },
  confirmText: { create: l("admin.general.create_button"), update: l("admin.general.save_button") },
  cancelText: l("admin.general.cancel_button"),
});
