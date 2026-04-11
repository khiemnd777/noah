import { registerSlot } from "@core/module/registry";
import SearchBox from "@core/search/search-box";
import { navigate } from "@core/navigation/navigate";
import { useAdminI18n } from "@root/core/i18n/use-admin-i18n";
import { Box } from "@mui/material";
import type { SearchModel } from "@core/search/search.model";

export function SearchBoxWidget() {
  const { t } = useAdminI18n();

  const handleSelect = (_: SearchModel, href: string | void) => {
    if (typeof href === "string" && href.trim() !== "") {
      navigate(href);
    }
  };

  return (
    <>
      <Box>
        <SearchBox
          placeholder={t(
            "admin.search.placeholder",
            "Tìm kiếm theo tên sản phẩm, đơn hàng, vật tư, nhân sự và nha khoa..."
          )}
          onSelect={handleSelect}
          minChars={2}
          debounceMs={300}
          autoFocus
          fullWidth
        />
      </Box>
    </>
  );
}

registerSlot({
  id: "search",
  name: "search:left",
  render: () => <SearchBoxWidget />,
});
