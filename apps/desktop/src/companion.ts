export const COMPANION = "http://127.0.0.1:17373";

export type LocalStatus = {
  sessionAlive?: boolean;
  agent?: string;
  network?: string;
  sign?: boolean;
  trade?: boolean;
};

export async function localStatus(): Promise<LocalStatus | null> {
  try {
    const r = await fetch(`${COMPANION}/local/status`);
    if (!r.ok) return null;
    const body = (await r.json()) as LocalStatus;
    if (body.sign || body.trade) return null;
    return body;
  } catch {
    return null;
  }
}
