import { PostCard } from "@/components/post-card";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { getServerApi } from "@/lib/api-server";
import { formatTimeAgo } from "@/lib/format";
import type { User, PopulatedPost } from "@nexus/zod";

export default async function FollowingPage() {
  const api = await getServerApi();
  const fetchOptions = { cache: "no-store" as const };

  const feedRes = await api.Post.getFeed({ query: {}, fetchOptions });

  const posts: PopulatedPost[] =
    feedRes.status === 200
      ? dedupeAndSort([...feedRes.body.followingUsersPosts, ...feedRes.body.followingCommunitiesPosts])
      : [];

  const authorIds = [...new Set(posts.map((p) => p.authorId))];
  const communityIds = [...new Set(posts.map((p) => p.communityId).filter((id): id is string => !!id))];

  const [authorResults, communityResults] = await Promise.all([
    Promise.all(
      authorIds.map((id) => api.User.getUserById({ params: { id }, fetchOptions }).catch(() => null))
    ),
    Promise.all(
      communityIds.map((id) =>
        api.Community.getCommunityByIdOrSlug({ params: { idOrSlug: id }, fetchOptions }).catch(() => null)
      )
    ),
  ]);

  const authorsById = new Map<string, User>();
  authorResults.forEach((r) => {
    if (r && r.status === 200) authorsById.set(r.body.id, r.body);
  });
  const communitiesById = new Map<string, { name: string; slug: string }>();
  communityResults.forEach((r) => {
    if (r && r.status === 200) communitiesById.set(r.body.id, r.body);
  });

  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr]">
      <DashboardSidebar active="/dashboard/following" />

      <main className="mx-auto flex w-full max-w-[640px] flex-col gap-4 px-6 py-8">
        <h1 className="font-display text-[1.4rem] font-extrabold">Following</h1>

        {posts.map((p) => {
          const author = authorsById.get(p.authorId);
          const community = p.communityId ? communitiesById.get(p.communityId) : undefined;
          return (
            <PostCard
              key={p.id}
              community={community ? `n/${community.slug}` : "general"}
              author={author ? `@${author.username}` : "@unknown"}
              timeAgo={formatTimeAgo(p.createdAt)}
              title={p.title ?? p.content ?? ""}
              snippet={p.title ? (p.content ?? "") : ""}
              votes={p.upvotes - p.downvotes}
              comments={p.commentCount}
            />
          );
        })}

        {posts.length === 0 && (
          <div className="rounded-2xl border border-border bg-surface p-6 text-center text-[0.86rem] text-text-muted">
            Nothing here yet — follow some people or communities to see their posts.
          </div>
        )}
      </main>
    </div>
  );
}

function dedupeAndSort(posts: PopulatedPost[]): PopulatedPost[] {
  const byId = new Map<string, PopulatedPost>();
  for (const p of posts) byId.set(p.id, p);
  return [...byId.values()].sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );
}
