import React from "react";
import { SectionCard } from "@root/shared/components/ui/section-card";
import type { AutoFormRef } from "@root/core/form/form.types";
import { AutoForm } from "@root/core/form/auto-form";
import { SafeButton } from "@shared/components/button/safe-button";
import { registerSlot } from "@root/core/module/registry";
import ChangeCircleOutlinedIcon from '@mui/icons-material/ChangeCircleOutlined';
import { useI18n } from "@root/core/i18n/use-i18n";

function ChangePasswordWidget() {
  const formAccountChangePasswordRef = React.useRef<AutoFormRef>(null);
  const { t } = useI18n();
  return (
    <SectionCard title={t("admin.auth.change_password.section_title")} extra={
      <SafeButton variant="contained" startIcon={<ChangeCircleOutlinedIcon />} onClick={() => formAccountChangePasswordRef.current?.submit()}>
        {t("admin.auth.change_password.submit_button")}
      </SafeButton>
    }>
      <AutoForm name="account-change-password" ref={formAccountChangePasswordRef} />
    </SectionCard>
  );
}

registerSlot({
  id: "change-password",
  name: "auth:right",
  priority: 1,
  render: () => <ChangePasswordWidget />,
});
