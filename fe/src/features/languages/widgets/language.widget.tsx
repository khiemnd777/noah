import AddIcon from "@mui/icons-material/Add";
import { Button } from "@mui/material";
import { IfPermission } from "@root/core/auth/if-permission";
import { openFormDialog } from "@root/core/form/form-dialog.service";
import { useAdminI18n } from "@root/core/i18n/use-admin-i18n";
import { registerSlot } from "@root/core/module/registry";
import { AutoTable } from "@root/core/table/auto-table";
import { SectionCard } from "@root/shared/components/ui/section-card";

export function LanguageWidget() {
  const { t } = useAdminI18n();

  return (
    <SectionCard
      title={t("admin.languages.list.section_title", "Ngôn ngữ")}
      extra={
        <IfPermission permissions={["languages.create"]}>
          <Button variant="outlined" startIcon={<AddIcon />} onClick={() => openFormDialog("language")}>
            {t("admin.languages.actions.create", "Thêm ngôn ngữ")}
          </Button>
        </IfPermission>
      }
    >
      <AutoTable name="languages" />
    </SectionCard>
  );
}

registerSlot({
  id: "languages",
  name: "languages:left",
  render: () => <LanguageWidget />,
});
