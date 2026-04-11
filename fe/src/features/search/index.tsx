import type { ModuleDescriptor } from "@root/core/module/types";
import { l } from "@root/core/i18n/localized-text";
import { registerModule } from "@root/core/module/registry";
import SearchOutlinedIcon from '@mui/icons-material/SearchOutlined';

const mod: ModuleDescriptor = {
  id: "search",
  routes: [
    {
      key: "search",
      label: l("admin.search.page_title"),
      title: l("admin.search.page_title"),
      path: "/search",
      icon: <SearchOutlinedIcon />,
      hidden: true,
      priority: 9999,
    },
  ],
};

registerModule(mod);
