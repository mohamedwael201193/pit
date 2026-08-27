import { Navigate, Route, Routes } from "react-router-dom";
import { Landing } from "./landing/Landing";
import { WalletGate } from "./app/WalletGate";
import { AppShell } from "./app/AppShell";
import { Home } from "./app/Home";
import { StartFlow } from "./app/StartFlow";
import { VerifyPage } from "./app/VerifyPage";

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/signin" element={<WalletGate />} />
      <Route path="/app" element={<AppShell />}>
        <Route index element={<Home />} />
        <Route path="start" element={<StartFlow />} />
      </Route>
      <Route path="/verify" element={<AppShell />}>
        <Route index element={<VerifyPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
