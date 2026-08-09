import "server-only";

import { getServerApi } from "@/lib/api-server";
import { enrichPostsWith } from "@/lib/enrich-posts-shared";
import type { PostCardData } from "@/lib/post-card-data";
import type { PopulatedPost } from "@nexus/zod";

// Server-side wrapper. The isomorphic logic lives in enrich-posts-shared so
// the client feed can reuse it without importing server auth.
export async function enrichPosts(posts: PopulatedPost[]): Promise<PostCardData[]> {
  if (posts.length === 0) return [];
  const api = await getServerApi();
  return enrichPostsWith(api as never, posts);
}

export { dedupeAndSortByNewest } from "@/lib/enrich-posts-shared";
