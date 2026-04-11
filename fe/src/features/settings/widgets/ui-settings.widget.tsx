import { SectionCard } from "@shared/components/ui/section-card";
import { registerSlot } from "@root/core/module/registry";
import SettingsForm from "@features/settings/components/common-settings-form";
import { useI18n } from "@root/core/i18n/use-i18n";

function UISettingsWidget() {
  const { t } = useI18n();

  return (
    <>
      <SectionCard title={t("admin.settings.ui.section_title")}>
        <SettingsForm />
      </SectionCard>
    </>
  );
}

registerSlot({
  id: "ui-settings",
  name: "settings:right",
  render: () => <UISettingsWidget />,
});
