import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";

const RULES = [
  {
    title: "Web cannot authorize",
    body: "Orders are signed on desktop or CLI. This browser cannot hold a session.",
  },
  {
    title: "No seed field",
    body: "Connect a wallet you already control. PIT never asks for a private key.",
  },
  {
    title: "Networks do not mix",
    body: "MAINNET and TESTNET are separate workspaces. Capabilities are not copied.",
  },
] as const;

export function SettingsPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Settings"
        lede="Honest limits. If a Galileo feature is not proven, it stays unavailable."
      />
      <div className="mt-10 grid gap-4 md:grid-cols-2">
        {RULES.map((rule) => (
          <Bezel key={rule.title}>
            <p className="font-semibold text-[var(--guide-cream)]">{rule.title}</p>
            <p className="mt-2 text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.65)]">{rule.body}</p>
          </Bezel>
        ))}
      </div>
    </div>
  );
}
