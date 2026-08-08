"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import { VideoPlayer } from "@/components/media/video-player";
import { useLightbox } from "@/components/media/lightbox";
import type { PostMedia } from "@nexus/zod";

// Storage keys aren't URLs -- they're resolved to short-lived signed URLs on
// demand, batched into one request per gallery.
export function MediaGallery({ media }: { media: PostMedia[] }) {
  const api = useApi();
  const openLightbox = useLightbox();
  const [urls, setUrls] = useState<Record<string, string>>({});

  const keys = media.map((m) => m.downloadKey).join(",");

  useEffect(() => {
    if (media.length === 0) return;
    let cancelled = false;
    api.Storage.presignDownloads({ body: { keys: keys.split(",") } })
      .then((res) => {
        if (!cancelled && res.status === 200) setUrls(res.body.urls);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [keys, api]);

  if (media.length === 0) return null;

  return (
    <div className={`mt-3 grid gap-2 ${media.length > 1 ? "sm:grid-cols-2" : "grid-cols-1"}`}>
      {media.map((m) => {
        const url = urls[m.downloadKey];
        const isVideo = m.mimeType.startsWith("video/");

        if (!url) {
          return (
            <div
              key={m.id}
              className="flex h-40 items-center justify-center border border-border bg-surface font-mono text-[0.7rem] text-text-faint"
            >
              loading media…
            </div>
          );
        }

        if (isVideo) return <VideoPlayer key={m.id} src={url} />;

        return (
          // Signed URLs from a runtime-configured host can't be known at build
          // time, so next/image's allowlist doesn't apply here.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            key={m.id}
            src={url}
            alt=""
            onClick={() => openLightbox({ url, kind: "image" })}
            className="max-h-[520px] w-full cursor-zoom-in border border-border object-contain"
          />
        );
      })}
    </div>
  );
}
