export function Avatar({
  online,
  size = 36,
  rounded = "full",
}: {
  online?: boolean;
  size?: number;
  rounded?: "full" | "9px";
}) {
  return (
    <div className="relative shrink-0" style={{ width: size, height: size }}>
      <div
        className={`h-full w-full bg-gradient-to-br from-accent to-down ${
          rounded === "full" ? "rounded-full" : "rounded-[9px]"
        }`}
      />
      {online !== undefined && (
        <span
          className={`absolute right-0 bottom-0 h-2.5 w-2.5 rounded-full border-2 border-surface ${
            online ? "bg-up" : "bg-text-faint"
          }`}
        />
      )}
    </div>
  );
}
