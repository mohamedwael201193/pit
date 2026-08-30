import { useEffect } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { WatchProvider } from "../public/Watch";
import { DeskFrame } from "../desk/DeskFrame";
import { WalletGate } from "./WalletGate";

export function AppShell() {
  const { ready, authenticated } = usePrivy();
  const location = useLocation();

  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", "#D82F2F");
  }, []);

  if (location.pathname === "/app/start" || location.pathname.startsWith("/app/start/")) {
    return <Navigate to="/protect" replace />;
  }

  if (!ready) {
    return (
      <div className="guide-shell grid min-h-[100dvh] place-items-center">
        <p>Loading wallet connect</p>
      </div>
    );
  }
  if (!authenticated) return <WalletGate />;

  return (
    <WatchProvider>
      <DeskFrame />
    </WatchProvider>
  );
}
