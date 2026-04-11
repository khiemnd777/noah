import SettingsRoundedIcon from "@mui/icons-material/SettingsRounded";
import { Box, Tooltip } from "@mui/material";
import { useI18n } from "@root/core/i18n/use-i18n";
import { navigate } from "@root/core/navigation/navigate";

export default function SettingIcon() {
  const { t } = useI18n();

  return (
    <Tooltip title={t("admin.toolbar.settings.tooltip", "Cài đặt")}>
      <Box
        aria-label={t("admin.toolbar.settings.tooltip", "Cài đặt")}
        onClick={() => navigate("/settings")}
        sx={{
          position: "relative",
          display: "inline-flex",
          alignItems: "center",
          justifyContent: "center",
          cursor: "pointer",
        }}
      >
        <SettingsRoundedIcon />
      </Box>
    </Tooltip>
  );
}
