import { toPostCardData, type PostCardData } from "@/lib/post-card-data";
import type { PopulatedPost } from "@nexus/zod";

// No server-only imports in this file on purpose: the client feed pages
// through it too, and pulling in getServerApi (which reaches for Clerk's
// server auth()) would drag server code into the browser bundle.

type EnrichClient = {
  User: {
    getUserById: (args: { params: { id: string } }) => Promise<{ status: number; body: unknown }>;
  };
  Community: {
    getCommunityByIdOrSlug: (args: {
      params: { idOrSlug: string };
    }) => Promise<{ status: number; body: unknown }>;
  };
};

// The feed endpoints return raw author and community ids. Rather than
// N+1'ing per card, resolve the unique set in parallel once per batch.
export async function enrichPostsWith(
  api: EnrichClient,
  posts: PopulatedPost[]
): Promise<PostCardData[]> {
  if (posts.length === 0) return [];

  const authorIds = [...new Set(posts.map((p) => p.authorId))];
  const communityIds = [
    ...new Set(posts.map((p) => p.communityId).filter((id): id is string => !!id)),
  ];

  const [authorResults, communityResults] = await Promise.all([
    Promise.all(authorIds.map((id) => api.User.getUserById({ params: { id } }).catch(() => null))),
    Promise.all(
      communityIds.map((id) =>
        api.Community.getCommunityByIdOrSlug({ params: { idOrSlug: id } }).catch(() => null)
      )
    ),
  ]);

  const usernameById = new Map<string, string>();
  authorResults.forEach((r) => {
    if (r && r.status === 200) {
      const body = r.body as { id: string; username: string };
      usernameById.set(body.id, body.username);
    }
  });

  const slugById = new Map<string, string>();
  communityResults.forEach((r) => {
    if (r && r.status === 200) {
      const body = r.body as { id: string; slug: string };
      slugById.set(body.id, body.slug);
    }
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
