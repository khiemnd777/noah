import * as React from "react";
import { subscribeAuthEvents } from "@core/network/auth-session";
import { useAdminI18nStore } from "@store/admin-i18n-store";
import { useAuthStore } from "@store/auth-store";

export function AdminI18nBootstrap() {
  const userId = useAuthStore((state) => state.user?.id ?? null);
  const isLoggedIn = useAuthStore((state) => state.isLoggedIn);
  const bootstrap = useAdminI18nStore((state) => state.bootstrap);
  const clear = useAdminI18nStore((state) => state.clear);

  React.useEffect(() => {
    if (!isLoggedIn || !userId) {
      return;
    }

    void bootstrap();
  }, [bootstrap, isLoggedIn, userId]);

  React.useEffect(() => {
    const unsubscribe = subscribeAuthEvents((event) => {
      if (event.type === "logout" || event.type === "refresh_failed") {
        clear();
      }
    });
    return () => {
      unsubscribe();
    };
  }, [clear]);

  return null;
}
