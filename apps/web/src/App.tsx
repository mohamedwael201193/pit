import { Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "./app/AppShell";
import { Home } from "./app/Home";
import { Activity } from "./app/Activity";
import { PolicyPage } from "./app/PolicyPage";
import { AccountPage } from "./app/AccountPage";
import { SettingsPage } from "./app/SettingsPage";
import { VerifyPage } from "./app/VerifyPage";
import { PairPage } from "./PairPage";
import { PublicShell } from "./public/Shell";
import { Landing } from "./landing/Landing";
import { RadarPage } from "./public/Radar";
import { MarketPage } from "./public/Market";
import { CapitalPage } from "./public/Capital";
import { MissionsPage } from "./public/Missions";
import { ReplayPage } from "./public/Replay";
import { AutonomyPage } from "./public/Autonomy";
import { MissionDetailPage } from "./public/MissionDetail";
import { ProofPage } from "./public/Proof";
import { AgentPage } from "./public/Agent";
import { HowPage } from "./public/How";
import { DownloadPage } from "./public/Download";
import { ProtectPage } from "./public/Protect";

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route element={<PublicShell />}>
        <Route path="/radar" element={<RadarPage />} />
        <Route path="/radar/:coin" element={<MarketPage />} />
        <Route path="/capital" element={<CapitalPage />} />
        <Route path="/autonomy" element={<AutonomyPage />} />
        <Route path="/missions" element={<MissionsPage />} />
        <Route path="/missions/:id" element={<MissionDetailPage />} />
        <Route path="/missions/:id/replay" element={<ReplayPage />} />
        <Route path="/proof" element={<ProofPage />} />
        <Route path="/agent" element={<AgentPage />} />
        <Route path="/how-it-works" element={<HowPage />} />
        <Route path="/download" element={<DownloadPage />} />
        <Route path="/pair" element={<PairPage />} />
        <Route path="/protect" element={<ProtectPage />} />
        <Route path="/watch" element={<Navigate to="/radar" replace />} />
      </Route>
      <Route path="/signin" element={<Navigate to="/protect" replace />} />
      <Route path="/app" element={<AppShell />}>
        <Route index element={<Home />} />
        <Route path="start" element={<Navigate to="/protect" replace />} />
        <Route path="activity" element={<Activity />} />
        <Route path="policy" element={<PolicyPage />} />
        <Route path="account" element={<AccountPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="verify" element={<VerifyPage />} />
      </Route>
      <Route path="/verify" element={<Navigate to="/proof" replace />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
