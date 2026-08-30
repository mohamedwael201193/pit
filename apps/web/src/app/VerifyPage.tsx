import { Link } from "react-router-dom";
import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { ButtonLink } from "../ui/Button";

export function VerifyPage() {
  return (
    <div className="mx-auto max-w-[40rem]">
      <PageHead
        title="My Proof"
        lede="Verification runs against Aristotle in the public proof center. This page does not switch to Galileo."
      />
      <Bezel className="mt-8">
        <p className="max-w-[48ch] text-[0.9375rem] leading-6 text-[rgb(240_231_212/0.7)]">
          Session keys never live here. This browser does not run VerifyE2EE. Paste a 0G Chain hash on the public
          proof page and read the receipt yourself.
        </p>
        <div className="mt-6">
          <ButtonLink as={Link} to="/proof" size="lg">
            Open proof center
          </ButtonLink>
        </div>
      </Bezel>
    </div>
  );
}
