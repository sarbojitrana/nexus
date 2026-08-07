import { z } from "zod";

export const ZPresignUploadPayload = z.object({
  mimeType: z.string(),
});

export const ZPresignUploadResponse = z.object({
  uploadUrl: z.string(),
  key: z.string(),
});

export const ZPresignDownloadsPayload = z.object({
  keys: z.array(z.string()).max(100),
});

export const ZPresignDownloadsResponse = z.object({
  urls: z.record(z.string(), z.string()),
});
