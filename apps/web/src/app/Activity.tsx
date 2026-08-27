import { PageHead } from "../ui/PageHead";
import { Bezel } from "../ui/Surface";
import { ProgressStrip } from "../ProgressStrip";
import { NoSession } from "../NoSession";
import { DiagramAuthorize } from "../diagrams/pitGuide";

export function Activity() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Activity"
        lede="Named states from the real pipeline. Timers never fake progress."
      />
      <figure className="mt-10 max-w-[36rem] overflow-hidden border border-[rgb(240_231_212/0.28)]">
        <DiagramAuthorize className="aspect-[16/10] w-full" />
      </figure>
      <Bezel className="mt-8 max-w-[36rem]">
        <ProgressStrip current="WAITING_FOR_USER" />
        <NoSession />
      </Bezel>
    </div>
  );
}
