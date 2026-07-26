import { z } from "zod";
import { ZOrder, schemaWithCursorPagination, schemaWithPagination } from "@/utils.js";

export const ZPostType = z.enum(["comment", "post"]);
export const ZVoteType = z.enum(["upvote", "downvote"]);

export const ZPostMedia = z.object({
  id: z.string().uuid(),
  userId: z.string().uuid(),
  postId: z.string().uuid(),
  downloadKey: z.string(),
  fileSize: z.number(),
  mimeType: z.string(),
  createdAt: z.string().datetime(),
});

export const ZPost = z.object({
  id: z.string().uuid(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  authorId: z.string(),
  communityId: z.string().uuid().nullable(),
  parentPostId: z.string().uuid().nullable(),
  postType: ZPostType,
  title: z.string().nullable(),
  content: z.string().nullable(),
  upvotes: z.number(),
  downvotes: z.number(),
  commentCount: z.number(),
  engagement: z.number(),
  isPopularityUpdated: z.boolean(),
  deletedAt: z.string().datetime().nullable(),
});

export const ZPopulatedPost = ZPost.extend({
  postMedia: z.array(ZPostMedia),
});

// parentPostId requires postType "comment" and vice versa
export const ZCreatePostPayload = z.object({
  communityId: z.string().uuid().nullable().optional(),
  parentPostId: z.string().uuid().nullable().optional(),
  postType: ZPostType,
  title: z.string().nullable().optional(),
  content: z.string().nullable().optional(),
});

export const ZUpdatePostByIDPayload = z.object({
  title: z.string().nullable().optional(),
  content: z.string().nullable().optional(),
});

export const ZGetCommentsByPostIDQuery = z.object({
  cursorSortValue: z.string().optional(),
  cursorCreatedAt: z.string().datetime().optional(),
  sort: z.enum(["created_at", "popularity"]).optional(),
  order: ZOrder.optional(),
  dateCreatedStart: z.string().datetime().optional(),
  dateCreatedEnd: z.string().datetime().optional(),
});

export const ZGetRepliesByCommentIDQuery = z.object({
  page: z.coerce.number().min(1).optional(),
  limit: z.coerce.number().min(1).max(100).optional(),
});

export const ZReactToPostPayload = z.object({
  reaction: ZVoteType,
});

export const ZGetPostsQuery = z.object({
  referenceTime: z.string().datetime().optional(),
  trendingLimit: z.coerce.number().optional(),
  followingUsersLimit: z.coerce.number().optional(),
  followingCommunitiesLimit: z.coerce.number().optional(),
  trendingCursorValue: z.coerce.number().optional(),
  trendingCursorCreatedAt: z.string().datetime().optional(),
  followingUsersCursorCreatedAt: z.string().datetime().optional(),
  followingCommunitiesCursorCreatedAt: z.string().datetime().optional(),
});

export const ZGetPostsQueryResponse = z.object({
  referenceTime: z.string().datetime(),

  trendingPosts: z.array(ZPopulatedPost),
  nextTrendingCursorValue: z.number().nullable(),
  nextTrendingCursorCreatedAt: z.string().datetime().nullable(),
  hasMoreTrending: z.boolean(),

  followingUsersPosts: z.array(ZPopulatedPost),
  nextFollowingUsersCursorCreatedAt: z.string().datetime().nullable(),
  hasMoreFollowingUsers: z.boolean(),

  followingCommunitiesPosts: z.array(ZPopulatedPost),
  nextFollowingCommunitiesCursorCreatedAt: z.string().datetime().nullable(),
  hasMoreFollowingCommunities: z.boolean(),
});

export const ZPostCommentsResponse = schemaWithCursorPagination(ZPopulatedPost);
export const ZPostRepliesResponse = schemaWithPagination(ZPopulatedPost);
