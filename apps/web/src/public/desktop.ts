import { COMPANION } from "./facts";
import type { DesktopProbe } from "./types";

export async function probeDesktop(signal?: AbortSignal): Promise<DesktopProbe> {
  try {
    const r = await fetch(`${COMPANION}/health`, { signal });
    const body = (await r.json()) as { ok?: boolean; sign?: boolean; version?: string };
    if (body.sign) return { present: false, refused: true };
    if (body.ok) return { present: true, version: body.version, refused: false };
  } catch {
    /* mixed-content on HTTPS production, or companion down */
  }
  return { present: false, refused: false };
}
