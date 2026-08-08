"use client";

import { useCallback } from "react";
import { useApi } from "@/lib/use-api";
import { apiErrorMessage } from "@/lib/api-error";

export type UploadedMedia = {
  storageKey: string;
  fileSize: number;
  mimeType: string;
};

export type UploadResult =
  | { ok: true; media: UploadedMedia }
  | { ok: false; error: string };

// Presign -> PUT straight to object storage -> hand back the key. File bytes
// never pass through the Nexus backend.
//
// The two steps fail for completely different reasons, so they report
// separately: presign failing means the backend has no storage credentials,
// while the PUT failing is almost always the bucket's own CORS policy
// rejecting a browser request.
export function useUpload() {
  const api = useApi();

  return useCallback(
    async (file: File): Promise<UploadResult> => {
      const mimeType = file.type || "application/octet-stream";

      const presign = await api.Storage.presignUpload({ body: { mimeType } }).catch(() => null);
      if (!presign || presign.status !== 200) {
        return { ok: false, error: apiErrorMessage(presign, "Couldn't prepare the upload.") };
      }

      try {
        const res = await fetch(presign.body.uploadUrl, {
          method: "PUT",
          headers: { "Content-Type": mimeType },
          body: file,
        });

        if (!res.ok) {
          return {
            ok: false,
            error: `Storage rejected the upload (HTTP ${res.status}). Check the bucket's permissions.`,
          };
        }
      } catch {
        // fetch throws (rather than returning a status) when the browser
        // blocks the response -- for a direct-to-bucket PUT that's the
        // bucket's CORS policy not allowing this origin.
        return {
          ok: false,
          error:
            "Upload blocked by the storage bucket's CORS policy — allow PUT from this site in your R2 bucket settings.",
        };
      }

      return { ok: true, media: { storageKey: presign.body.key, fileSize: file.size, mimeType } };
    },
    [api]
  );
}
