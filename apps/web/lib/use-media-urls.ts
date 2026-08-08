"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";

// Storage keys aren't URLs. This batches however many keys a screen needs into
// a single presign request and hands back a key -> signed URL map.
export function useMediaUrls(keys: (string | null | undefined)[]) {
  const api = useApi();
  const [urls, setUrls] = useState<Record<string, string>>({});

  const present = keys.filter((k): k is string => !!k);
  const cacheKey = present.slice().sort().join(",");

  useEffect(() => {
    if (!cacheKey) return;
    let cancelled = false;

    api.Storage.presignDownloads({ body: { keys: cacheKey.split(",") } })
      .then((res) => {
        if (!cancelled && res.status === 200) {
          setUrls((prev) => ({ ...prev, ...res.body.urls }));
        }
      })
      .catch(() => {});

    return () => {
      cancelled = true;
    };
  }, [cacheKey, api]);

  return urls;
}
