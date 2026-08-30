import { useEffect } from "react";
import { WatchProvider } from "./Watch";
import { DeskFrame } from "../desk/DeskFrame";

export function PublicShell() {
  useEffect(() => {
    document.documentElement.dataset.theme = "guide";
  }, []);

  return (
    <WatchProvider>
      <DeskFrame />
    </WatchProvider>
  );
}
