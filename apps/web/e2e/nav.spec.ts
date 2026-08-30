export function assertVerifyNav(href: string) {
  if (href !== "/proof") {
    throw new Error("verify nav");
  }
}
