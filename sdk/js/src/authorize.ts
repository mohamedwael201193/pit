export function refuseAuthorize(): never {
  throw new Error("authorize_denied");
}
