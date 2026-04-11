import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";

const mod: ModuleDescriptor = {
  id: "auth",
  routes: [
    {
      key: "auth",
      title: l("admin.auth.page_title"),
      subtitle: l("admin.auth.page_subtitle"),
      path: "/account",
      hidden: true,
    },
  ],
};

registerModule(mod);
