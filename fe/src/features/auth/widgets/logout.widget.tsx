import { useAuthStore } from "@store/auth-store";
import { LogoutRounded } from "@mui/icons-material";
import { SafeButton } from "@shared/components/button/safe-button";
import { registerSlot } from "@core/module/registry";
import { useI18n } from "@root/core/i18n/use-i18n";

function LogoutWidget() {
  const logout = useAuthStore((s) => s.logout);
  const { t } = useI18n();
  return (
    <>
      <SafeButton variant="contained" color="error" startIcon={<LogoutRounded />} onClick={async () => await logout()}>
        {t("admin.general.logout_button")}
      </SafeButton>
    </>
  );
}

registerSlot({
  id: "logout",
  name: "auth:actions",
  priority: 2,
  render: () => <LogoutWidget />,
});
