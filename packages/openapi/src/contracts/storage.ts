import { initContract } from "@ts-rest/core";
import {
  ZPresignUploadPayload,
  ZPresignUploadResponse,
  ZPresignDownloadsPayload,
  ZPresignDownloadsResponse,
} from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const storageContract = c.router(
  {
    presignUpload: {
      summary: "Get a presigned URL to upload a file",
      method: "POST",
      path: "/uploads/presign",
      body: ZPresignUploadPayload,
      responses: {
        200: ZPresignUploadResponse,
      },
      metadata: getSecurityMetadata(),
    },
    presignDownloads: {
      summary: "Resolve storage keys to temporary read URLs",
      method: "POST",
      path: "/uploads/download-urls",
      body: ZPresignDownloadsPayload,
      responses: {
        200: ZPresignDownloadsResponse,
      },
      metadata: getSecurityMetadata(),
    },
  },
  { pathPrefix: "/api/v1" }
);
