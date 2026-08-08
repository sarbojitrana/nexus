"use client";

import { useMediaUrls } from "@/lib/use-media-urls";

// Avatars and banners are stored as keys, so they need resolving before they
// can render. Both fall back to the accent block the app used before any
// images existed, which keeps layout stable while a URL is in flight.

export function RemoteAvatar({
  storageKey,
  url: directUrl,
  size = 40,
  className = "",
  onClick,
}: {
  storageKey?: string | null;
  /** Clerk's hosted avatar. Takes precedence -- Clerk's Manage Account is the
   *  primary place people change their picture, so its copy is the freshest. */
  url?: string | null;
  size?: number;
  className?: string;
  onClick?: () => void;
}) {
  // Skip the presign round trip entirely when Clerk already gave us a URL.
  const urls = useMediaUrls(directUrl ? [] : [storageKey]);
  const url = directUrl ?? (storageKey ? urls[storageKey] : undefined);

  return (
    <span
      className={`block shrink-0 overflow-hidden bg-accent ${onClick && url ? "cursor-zoom-in" : ""} ${className}`}
      style={{ width: size, height: size }}
      onClick={url ? onClick : undefined}
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
  onClick,
}: {
  storageKey?: string | null;
  className?: string;
  onClick?: () => void;
}) {
  const urls = useMediaUrls([storageKey]);
  const url = storageKey ? urls[storageKey] : undefined;

  return (
    <div
      className={`w-full overflow-hidden bg-gradient-to-r from-accent to-down ${onClick && url ? "cursor-zoom-in" : ""} ${className}`}
      onClick={url ? onClick : undefined}
    >
      {url && (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={url} alt="" className="h-full w-full object-cover" />
      )}
    </div>
  );
}
