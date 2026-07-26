export function Logo({ className = "" }: { className?: string }) {
  return (
    <span
      className={`inline-flex items-center gap-2 font-display text-[1.15rem] font-extrabold tracking-tight ${className}`}
    >
      <span className="h-[9px] w-[9px] rounded-full bg-accent shadow-[0_0_0_4px_color-mix(in_srgb,var(--color-accent)_13%,transparent)]" />
      nexus
    </span>
  );
}
