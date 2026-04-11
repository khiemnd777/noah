import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import BadgeIcon from '@mui/icons-material/Badge';
import OneColumnPage from "@root/core/pages/one-column-page";

const mod: ModuleDescriptor = {
  id: "staff",
  routes: [
    {
      key: "staff",
      permissions: ["staff.view"],
      label: l("admin.staff.page_title"),
      title: l("admin.staff.page_title"),
      subtitle: l("admin.staff.page_subtitle"),
      path: "/staff",
      element: <OneColumnPage />,
      icon: <BadgeIcon />,
      priority: 94,
      children: [
        {
          hidden: true,
          key: "staff-detail",
          permissions: ["staff.view", "staff.update"],
          label: l("admin.staff.detail_title"),
          title: l("admin.staff.detail_title"),
          subtitle: l("admin.staff.detail_subtitle"),
          path: "/staff/:staffId",
          icon: <BadgeIcon />,
          element: <OneColumnPage />,
          priority: 99,
        },
      ],
    },
  ],
};

registerModule(mod);
