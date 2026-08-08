"use client";

import { RemoteAvatar, RemoteBanner } from "@/components/media/remote-image";
import { useLightbox } from "@/components/media/lightbox";
import { useMediaUrls } from "@/lib/use-media-urls";

// Wraps the avatar/banner so a click opens it full screen. The URL has to be
// resolved here too, since only the resolved URL can be handed to the viewer.

export function ZoomableAvatar({
  storageKey,
  url,
  size = 40,
  className = "",
}: {
  storageKey?: string | null;
  url?: string | null;
  size?: number;
  className?: string;
}) {
  const open = useLightbox();
  const urls = useMediaUrls(url ? [] : [storageKey]);
  const resolved = url ?? (storageKey ? urls[storageKey] : undefined);

  return (
    <RemoteAvatar
      storageKey={storageKey}
      url={url}
      size={size}
      className={className}
      onClick={resolved ? () => open({ url: resolved, kind: "image" }) : undefined}
    />
  );
}

export function ZoomableBanner({
  storageKey,
  className = "",
}: {
  storageKey?: string | null;
  className?: string;
}) {
  const open = useLightbox();
  const urls = useMediaUrls([storageKey]);
  const resolved = storageKey ? urls[storageKey] : undefined;

  return (
    <RemoteBanner
      storageKey={storageKey}
      className={className}
      onClick={resolved ? () => open({ url: resolved, kind: "image" }) : undefined}
    />
  );
}
