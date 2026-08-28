import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { ButtonLink } from "../ui/Button";
import { NoSession } from "../NoSession";
import { DiagramAuthorize } from "../diagrams/pitGuide";

export function Activity() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Activity"
        lede="This browser cannot hold a session or invent an OID. Exact orders live on PIT Desktop."
      />
      <figure className="mt-10 max-w-[36rem] overflow-hidden border border-[rgb(240_231_212/0.28)]">
        <DiagramAuthorize className="aspect-[16/10] w-full" />
      </figure>
      <Bezel className="mt-8 max-w-[36rem]">
        <p className="text-[0.75rem] tracking-[0.14em] text-[#d82f2f]">STATUS</p>
        <p className="mt-2 text-[1.25rem] font-semibold">Execution lives on PIT Desktop</p>
        <p className="mt-3 max-w-[48ch] text-[0.975rem] leading-6 text-[rgb(240_231_212/0.75)]">
          WHY: this site never receives a session key. NEXT: open Activity in PIT Desktop after you type AUTHORIZE.
        </p>
        <div className="mt-6 flex flex-wrap gap-3">
          <ButtonLink href="https://app.hyperliquid.xyz/portfolio" target="_blank" rel="noreferrer" size="lg">
            Open Hyperliquid portfolio
          </ButtonLink>
          <ButtonLink href="https://app.hyperliquid.xyz/API" target="_blank" rel="noreferrer" variant="secondary" size="lg">
            Open Hyperliquid API
          </ButtonLink>
        </div>
        <NoSession />
      </Bezel>
    </div>
  );
}
