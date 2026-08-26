const CARDS = [
  {
    title: "YOUR WALLET",
    body: "Connect. PIT never asks for a seed phrase. The browser never holds a session.",
  },
  {
    title: "YOUR SESSION",
    body: "Authorize the exact preview on this machine. Order and cancel only. Walk away if the card is wrong.",
  },
] as const;

export function StartCards() {
  return (
    <div className="start">
      {CARDS.map((c) => (
        <article key={c.title}>
          <p className="label">{c.title}</p>
          <p>{c.body}</p>
        </article>
      ))}
    </div>
  );
}
