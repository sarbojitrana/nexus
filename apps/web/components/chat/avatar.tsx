"use client";

import { RemoteAvatar } from "@/components/media/remote-image";

// Chat avatars were a placeholder gradient -- they now render the same picture
// as everywhere else, with the presence dot layered on top.
export function Avatar({
  online,
  size = 36,
  storageKey,
  url,
}: {
  online?: boolean;
  size?: number;
  storageKey?: string | null;
  url?: string | null;
}) {
  return (
    <div className="relative shrink-0" style={{ width: size, height: size }}>
      <RemoteAvatar storageKey={storageKey} url={url} size={size} />
      {online !== undefined && (
        <span
          className={`absolute right-0 bottom-0 h-2.5 w-2.5 border-2 border-surface ${
            online ? "bg-up" : "bg-text-faint"
          }`}
        />
      )}
    </div>
  );
}
