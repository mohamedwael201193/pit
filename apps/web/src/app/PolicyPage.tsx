import { PageHead } from "../ui/PageHead";
import { PolicyPanel } from "../PolicyPanel";
import { DiagramPolicy } from "../diagrams/pitGuide";

export function PolicyPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Policy"
        lede="Readable limits. The model cannot raise clip, leverage, or permissions."
      />
      <figure className="mt-10 max-w-[36rem] overflow-hidden border border-[rgb(240_231_212/0.28)]">
        <DiagramPolicy className="aspect-[16/10] w-full" />
      </figure>
      <PolicyPanel />
    </div>
  );
}
