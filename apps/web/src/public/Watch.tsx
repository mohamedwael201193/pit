import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";
import { fetchHealth, fetchWatch, UnsafeWatchError } from "./api";
import { probeDesktop } from "./desktop";
import type { DesktopProbe, HealthView, PublicCoin, WatchView } from "./types";

type WatchState = {
  watch: WatchView | null;
  health: HealthView | null;
  error: string | null;
  loading: boolean;
  desktop: DesktopProbe;
};

const WatchCtx = createContext<WatchState | null>(null);

export function WatchProvider({ children }: { children: ReactNode }) {
  const [watch, setWatch] = useState<WatchView | null>(null);
  const [health, setHealth] = useState<HealthView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [desktop, setDesktop] = useState<DesktopProbe>({ present: false, refused: false });

  useEffect(() => {
    const ac = new AbortController();
    setLoading(true);
    Promise.all([fetchWatch("mainnet", ac.signal), fetchHealth(ac.signal)])
      .then(([w, h]) => {
        setWatch(w);
        setHealth(h);
        setError(null);
      })
      .catch((err: unknown) => {
        if (ac.signal.aborted) return;
        setWatch(null);
        setHealth(null);
        setError(err instanceof UnsafeWatchError ? err.message : "Public health is unreachable. PIT will not invent live counts.");
      })
      .finally(() => {
        if (!ac.signal.aborted) setLoading(false);
      });
    return () => ac.abort();
  }, []);

  useEffect(() => {
    const ac = new AbortController();
    void probeDesktop(ac.signal).then(setDesktop);
    const t = window.setInterval(() => {
      void probeDesktop().then(setDesktop);
    }, 8000);
    return () => {
      ac.abort();
      window.clearInterval(t);
    };
  }, []);

  const value = useMemo(() => ({ watch, health, error, loading, desktop }), [watch, health, error, loading, desktop]);
  return <WatchCtx.Provider value={value}>{children}</WatchCtx.Provider>;
}

export function useWatch(): WatchState {
  const ctx = useContext(WatchCtx);
  if (!ctx) throw new Error("useWatch");
  return ctx;
}

export function useWatchSafe(): WatchState | null {
  return useContext(WatchCtx);
}

export function eligibleCoins(watch: WatchView | null): PublicCoin[] {
  return (watch?.coins ?? []).filter((c) => c.eligible);
}

export function blockedCoins(watch: WatchView | null): PublicCoin[] {
  return (watch?.coins ?? []).filter((c) => Boolean(c.block));
}

export function watchCoins(watch: WatchView | null): PublicCoin[] {
  return (watch?.coins ?? []).filter((c) => !c.eligible).slice(0, 24);
}

export function actionableCoins(watch: WatchView | null): PublicCoin[] {
  return (watch?.coins ?? []).filter((c) => c.executionFeasible);
}

export function findCoin(watch: WatchView | null, id: string): PublicCoin | undefined {
  const want = id.trim().toUpperCase();
  return (watch?.coins ?? []).find((c) => c.coin.toUpperCase() === want);
}
