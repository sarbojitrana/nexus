import { z } from "zod";
import { ZOrder } from "@/utils.js";

export const ZUser = z.object({
  id: z.string(),
  username: z.string(),
  emailId: z.string(),
  displayName: z.string(),
  bio: z.string().nullable(),
  avatarKey: z.string().nullable(),
  bannerKey: z.string().nullable(),
  followerCount: z.number(),
  followingCount: z.number(),
  postsCount: z.number(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const ZMiniUser = z.object({
  id: z.string(),
  username: z.string(),
  displayName: z.string(),
  avatarKey: z.string(),
  bio: z.string(),
  followerCount: z.number(),
  createdAt: z.string().datetime(),
});

export const ZUpdateUserPayload = z.object({
  username: z.string().max(50).nullable().optional(),
  displayName: z.string().max(50).nullable().optional(),
  bio: z.string().max(1000).nullable().optional(),
  avatarKey: z.string().nullable().optional(),
  bannerKey: z.string().nullable().optional(),
});

export const ZGetUsersQuery = z.object({
  cursorSortValue: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
  sort: z.enum(["created_at", "follower_count"]).optional(),
  order: ZOrder.optional(),
  name: z.string().optional(),
  dateJoinedStart: z.string().datetime().optional(),
  dateJoinedEnd: z.string().datetime().optional(),
});

export const ZGetPostsByUserIDQuery = z.object({
  cursorSortValue: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
  sort: z.enum(["created_at", "upvotes"]).optional(),
  order: ZOrder.optional(),
  dateCreatedStart: z.string().datetime().optional(),
  dateCreatedEnd: z.string().datetime().optional(),
});
