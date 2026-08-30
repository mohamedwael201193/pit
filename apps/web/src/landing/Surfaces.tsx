import { Cube, Desktop, Plugs, TerminalWindow } from "@phosphor-icons/react";
import { Reveal } from "../ui/Reveal";

const SURFACES = [
  {
    name: "Desktop",
    Icon: Desktop,
    tone: "coral" as const,
    body: "Session lives in the OS keychain. You type AUTHORIZE on the exact preview. Chat cannot.",
  },
  {
    name: "CLI",
    Icon: TerminalWindow,
    tone: "cream" as const,
    body: "pit.exe is the same host. TTY AUTHORIZE. It refuses to print the session key.",
  },
  {
    name: "MCP",
    Icon: Plugs,
    tone: "ink" as const,
    body: "Cursor can read markets and status. It cannot order, cancel, pin, or enable autonomy.",
  },
  {
    name: "SDK",
    Icon: Cube,
    tone: "line" as const,
    body: "0G Direct TeeML seals the private book. This website never holds a Hyperliquid session.",
  },
] as const;

const TONE: Record<(typeof SURFACES)[number]["tone"], string> = {
  coral: "bg-[#d82f2f] text-black",
  cream: "bg-[#f0e7d4] text-black",
  ink: "bg-[#141414] text-[var(--guide-cream)]",
  line: "bg-[#1a1a1a] text-[var(--guide-cream)]",
};

export function Surfaces() {
  return (
    <section className="border-t border-[rgb(240_231_212/0.25)] py-20 md:py-28">
      <div className="container-pit">
        <Reveal>
          <h2 className="guide-display max-w-[12ch]">Four doors. One seat.</h2>
          <p className="mt-6 max-w-[40ch] text-[1.2rem] leading-8 text-[rgb(240_231_212/0.78)]">
            Research can live on the web. Execution stays on this computer.
          </p>
        </Reveal>
        <div className="mt-12 grid grid-cols-1 gap-px bg-[rgb(240_231_212/0.28)] md:grid-cols-2 md:grid-flow-dense">
          {SURFACES.map(({ name, Icon, tone, body }, i) => (
            <Reveal
              key={name}
              index={i}
              as="article"
              className={`flex min-h-[15rem] flex-col justify-between p-8 md:min-h-[17rem] md:p-10 ${TONE[tone]} ${name === "Desktop" || name === "SDK" ? "md:col-span-2" : ""}`}
            >
              <Icon size={32} weight="regular" aria-hidden="true" />
              <div className="mt-10">
                <h3 className="text-[2rem] font-bold tracking-[-0.04em]">{name}</h3>
                <p className="mt-3 max-w-[42ch] text-[1.05rem] leading-7 opacity-80">{body}</p>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  );
}
