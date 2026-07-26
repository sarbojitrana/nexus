import { z } from "zod";
import { ZOrder, schemaWithCursorPagination } from "@/utils.js";
import { ZMiniUser } from "@/user.js";

export const ZUserFollow = z.object({
  createdAt: z.string().datetime(),
  followerId: z.string(),
  followingId: z.string(),
});

export const ZCommunityFollow = z.object({
  createdAt: z.string().datetime(),
  followerId: z.string(),
  communityId: z.string().uuid(),
});

export const ZGetFollowersQuery = z.object({
  cursorCreatedAt: z.string().datetime().optional(),
  order: ZOrder.optional(),
});

export const ZGetFollowersResponse = schemaWithCursorPagination(ZMiniUser);
