"use client";

import { useCallback } from "react";
import { useApi } from "@/lib/use-api";

export type UploadedMedia = {
  storageKey: string;
  fileSize: number;
  mimeType: string;
};

// Presign -> PUT straight to object storage -> hand back the key. File bytes
// never pass through the Nexus backend.
export function useUpload() {
  const api = useApi();

  return useCallback(
    async (file: File): Promise<UploadedMedia | null> => {
      const mimeType = file.type || "application/octet-stream";

      const presign = await api.Storage.presignUpload({ body: { mimeType } }).catch(() => null);
      if (!presign || presign.status !== 200) return null;

      const ok = await fetch(presign.body.uploadUrl, {
        method: "PUT",
        headers: { "Content-Type": mimeType },
        body: file,
      })
        .then((r) => r.ok)
        .catch(() => false);

      if (!ok) return null;
      return { storageKey: presign.body.key, fileSize: file.size, mimeType };
    },
    [api]
  );
}
