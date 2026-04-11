import { SectionCard } from "@shared/components/ui/section-card";
import React from "react";
import type { AutoFormRef } from "@core/form/form.types";
import { AutoForm } from "@root/core/form/auto-form";
import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined';
import { SafeButton } from "@shared/components/button/safe-button";
import { registerSlot } from "@root/core/module/registry";
import { IfPermission } from "@root/core/auth/if-permission";
import { useI18n } from "@root/core/i18n/use-i18n";

function DepartmentSettingsWidget() {
  const formRef = React.useRef<AutoFormRef>(null);
  const { t } = useI18n();

  return (
    <>
      <SectionCard title={t("admin.settings.department.section_title")} extra={
        <IfPermission permissions={["settings.update"]}>
          <SafeButton variant="contained" startIcon={<SaveOutlinedIcon />} onClick={() => formRef.current?.submit()}>{t("admin.general.save_button")}</SafeButton>
        </IfPermission>
      }>
        <AutoForm name="department-settings" ref={formRef} />
      </SectionCard>
    </>
  );
}

registerSlot({
  id: "department-settings",
  name: "settings:left",
  render: () => <DepartmentSettingsWidget />,
});
