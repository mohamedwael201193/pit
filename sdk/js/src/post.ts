export const canPostExchange = false;

export function refuseUnsignedPost(): never {
  throw new Error("exchange_unsigned");
}
