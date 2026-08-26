export function confirmAuthorize(typed: string, sessionAlive: boolean): string | null {
  if (!sessionAlive) {
    return "session_expired";
  }
  if (typed.trim() !== "AUTHORIZE") {
    return "authorize_refused";
  }
  return null;
}
