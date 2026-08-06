import { PostCard } from "@/components/post-card";
import { FollowButton } from "@/components/follow-button";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { SearchBar } from "@/components/search-bar";
import { NewPostButton } from "@/components/new-post-modal";
import { getServerApi, getServerUserId } from "@/lib/api-server";
import { formatTimeAgo } from "@/lib/format";
import type { User, PopulatedPost, MiniCommunity } from "@nexus/zod";

export default async function DashboardPage() {
  const [api, viewerId] = await Promise.all([getServerApi(), getServerUserId()]);
  const fetchOptions = { cache: "no-store" as const };

  const [feedRes, trendingCommunitiesRes, peopleRes] = await Promise.all([
    api.Post.getFeed({ query: {}, fetchOptions }),
    api.Community.getCommunities({
      query: { sort: "members_count", order: "desc" },
      fetchOptions,
    }),
    api.User.getUsers({ query: { sort: "follower_count", order: "desc" }, fetchOptions }),
  ]);

  const posts: PopulatedPost[] =
    feedRes.status === 200
      ? dedupeAndSort([
          ...feedRes.body.trendingPosts,
          ...feedRes.body.followingUsersPosts,
          ...feedRes.body.followingCommunitiesPosts,
        ]).slice(0, 15)
      : [];

  const trendingCommunities: MiniCommunity[] =
    trendingCommunitiesRes.status === 200 ? trendingCommunitiesRes.body.data.slice(0, 3) : [];
  const people =
    peopleRes.status === 200
      ? peopleRes.body.data.filter((u) => u.id !== viewerId).slice(0, 2)
      : [];

  const authorIds = [...new Set(posts.map((p) => p.authorId))];
  const communityIds = [...new Set(posts.map((p) => p.communityId).filter((id): id is string => !!id))];

  const [authorResults, communityResults] = await Promise.all([
    Promise.all(
      authorIds.map((id) =>
        api.User.getUserById({ params: { id }, fetchOptions }).catch(() => null)
      )
    ),
    Promise.all(
      communityIds.map((id) =>
        api.Community.getCommunityByIdOrSlug({ params: { idOrSlug: id }, fetchOptions }).catch(
          () => null
        )
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
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr_268px]">
      <DashboardSidebar active="/dashboard" />

      <main className="flex max-w-[640px] flex-col gap-4 px-6 py-5.5">
        <SearchBar />

        <div className="flex items-center justify-between">
          <div className="flex gap-1">
            <span className="rounded-full bg-accent/10 px-3 py-1.5 text-[0.8rem] font-bold text-accent-strong">
              For you
            </span>
          </div>
          <NewPostButton />
        </div>

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
            No posts yet. Follow some people or communities to fill your feed.
          </div>
        )}
      </main>

      <aside className="hidden flex-col gap-5.5 border-l border-border-soft px-4.5 py-5.5 lg:flex">
        <div className="rounded-2xl border border-border bg-surface p-4">
          <h4 className="mb-3 text-[0.82rem] font-bold">Trending communities</h4>
          <div className="flex flex-col gap-3">
            {trendingCommunities.map((t) => (
              <div key={t.communityId} className="flex items-center gap-2.5">
                <div className="h-[30px] w-[30px] shrink-0 rounded-[9px] bg-gradient-to-br from-accent to-down" />
                <div className="flex min-w-0 flex-col gap-px">
                  <strong className="text-[0.82rem] font-bold">n/{t.communityName}</strong>
                  <span className="text-[0.72rem] text-text-faint">
                    {t.membersCount} members
                  </span>
                </div>
                <FollowButton kind="community" id={t.communityId} />
              </div>
            ))}
            {trendingCommunities.length === 0 && (
              <span className="text-[0.8rem] text-text-faint">Nothing trending yet</span>
            )}
          </div>
        </div>
        <div className="rounded-2xl border border-border bg-surface p-4">
          <h4 className="mb-3 text-[0.82rem] font-bold">People to follow</h4>
          <div className="flex flex-col gap-3">
            {people.map((p) => (
              <div key={p.id} className="flex items-center gap-2.5">
                <div className="h-[30px] w-[30px] shrink-0 rounded-full bg-gradient-to-br from-accent to-down" />
                <div className="flex min-w-0 flex-col gap-px">
                  <strong className="text-[0.82rem] font-bold">@{p.username}</strong>
                  <span className="text-[0.72rem] text-text-faint">
                    {p.followerCount} followers
                  </span>
                </div>
                <FollowButton kind="user" id={p.id} />
              </div>
            ))}
            {people.length === 0 && (
              <span className="text-[0.8rem] text-text-faint">No suggestions yet</span>
            )}
          </div>
        </div>
      </aside>
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
