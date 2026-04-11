import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import SettingsRoundedIcon from "@mui/icons-material/SettingsRounded";

const mod: ModuleDescriptor = {
  id: "settings",
  routes: [
    {
      key: "settings",
      permissions: ["settings.view"],
      label: l("admin.settings.page_title"),
      title: l("admin.settings.page_title"),
      subtitle: l("admin.settings.page_subtitle"),
      path: "/settings",
      icon: <SettingsRoundedIcon />,
      hidden: true,
      priority: 0,
    },
  ],
};

registerModule(mod);
