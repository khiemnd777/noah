import { ThemeToggle } from "@features/settings/components/theme-toggle";
import { Typography } from "@mui/material";
import { AutoGrid } from "@root/shared/components/ui/auto-grid";
import { useI18n } from "@root/core/i18n/use-i18n";

export default function SettingsForm() {
  const { t } = useI18n();

  return (
    <>
      <AutoGrid scheme="equal">
        <Typography fontWeight={500}>{t("admin.settings.ui.theme.label")}</Typography>
        <ThemeToggle />
      </AutoGrid>
    </>
  );
}
