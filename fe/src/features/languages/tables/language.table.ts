import { openFormDialog } from "@core/form/form-dialog.service";
import { l } from "@root/core/i18n/localized-text";
import { navigate } from "@core/navigation/navigate";
import { reloadTable } from "@core/table/table-reload";
import { registerTable } from "@core/table/table-registry";
import { createTableSchema, type ColumnDef, type FetchTableOpts } from "@core/table/table.types";
import { deleteLanguage, listLanguages } from "@features/languages/api/language.api";
import type { LanguageModel } from "@features/languages/model/language.model";

const columns: ColumnDef<LanguageModel>[] = [
  { key: "code", header: l("admin.languages.table.columns.code"), sortable: true, labelField: true },
  { key: "name", header: l("admin.languages.table.columns.name"), sortable: true },
  { key: "nativeName", header: l("admin.languages.table.columns.native_name"), sortable: true },
  { key: "isDefault", header: l("admin.languages.table.columns.is_default"), type: "boolean", sortable: true },
  { key: "active", header: l("admin.languages.table.columns.active"), type: "boolean", sortable: true },
  { key: "updatedAt", header: l("admin.languages.table.columns.updated_at"), type: "datetime", sortable: true },
];

registerTable("languages", () =>
  createTableSchema<LanguageModel>({
    columns,
    fetch: async (opts: FetchTableOpts) => await listLanguages(opts),
    initialPageSize: 10,
    initialSort: { by: "id", dir: "asc" },
    allowUpdating: ["languages.update"],
    allowDeleting: ["languages.delete"],
    onView(row) {
      navigate(`/languages/${row.id}`);
    },
    onEdit(row) {
      openFormDialog("language", { initial: { id: row.id } });
    },
    async onDelete(row) {
      await deleteLanguage(Number(row.id));
      reloadTable("languages");
    },
  })
);
