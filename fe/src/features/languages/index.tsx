import type { ModuleDescriptor } from "@root/core/module/types";
import { registerModule } from "@root/core/module/registry";
import TranslateIcon from "@mui/icons-material/Translate";
import OneColumnPage from "@root/core/pages/one-column-page";

const mod: ModuleDescriptor = {
  id: "languages",
  routes: [
    {
      key: "languages",
      permissions: ["languages.view"],
      label: "Ngôn ngữ",
      title: "Ngôn ngữ",
      subtitle: "Quản lý ngôn ngữ và resource i18n cho Admin Panel.",
      extra: {
        i18nTitleKey: "admin.languages.page_title",
        i18nSubtitleKey: "admin.languages.page_subtitle",
      },
      path: "/languages",
      element: <OneColumnPage />,
      icon: <TranslateIcon />,
      priority: 93,
      children: [
        {
          hidden: true,
          key: "languages-detail",
          permissions: ["languages.view"],
          label: "Chi tiết ngôn ngữ",
          title: "Chi tiết ngôn ngữ",
          subtitle: "Cập nhật metadata và resource admin.{module}.*.",
          extra: {
            i18nTitleKey: "admin.languages.detail_title",
            i18nSubtitleKey: "admin.languages.detail_subtitle",
          },
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
