import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import { Box, Tooltip } from "@mui/material";
import { useAdminI18n } from "@root/core/i18n/use-admin-i18n";
import { registerSlot } from "@root/core/module/registry";
import { useAuthStore } from "@store/auth-store";

export function LogoutToolbarWidget() {
  const { t } = useAdminI18n();
  const logout = useAuthStore((s) => s.logout);

  return (
    <Tooltip title={t("admin.toolbar.logout.tooltip", "Đăng xuất")}>
      <Box
        aria-label={t("admin.toolbar.logout.tooltip", "Đăng xuất")}
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
