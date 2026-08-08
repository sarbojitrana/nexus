import type { PopulatedPost, PostMedia } from "@nexus/zod";

// Deliberately kept out of post-card.tsx. That file is a client component, and
// anything exported from a "use client" module becomes a client reference --
// so calling this from a Server Component would throw. Server code enriching
// posts for render needs it, so the plain data shaping lives here instead.

export type PostCardData = {
  id: string;
  communitySlug: string | null;
  communityId: string | null;
  authorUsername: string;
  authorId: string;
  createdAt: string;
  title: string | null;
  content: string | null;
  upvotes: number;
  downvotes: number;
  commentCount: number;
  media: PostMedia[];
};

export function toPostCardData(
  post: PopulatedPost,
  authorUsername: string,
  communitySlug: string | null
): PostCardData {
  return {
    id: post.id,
    communitySlug,
    communityId: post.communityId,
    authorUsername,
    authorId: post.authorId,
    createdAt: post.createdAt,
    title: post.title,
    content: post.content,
    upvotes: post.upvotes,
    downvotes: post.downvotes,
    commentCount: post.commentCount,
    media: post.postMedia ?? [],
  };
}
