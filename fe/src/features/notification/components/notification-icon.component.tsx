import { useAsync } from "@root/core/hooks/use-async";
import { countUnread } from "@core/notification/notification.api";
import { NotifierChip } from "@root/shared/components/notification/notifier-chip";
import NotificationsIcon from "@mui/icons-material/Notifications";
import { Box, Tooltip } from "@mui/material";
import { useI18n } from "@root/core/i18n/use-i18n";
import { navigate } from "@root/core/navigation/navigate";

export default function NotificationIcon() {
  const { t } = useI18n();
  const { data: count } = useAsync<number>(() => countUnread(), [], {
    key: "notification-unread-count",
  });

  return (
    <Tooltip title={t("admin.toolbar.notifications.tooltip", "Thông báo")}>
      <Box
        aria-label={t("admin.toolbar.notifications.tooltip", "Thông báo")}
        onClick={() => navigate("/notification")}
        sx={{
          position: "relative",
          display: "inline-flex",
          alignItems: "center",
          justifyContent: "center",
          cursor: "pointer",
        }}
      >
        <NotificationsIcon />
        <Box
          sx={{
            position: "absolute",
            top: 0,
            right: 0,
            transform: "translate(50%, -50%)",
            display: "flex",
            alignItems: "center",
            justifyContent: "flex-start",
            pointerEvents: "none",
          }}
        >
          <NotifierChip count={count} />
        </Box>
      </Box>
    </Tooltip>
  );
}
