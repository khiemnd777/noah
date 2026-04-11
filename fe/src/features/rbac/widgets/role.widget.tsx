import { SectionCard } from "@root/shared/components/ui/section-card";
import { AutoTable } from "@root/core/table/auto-table";
import { Button } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import { openFormDialog } from "@root/core/form/form-dialog.service";
import { useI18n } from "@root/core/i18n/use-i18n";
import { registerSlot } from "@root/core/module/registry";

function RoleWidget() {
  const { t } = useI18n();

  return (
    <>
      <SectionCard title={t("admin.rbac.roles.section_title")} extra={
        <>
          <Button variant="outlined" startIcon={<AddIcon />} onClick={() => {
            openFormDialog("role");
          }} >{t("admin.general.create_button")}</Button>
        </>
      }>
        <AutoTable name="roles" />
      </SectionCard>
    </>
  );
}

registerSlot({
  id: "role",
  name: "rbac:left",
  priority: 2,
  render: () => RoleWidget(),
});
