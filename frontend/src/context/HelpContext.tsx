import { createContext, ReactNode, useCallback, useContext, useMemo, useState } from 'react';

interface HelpContextValue {
  drawerOpen: boolean;
  openDrawer: () => void;
  closeDrawer: () => void;
  openArticleInHub: (articleId: string) => void;
  pendingArticleId: string | null;
  clearPendingArticle: () => void;
}

const HelpContext = createContext<HelpContextValue | null>(null);

export function HelpProvider({ children }: { children: ReactNode }) {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [pendingArticleId, setPendingArticleId] = useState<string | null>(null);

  const openDrawer = useCallback(() => setDrawerOpen(true), []);
  const closeDrawer = useCallback(() => setDrawerOpen(false), []);

  const openArticleInHub = useCallback((articleId: string) => {
    setPendingArticleId(articleId);
    setDrawerOpen(false);
  }, []);

  const clearPendingArticle = useCallback(() => setPendingArticleId(null), []);

  const value = useMemo(
    () => ({
      drawerOpen,
      openDrawer,
      closeDrawer,
      openArticleInHub,
      pendingArticleId,
      clearPendingArticle,
    }),
    [drawerOpen, openDrawer, closeDrawer, openArticleInHub, pendingArticleId, clearPendingArticle],
  );

  return <HelpContext.Provider value={value}>{children}</HelpContext.Provider>;
}

export function useHelp() {
  const context = useContext(HelpContext);
  if (!context) {
    throw new Error('useHelp must be used within HelpProvider');
  }
  return context;
}
