import { PageHead } from "../ui/PageHead";
import { ProgressStrip } from "../ProgressStrip";
import { NoSession } from "../NoSession";

export function Activity() {
  return (
    <div className="mx-auto max-w-[80rem]">
      <PageHead
        title="Activity"
        lede="Named states from the real pipeline. Timers never fake progress."
      />
      <ProgressStrip current="WAITING_FOR_USER" />
      <NoSession />
    </div>
  );
}
