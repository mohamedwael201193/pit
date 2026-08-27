export function SiweBind({ connected }: { connected: boolean }) {
  if (!connected) {
    return (
      <p className="mt-4 max-w-[40ch] text-sm opacity-80">
        After you connect, you sign a bind message. No seed field exists.
      </p>
    );
  }
  return (
    <p className="mt-4 max-w-[40ch] text-sm opacity-80">
      Sign the bind message in your wallet to attach this address to your workspace.
    </p>
  );
}
