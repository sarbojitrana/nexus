import { getServerApi } from "@/lib/api-server";
import { toPostCardData, type PostCardData } from "@/lib/post-card-data";
import type { PopulatedPost } from "@nexus/zod";

// The feed/post endpoints return raw author and community ids. Rather than
// N+1'ing per card, resolve the unique set in parallel once per page.
export async function enrichPosts(posts: PopulatedPost[]): Promise<PostCardData[]> {
  if (posts.length === 0) return [];

  const api = await getServerApi();
  const fetchOptions = { cache: "no-store" as const };

  const authorIds = [...new Set(posts.map((p) => p.authorId))];
  const communityIds = [
    ...new Set(posts.map((p) => p.communityId).filter((id): id is string => !!id)),
  ];

  const [authorResults, communityResults] = await Promise.all([
    Promise.all(
      authorIds.map((id) => api.User.getUserById({ params: { id }, fetchOptions }).catch(() => null))
    ),
    Promise.all(
      communityIds.map((id) =>
        api.Community.getCommunityByIdOrSlug({ params: { idOrSlug: id }, fetchOptions }).catch(
          () => null
        )
      )
    ),
  ]);

  const usernameById = new Map<string, string>();
  authorResults.forEach((r) => {
    if (r && r.status === 200) usernameById.set(r.body.id, r.body.username);
  });

  const slugById = new Map<string, string>();
  communityResults.forEach((r) => {
    if (r && r.status === 200) slugById.set(r.body.id, r.body.slug);
  });

  return posts.map((p) =>
    toPostCardData(
      p,
      usernameById.get(p.authorId) ?? "unknown",
      p.communityId ? (slugById.get(p.communityId) ?? null) : null
    )
  );
}

export function dedupeAndSortByNewest(posts: PopulatedPost[]): PopulatedPost[] {
  const byId = new Map<string, PopulatedPost>();
  for (const p of posts) byId.set(p.id, p);
  return [...byId.values()].sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );
}
