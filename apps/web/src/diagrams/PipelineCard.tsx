import { DiagramHeroPostcard } from "./pitGuide";

export { PIPELINE } from "./pipeline";

export function PipelineCard({ className = "" }: { className?: string }) {
  return (
    <figure className={`border border-black ${className}`}>
      <DiagramHeroPostcard className="aspect-[4/3] w-full" />
    </figure>
  );
}
