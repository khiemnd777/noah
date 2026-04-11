import { Button, Tooltip } from "@mui/material";
import { useColorScheme } from "@mui/material/styles";
import DarkModeRoundedIcon from "@mui/icons-material/DarkModeRounded";
import LightModeRoundedIcon from "@mui/icons-material/LightModeRounded";
import SettingsBrightnessRoundedIcon from "@mui/icons-material/SettingsBrightnessRounded";
import { useI18n } from "@root/core/i18n/use-i18n";

export function ThemeToggle() {
  const { t } = useI18n();
  const { mode, setMode, systemMode } = useColorScheme();
  const effective = mode === "system" ? systemMode : mode; // light | dark

  const next =
    mode === "light" ? "dark" : mode === "dark" ? "system" : "light";

  const icon =
    effective === "dark" ? (
      <DarkModeRoundedIcon fontSize="small" />
    ) : effective === "light" ? (
      <LightModeRoundedIcon fontSize="small" />
    ) : (
      <SettingsBrightnessRoundedIcon fontSize="small" />
    );

  const currentModeLabel = t(`admin.settings.ui.theme.mode.${mode}`);
  const nextModeLabel = t(`admin.settings.ui.theme.mode.${next}`);

  return (
    <Tooltip title={t("admin.settings.ui.theme.tooltip").replace("{current}", currentModeLabel).replace("{next}", nextModeLabel)}>
      <Button
        size="small"
        variant="outlined"
        onClick={() => setMode(next as any)}
        startIcon={icon}
        sx={{ textTransform: "capitalize" }}
      >
        {currentModeLabel}
      </Button>
    </Tooltip>
  );
}
