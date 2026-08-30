import { Link } from "react-router-dom";
import { usePrivy } from "@privy-io/react-auth";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { IsolateNote } from "../IsolateNote";
import { KillNote } from "../KillNote";
import { TransferNote } from "../TransferNote";
import { NoSession } from "../NoSession";
import { ButtonLink } from "../ui/Button";

export function AccountPage() {
  const { user } = usePrivy();
  const addr = user?.wallet?.address;

  return (
    <div className="mx-auto flex max-w-[80rem] flex-col gap-8">
      <PageHead
        title="My Agent"
        lede="This address is login identity. Session keys never live here. Bind and mint stay on desktop."
      />

      <Bezel>
        {addr ? <p className="font-mono text-[0.9375rem] break-all text-[var(--guide-cream)]">{addr}</p> : null}
        <p className="mt-4 max-w-[48ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.7)]">
          This page does not collect a SIWE or a seed. Pairing and Protect my strategy run on PIT Desktop after you have
          seen radar and proof.
        </p>
        <NoSession />
        <div className="mt-6 flex flex-wrap gap-2">
          <ButtonLink as={Link} to="/agent" size="lg">
            Public passport
          </ButtonLink>
          <ButtonLink as={Link} to="/proof" variant="secondary" size="lg">
            Proof
          </ButtonLink>
        </div>
      </Bezel>

      <Bezel>
        <IsolateNote />
        <KillNote />
        <TransferNote />
      </Bezel>
    </div>
  );
}
