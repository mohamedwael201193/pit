export function askNotify() {
  if (!("Notification" in window)) return;
  if (Notification.permission === "default") void Notification.requestPermission();
}

export function deskNotify(title: string, body: string) {
  if (!("Notification" in window)) return;
  if (Notification.permission !== "granted") return;
  if (document.hasFocus()) return;
  try {
    new Notification(title, { body });
  } catch {
    /* webview may refuse */
  }
}
