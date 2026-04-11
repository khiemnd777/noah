import SearchOutlinedIcon from '@mui/icons-material/SearchOutlined';
import { Box, Tooltip } from "@mui/material";
import { useAdminI18n } from "@root/core/i18n/use-admin-i18n";
import { navigate } from "@root/core/navigation/navigate";

export default function SearchIcon() {
  const { t } = useAdminI18n();

  return (
    <Tooltip title={t("admin.toolbar.search.tooltip", "Tìm kiếm")}>
      <Box
        aria-label={t("admin.toolbar.search.tooltip", "Tìm kiếm")}
        onClick={() => navigate("/search")}
        sx={{
          position: "relative",
          display: "inline-flex",
          alignItems: "center",
          justifyContent: "center",
          cursor: "pointer",
        }}
      >
        <SearchOutlinedIcon />
      </Box>
    </Tooltip>
  );
}
