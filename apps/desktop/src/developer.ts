const KEY = "pit.developer";

export function isDeveloper(): boolean {
  try {
    if (typeof localStorage === "undefined") return false;
    if (localStorage.getItem(KEY) === "1") return true;
    if (new URLSearchParams(window.location.search).get("dev") === "1") {
      localStorage.setItem(KEY, "1");
      return true;
    }
  } catch {
    return false;
  }
  return false;
}

export function armDeveloper(): boolean {
  try {
    const n = Number(sessionStorage.getItem("pit.dev.clicks") || "0") + 1;
    sessionStorage.setItem("pit.dev.clicks", String(n));
    if (n >= 7) {
      localStorage.setItem(KEY, "1");
      return true;
    }
  } catch {
    return false;
  }
  return isDeveloper();
}
