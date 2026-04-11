import * as React from "react";
import { useShallow } from "zustand/react/shallow";
import { useAdminI18nStore } from "@store/admin-i18n-store";

export function useAdminI18n() {
  const state = useAdminI18nStore(
    useShallow((store) => ({
      languages: store.languages,
      currentLanguageCode: store.currentLanguageCode,
      resources: store.resources,
      isBootstrapping: store.isBootstrapping,
      isReady: store.isReady,
      bootstrap: store.bootstrap,
      changeLanguage: store.changeLanguage,
      clear: store.clear,
    }))
  );

  const t = React.useCallback(
    (key: string, fallback?: string) => {
      return state.resources[key] ?? fallback ?? key;
    },
    [state.resources]
  );

  const currentLanguage = React.useMemo(
    () => state.languages.find((item) => item.code === state.currentLanguageCode) ?? null,
    [state.currentLanguageCode, state.languages]
  );

  return {
    ...state,
    currentLanguage,
    t,
  };
}
