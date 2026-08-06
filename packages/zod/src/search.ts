import { z } from "zod";

export const ZSearchQuery = z.object({
  q: z.string().min(1).max(200),
});

export const ZPostDoc = z.object({
  id: z.string().uuid(),
  authorId: z.string(),
  communityId: z.string().uuid().nullable(),
  title: z.string().nullable(),
  content: z.string().nullable(),
  upvotes: z.number(),
  commentCount: z.number(),
  createdAt: z.string().datetime(),
});

export const ZUserDoc = z.object({
  id: z.string(),
  username: z.string(),
  displayName: z.string(),
  bio: z.string().nullable(),
  followerCount: z.number(),
  createdAt: z.string().datetime(),
});

export const ZCommunityDoc = z.object({
  id: z.string().uuid(),
  name: z.string(),
  slug: z.string(),
  description: z.string().nullable(),
  membersCount: z.number(),
  createdAt: z.string().datetime(),
});

export const ZSearchResults = z.object({
  posts: z.array(ZPostDoc),
  users: z.array(ZUserDoc),
  communities: z.array(ZCommunityDoc),
});
export type SearchResults = z.infer<typeof ZSearchResults>;
export type PostDoc = z.infer<typeof ZPostDoc>;
export type UserDoc = z.infer<typeof ZUserDoc>;
export type CommunityDoc = z.infer<typeof ZCommunityDoc>;
