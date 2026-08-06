import { z } from "zod";

export const ZPresignUploadPayload = z.object({
  mimeType: z.string(),
});

export const ZPresignUploadResponse = z.object({
  uploadUrl: z.string(),
  key: z.string(),
});
