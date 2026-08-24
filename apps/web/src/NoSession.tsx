export function NoSession() {
  return (
    <p className="mt-4 max-w-[48ch] text-sm opacity-80">
      This browser never holds a Hyperliquid session. Disconnect does not export a key.
    </p>
  );
}
