export function Logo({
  className = "",
  stacked = true,
}: {
  className?: string;
  stacked?: boolean;
}) {
  if (!stacked) {
    return (
      <span
        className={`inline-flex items-center gap-2 font-mono text-[0.95rem] font-bold tracking-[0.02em] ${className}`}
      >
        <span className="h-2 w-2 bg-accent" />
        <span>
          NEXUS<span className="text-accent">NETWORK</span>
        </span>
      </span>
    );
  }

  return (
    <span className={`inline-flex flex-col font-mono leading-[1.05] font-bold ${className}`}>
      <span className="flex items-center gap-2">
        <span className="h-2 w-2 shrink-0 bg-accent" />
        <span className="text-[0.92rem] tracking-[0.02em] text-text">NEXUS</span>
      </span>
      <span className="ml-4 text-[0.92rem] tracking-[0.02em] text-accent">NETWORK</span>
    </span>
  );
}

export function HudCorners() {
  return (
    <>
      <span className="hud-corner hud-corner--tl" />
      <span className="hud-corner hud-corner--tr" />
      <span className="hud-corner hud-corner--bl" />
      <span className="hud-corner hud-corner--br" />
    </>
  );
}
