import { openFormDialog } from "@core/form/form-dialog.service";
import { navigate } from "@core/navigation/navigate";
import { reloadTable } from "@core/table/table-reload";
import { registerTable } from "@core/table/table-registry";
import { createTableSchema, type ColumnDef, type FetchTableOpts } from "@core/table/table.types";
import { deleteLanguage, listLanguages } from "@features/languages/api/language.api";
import type { LanguageModel } from "@features/languages/model/language.model";

const columns: ColumnDef<LanguageModel>[] = [
  { key: "code", header: "Mã", sortable: true, labelField: true },
  { key: "name", header: "Tên", sortable: true },
  { key: "nativeName", header: "Tên bản địa", sortable: true },
  { key: "isDefault", header: "Mặc định", type: "boolean", sortable: true },
  { key: "active", header: "Kích hoạt", type: "boolean", sortable: true },
  { key: "updatedAt", header: "Cập nhật lúc", type: "datetime", sortable: true },
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
