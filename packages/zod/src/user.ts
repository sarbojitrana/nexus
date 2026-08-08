import { z } from "zod";
import { ZOrder } from "@/utils.js";

export const ZProfileVisibility = z.enum(["public", "followers_only", "private"]);
export const ZGroupInvitePermission = z.enum(["everyone", "followers_only", "no_one"]);

export const ZUser = z.object({
  id: z.string(),
  username: z.string(),
  emailId: z.string(),
  displayName: z.string(),
  bio: z.string().nullable(),
  avatarKey: z.string().nullable(),
  avatarUrl: z.string().nullable(),
  bannerKey: z.string().nullable(),
  followerCount: z.number(),
  followingCount: z.number(),
  postsCount: z.number(),
  profileVisibility: ZProfileVisibility,
  showOnlineStatus: z.boolean(),
  groupInvitePermission: ZGroupInvitePermission,
  shareReadReceipts: z.boolean(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});
export type User = z.infer<typeof ZUser>;

export const ZUpdateUserSettingsPayload = z.object({
  profileVisibility: ZProfileVisibility.optional(),
  showOnlineStatus: z.boolean().optional(),
  groupInvitePermission: ZGroupInvitePermission.optional(),
  shareReadReceipts: z.boolean().optional(),
});

export const ZMiniUser = z.object({
  id: z.string(),
  username: z.string(),
  displayName: z.string(),
  avatarKey: z.string().nullable(),
  avatarUrl: z.string().nullable(),
  bio: z.string().nullable(),
  followerCount: z.number(),
  createdAt: z.string().datetime(),
});
export type MiniUser = z.infer<typeof ZMiniUser>;

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
