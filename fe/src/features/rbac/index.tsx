import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import KeyIcon from '@mui/icons-material/Key';
import OneColumnPage from "@root/core/pages/one-column-page";

const mod: ModuleDescriptor = {
  id: "rbac",
  routes: [
    {
      key: "rbac",
      permissions: ["rbac.manage"],
      label: l("admin.rbac.page_title"),
      title: l("admin.rbac.page_title"),
      subtitle: l("admin.rbac.page_subtitle"),
      path: "/rbac",
      element: <OneColumnPage />,
      icon: <KeyIcon />,
      priority: 1,
    },
  ],
};

registerModule(mod);
