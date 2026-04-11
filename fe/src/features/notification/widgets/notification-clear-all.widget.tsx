import { registerSlot } from "@root/core/module/registry";
import { useAsync } from "@root/core/hooks/use-async";
import { useI18n } from "@root/core/i18n/use-i18n";
import { SafeButton } from "@shared/components/button/safe-button";
import DeleteIcon from "@mui/icons-material/Delete";
import { deleteAllNotifications, shortList } from "@core/notification/notification.api";

function NotificationClearAllWidget() {
  const { t } = useI18n();
  const { data } = useAsync(
    () => shortList(),
    [],
    { key: "notification-list-for-clear-all" }
  );
  const hasNotifications = (data ?? []).length > 0;

  return (
    <SafeButton
      variant="contained"
      color="error"
      startIcon={<DeleteIcon />}
      disabled={!hasNotifications}
      onClick={async () => {
        await deleteAllNotifications();
      }}
    >
      {t("admin.notification.clear_all_button")}
    </SafeButton>
  );
}

registerSlot({
  id: "notification-clear-all",
  name: "notification:actions",
  render: () => <NotificationClearAllWidget />,
});
