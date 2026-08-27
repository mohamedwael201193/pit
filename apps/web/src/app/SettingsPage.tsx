import { PageHead } from "../ui/PageHead";

export function SettingsPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Settings"
        lede="Honest limits. If a Galileo feature is not proven, it stays unavailable."
      />
      <ul className="mt-10 max-w-[48ch] divide-y divide-[rgb(240_231_212/0.2)] border-y border-[rgb(240_231_212/0.2)]">
        <li className="py-4">
          <p className="font-semibold">Web cannot authorize</p>
          <p className="mt-1 text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
            Orders are signed on desktop or CLI. This browser cannot hold a session.
          </p>
        </li>
        <li className="py-4">
          <p className="font-semibold">No seed field</p>
          <p className="mt-1 text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
            Connect a wallet you already control. PIT never asks for a private key.
          </p>
        </li>
        <li className="py-4">
          <p className="font-semibold">Networks do not mix</p>
          <p className="mt-1 text-[0.9375rem] text-[rgb(240_231_212/0.65)]">
            MAINNET and TESTNET are separate workspaces. Capabilities are not copied.
          </p>
        </li>
      </ul>
    </div>
  );
}
