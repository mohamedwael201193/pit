import { Navigate, Route, Routes } from "react-router-dom";
import { Landing } from "./landing/Landing";
import { WalletGate } from "./app/WalletGate";
import { AppShell } from "./app/AppShell";
import { Home } from "./app/Home";
import { StartFlow } from "./app/StartFlow";
import { Activity } from "./app/Activity";
import { PolicyPage } from "./app/PolicyPage";
import { AccountPage } from "./app/AccountPage";
import { SettingsPage } from "./app/SettingsPage";
import { VerifyPage } from "./app/VerifyPage";

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/signin" element={<WalletGate />} />
      <Route path="/app" element={<AppShell />}>
        <Route index element={<Home />} />
        <Route path="start" element={<StartFlow />} />
        <Route path="activity" element={<Activity />} />
        <Route path="policy" element={<PolicyPage />} />
        <Route path="account" element={<AccountPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="verify" element={<VerifyPage />} />
      </Route>
      <Route path="/verify" element={<AppShell />}>
        <Route index element={<VerifyPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
