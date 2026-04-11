import React from "react";
import { SectionCard } from "@root/shared/components/ui/section-card";
import type { AutoFormRef } from "@root/core/form/form.types";
import { AutoForm } from "@root/core/form/auto-form";
import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined';
import { SafeButton } from "@shared/components/button/safe-button";
import { registerSlot } from "@root/core/module/registry";
import { useI18n } from "@root/core/i18n/use-i18n";

function AccountWidget() {
  const formAccountRef = React.useRef<AutoFormRef>(null);
  const { t } = useI18n();
  return (
    <SectionCard title={t("admin.auth.account.section_title")} extra={
      <SafeButton variant="contained" startIcon={<SaveOutlinedIcon />} onClick={() => formAccountRef.current?.submit()}>
        {t("admin.general.save_button")}
      </SafeButton>
    }>
      <AutoForm name="account" ref={formAccountRef} />
    </SectionCard>
  );
}

registerSlot({
  id: "account",
  name: "auth:left",
  priority: 2,
  render: () => <AccountWidget />,
});
