"use client";

import { useMediaUrls } from "@/lib/use-media-urls";

// Avatars and banners are stored as keys, so they need resolving before they
// can render. Both fall back to the accent block the app used before any
// images existed, which keeps layout stable while a URL is in flight.

export function RemoteAvatar({
  storageKey,
  size = 40,
  className = "",
}: {
  storageKey?: string | null;
  size?: number;
  className?: string;
}) {
  const urls = useMediaUrls([storageKey]);
  const url = storageKey ? urls[storageKey] : undefined;

  return (
    <span
      className={`block shrink-0 overflow-hidden bg-accent ${className}`}
      style={{ width: size, height: size }}
    >
      {url && (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={url} alt="" className="h-full w-full object-cover" />
      )}
    </span>
  );
}

export function RemoteBanner({
  storageKey,
  className = "",
}: {
  storageKey?: string | null;
  className?: string;
}) {
  const urls = useMediaUrls([storageKey]);
  const url = storageKey ? urls[storageKey] : undefined;

  return (
    <div
      className={`w-full overflow-hidden bg-gradient-to-r from-accent to-down ${className}`}
    >
      {url && (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={url} alt="" className="h-full w-full object-cover" />
      )}
    </div>
  );
}
