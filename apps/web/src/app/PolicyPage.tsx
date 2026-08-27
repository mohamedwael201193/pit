import { PageHead } from "../ui/PageHead";
import { PolicyPanel } from "../PolicyPanel";

export function PolicyPage() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Policy"
        lede="Readable limits. The model cannot raise clip, leverage, or permissions."
      />
      <PolicyPanel />
    </div>
  );
}
