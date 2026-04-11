import { registerTable } from "@core/table/table-registry";
import { createTableSchema, type ColumnDef, type FetchTableOpts } from "@core/table/table.types";
import { l } from "@root/core/i18n/localized-text";
import type { RoleModel } from "@features/rbac/model/role.model";
import { fetchRoles } from "@root/features/rbac/api/rbac.api";
import { openFormDialog } from "@root/core/form/form-dialog.service";

const columns: ColumnDef<RoleModel>[] = [
  // { key: "id", header: "ID", width: 80, sortable: true },
  { key: "displayName", header: l("admin.rbac.roles.table.columns.display_name"), sortable: true },
  { key: "roleName", header: l("admin.rbac.roles.table.columns.role_name"), width: 220, sortable: true, },
  { key: "brief", header: l("admin.rbac.roles.table.columns.brief") },
];

registerTable("roles", () =>
  createTableSchema<RoleModel>({
    columns,
    fetch: async (opts: FetchTableOpts) => await fetchRoles(opts),
    initialPageSize: 10,
    initialSort: { by: "id", dir: "asc" },
    onEdit(row) {
      openFormDialog("role", { initial: { id: row.id } });
    },
  })
);
