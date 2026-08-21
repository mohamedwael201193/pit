import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { PrivyProvider } from "@privy-io/react-auth";
import { App } from "./App";
import "./index.css";

const appId = import.meta.env.VITE_PRIVY_APP_ID as string;

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <PrivyProvider
      appId={appId}
      config={{
        appearance: { theme: "dark", accentColor: "#D82F2F" },
        embeddedWallets: { createOnLogin: "off" },
        supportedChains: [
          {
            id: 16661,
            name: "0G Aristotle",
            nativeCurrency: { name: "0G", symbol: "0G", decimals: 18 },
            rpcUrls: { default: { http: ["https://evmrpc.0g.ai"] } },
            blockExplorers: { default: { name: "ChainScan", url: "https://chainscan.0g.ai" } },
          },
          {
            id: 16602,
            name: "0G Galileo",
            nativeCurrency: { name: "0G", symbol: "0G", decimals: 18 },
            rpcUrls: { default: { http: ["https://evmrpc-testnet.0g.ai"] } },
            blockExplorers: { default: { name: "ChainScan Galileo", url: "https://chainscan-galileo.0g.ai" } },
          },
        ],
      }}
    >
      <App />
    </PrivyProvider>
  </StrictMode>,
);
