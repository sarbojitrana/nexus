import { z } from "zod";
import { ZOrder, schemaWithCursorPagination } from "@/utils.js";

export const ZCommunityRole = z.enum(["member", "moderator", "admin"]);
export const ZPostPermission = z.enum(["all", "moderator"]);
export const ZCommunityReportStatus = z.enum([
  "pending",
  "resolved",
  "dismissed",
]);

export const ZCommunity = z.object({
  id: z.string().uuid(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  adminId: z.string(),
  name: z.string(),
  slug: z.string(),
  description: z.string().nullable(),
  avatarKey: z.string().nullable(),
  bannerKey: z.string().nullable(),
  membersCount: z.number(),
  postsCount: z.number(),
  canPost: ZPostPermission.nullable(),
});

export const ZCommunityResponse = ZCommunity.extend({
  viewerRole: ZCommunityRole.nullable(),
});
export type Community = z.infer<typeof ZCommunity>;
export type CommunityResponse = z.infer<typeof ZCommunityResponse>;

export const ZCommunityMember = z.object({
  userId: z.string(),
  communityId: z.string().uuid(),
  role: ZCommunityRole,
  joinedAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const ZMiniCommunity = z.object({
  communityId: z.string().uuid(),
  slug: z.string(),
  communityName: z.string(),
  communityAvatarKey: z.string().nullable(),
  membersCount: z.number(),
  postsCount: z.number(),
  createdAt: z.string().datetime(),
});
export type MiniCommunity = z.infer<typeof ZMiniCommunity>;

export const ZMiniCommunityUser = z.object({
  userId: z.string(),
  avatarKey: z.string().nullable(),
  name: z.string(),
  joinedAt: z.string().datetime(),
  role: ZCommunityRole,
});
export type MiniCommunityUser = z.infer<typeof ZMiniCommunityUser>;
export type CommunityRole = z.infer<typeof ZCommunityRole>;

export const ZCommunityReport = z.object({
  id: z.string().uuid(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  reporterId: z.string(),
  communityId: z.string().uuid(),
  postId: z.string().uuid(),
  reason: z.string(),
  status: ZCommunityReportStatus,
});
export type CommunityReport = z.infer<typeof ZCommunityReport>;

export const ZBannedFromCommunityUser = z.object({
  communityId: z.string().uuid(),
  userId: z.string(),
  createdAt: z.string().datetime(),
});

export const ZCreateCommunityPayload = z.object({
  name: z.string().max(50),
  slug: z.string().max(50),
  description: z.string().max(1000).nullable().optional(),
  avatarKey: z.string().nullable().optional(),
  bannerKey: z.string().nullable().optional(),
});

export const ZUpdateCommunitySettingsPayload = z.object({
  name: z.string().max(50).nullable().optional(),
  slug: z.string().max(50).nullable().optional(),
  description: z.string().max(1000).nullable().optional(),
  avatarKey: z.string().nullable().optional(),
  bannerKey: z.string().nullable().optional(),
});

export const ZChangeMemberRoleInCommunityPayload = z.object({
  newRole: ZCommunityRole,
});

export const ZGetCommunitiesQuery = z.object({
  cursorSortValue: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
  sort: z.enum(["created_at", "members_count"]).optional(),
  order: ZOrder.optional(),
  name: z.string().optional(),
});

export const ZGetCommunityMembersQuery = z.object({
  cursorSortValue: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
  order: ZOrder.optional(),
  role: z.enum(["all", "moderator", "member"]).optional(),
});

export const ZReportCommunityPostPayload = z.object({
  postId: z.string().uuid(),
  reason: z.string().max(1000),
});

export const ZResolveCommunityPostReportPayload = z.object({
  updatedStatus: z.enum(["resolved", "dismissed"]),
});

export const ZBanCommunityMemberPayload = z.object({
  userIdToBan: z.string(),
});

export const ZGetCommunityReportsQuery = z.object({
  status: ZCommunityReportStatus.optional(),
  reportedDateStart: z.string().optional(),
  reportedDateEnd: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
});

export const ZGetCommunitiesResponse = schemaWithCursorPagination(ZMiniCommunity);
export const ZGetCommunityMembersResponse =
  schemaWithCursorPagination(ZMiniCommunityUser);
export const ZGetCommunityReportsResponse =
  schemaWithCursorPagination(ZCommunityReport);
