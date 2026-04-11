/* eslint-disable react-refresh/only-export-components */
import * as React from "react";
import TranslateRoundedIcon from "@mui/icons-material/TranslateRounded";
import {
  Box,
  CircularProgress,
  MenuItem,
  Select,
  type SelectChangeEvent,
} from "@mui/material";
import { useAdminI18n } from "@root/core/i18n/use-admin-i18n";
import { registerSlot } from "@root/core/module/registry";

function LanguageToolbarWidget() {
  const {
    changeLanguage,
    currentLanguageCode,
    isBootstrapping,
    languages,
    t,
  } = useAdminI18n();
  const [isSaving, setIsSaving] = React.useState(false);

  const handleChange = React.useCallback(
    async (event: SelectChangeEvent<string>) => {
      const nextCode = event.target.value;
      setIsSaving(true);
      try {
        await changeLanguage(nextCode);
      } finally {
        setIsSaving(false);
      }
    },
    [changeLanguage]
  );

  const resolvedValue = React.useMemo(() => {
    if (currentLanguageCode) return currentLanguageCode;
    return languages[0]?.code ?? "";
  }, [currentLanguageCode, languages]);

  return (
      <Box
        sx={{
          minWidth: 148,
          display: "inline-flex",
          alignItems: "center",
          gap: 1,
        }}
      >
        {isSaving || isBootstrapping ? (
          <CircularProgress size={16} />
        ) : (
          <TranslateRoundedIcon fontSize="small" />
        )}

        <Select
          size="small"
          value={resolvedValue}
          displayEmpty
          disabled={languages.length === 0 || isSaving}
          onChange={handleChange}
          variant="standard"
          inputProps={{
            "aria-label": t("admin.toolbar.language.tooltip", "Ngôn ngữ giao diện"),
          }}
          sx={{
            minWidth: 112,
            fontSize: 14,
          }}
        >
          {languages.length === 0 ? (
            <MenuItem value="">
              {t("admin.toolbar.language.empty", "Chưa có ngôn ngữ")}
            </MenuItem>
          ) : (
            languages.map((language) => (
              <MenuItem key={language.code} value={language.code}>
                {language.nativeName || language.name || language.code}
              </MenuItem>
            ))
          )}
        </Select>
      </Box>
  );
}

registerSlot({
  id: "language-toolbar",
  name: "toolbar",
  priority: 97.5,
  render: () => <LanguageToolbarWidget />,
});
