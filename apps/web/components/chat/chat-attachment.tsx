"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import { VideoPlayer } from "@/components/media/video-player";

export function ChatAttachment({
  storageKey,
  mimeType,
  fileSize,
}: {
  storageKey: string;
  mimeType: string | null;
  fileSize: number | null;
}) {
  const api = useApi();
  const [url, setUrl] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    api.Storage.presignDownloads({ body: { keys: [storageKey] } })
      .then((res) => {
        if (!cancelled && res.status === 200) setUrl(res.body.urls[storageKey] ?? null);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [storageKey, api]);

  const sizeLabel = fileSize ? ` · ${Math.round(fileSize / 1024)}KB` : "";

  if (!url) {
    return (
      <div className="mt-1.5 bg-black/20 px-2 py-1 font-mono text-[0.7rem]">
        loading attachment…
      </div>
    );
  }

  if (mimeType?.startsWith("image/")) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img src={url} alt="" className="mt-1.5 max-h-[300px] w-full object-contain" />
    );
  }

  if (mimeType?.startsWith("video/")) {
    return (
      <div className="mt-1.5">
        <VideoPlayer src={url} />
      </div>
    );
  }

  return (
    <a
      href={url}
      target="_blank"
      rel="noreferrer"
      className="mt-1.5 flex items-center gap-1.5 bg-black/20 px-2 py-1 font-mono text-[0.7rem] underline"
    >
      📎 {mimeType ?? "file"}
      {sizeLabel}
    </a>
  );
}
