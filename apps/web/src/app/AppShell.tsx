import { useEffect } from "react";
import { usePrivy } from "@privy-io/react-auth";
import { WatchProvider } from "../public/Watch";
import { DeskFrame } from "../desk/DeskFrame";
import { WalletGate } from "./WalletGate";

export function AppShell() {
  const { ready, authenticated } = usePrivy();

  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", "#D82F2F");
  }, []);

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
