import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import NotificationsIcon from '@mui/icons-material/Notifications';
import { NotificationChip } from "./components/notification-chip.component";

const mod: ModuleDescriptor = {
  id: "notification",
  routes: [
    {
      key: "notification",
      label: l("admin.notification.page_title"),
      title: l("admin.notification.page_title"),
      path: "/notification",
      icon: <NotificationsIcon />,
      chip: <NotificationChip />,
      hidden: true,
      priority: 9998,
    },
  ],
};

registerModule(mod);
