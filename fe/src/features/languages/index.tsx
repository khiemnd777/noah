import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import TranslateIcon from "@mui/icons-material/Translate";
import OneColumnPage from "@root/core/pages/one-column-page";

const mod: ModuleDescriptor = {
  id: "languages",
  routes: [
    {
      key: "languages",
      permissions: ["languages.view"],
      label: l("admin.languages.page_title"),
      title: l("admin.languages.page_title"),
      subtitle: l("admin.languages.page_subtitle"),
      path: "/languages",
      element: <OneColumnPage />,
      icon: <TranslateIcon />,
      priority: 93,
      children: [
        {
          hidden: true,
          key: "languages-detail",
          permissions: ["languages.view"],
          label: l("admin.languages.detail_title"),
          title: l("admin.languages.detail_title"),
          subtitle: l("admin.languages.detail_subtitle"),
          path: "/languages/:languageId",
          element: <OneColumnPage />,
          icon: <TranslateIcon />,
          priority: 94,
        },
      ],
    },
  ],
};

registerModule(mod);
