import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import { Box, Tooltip } from "@mui/material";
import { useI18n } from "@root/core/i18n/use-i18n";
import { registerSlot } from "@root/core/module/registry";
import { useAuthStore } from "@store/auth-store";

export function LogoutToolbarWidget() {
  const { t } = useI18n();
  const logout = useAuthStore((s) => s.logout);

  return (
    <Tooltip title={t("admin.toolbar.logout.tooltip")}>
      <Box
        aria-label={t("admin.toolbar.logout.tooltip")}
        onClick={async () => await logout()}
        sx={{
          position: "relative",
          display: "inline-flex",
          alignItems: "center",
          justifyContent: "center",
          cursor: "pointer",
        }}
      >
        <LogoutRoundedIcon />
      </Box>
    </Tooltip>
  );
}

registerSlot({
  id: "logout-toolbar",
  name: "toolbar",
  render: () => <LogoutToolbarWidget />,
  priority: 96,
});
