import { Button } from "@mui/material";
import { SectionCard } from "@root/shared/components/ui/section-card";
import AddIcon from '@mui/icons-material/Add';
import { openFormDialog } from "@core/form/form-dialog.service";
import { AutoTable } from "@core/table/auto-table";
import { useI18n } from "@root/core/i18n/use-i18n";
import { registerSlot } from "@root/core/module/registry";
import { IfPermission } from "@root/core/auth/if-permission";

function StaffWidget() {
  const { t } = useI18n();

  return (
    <>
      <SectionCard extra={
        <>
          <IfPermission permissions={["staff.create"]}>
            <Button variant="outlined" startIcon={<AddIcon />} onClick={() => {
              openFormDialog("staff-create");
            }} >{t("admin.general.create_button")}</Button>
          </IfPermission>
        </>
      }>
        <AutoTable name="staffs" />
      </SectionCard>
    </>
  );
}

registerSlot({
  id: "staff",
  name: "staff:left",
  render: () => <StaffWidget />,
})
