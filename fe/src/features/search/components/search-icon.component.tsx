import SearchOutlinedIcon from '@mui/icons-material/SearchOutlined';
import { Box, Tooltip } from "@mui/material";
import { useI18n } from "@root/core/i18n/use-i18n";
import { navigate } from "@root/core/navigation/navigate";

export default function SearchIcon() {
  const { t } = useI18n();
  const tooltip = t("admin.toolbar.search.tooltip");

  return (
    <Tooltip title={tooltip}>
      <Box
        aria-label={tooltip}
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
