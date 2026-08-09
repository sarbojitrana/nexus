import { FeedList, type FeedCursors } from "@/components/feed-list";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { getServerApi } from "@/lib/api-server";
import { enrichPosts, dedupeAndSortByNewest } from "@/lib/enrich-posts";

export default async function FollowingPage() {
  const api = await getServerApi();
  const feedRes = await api.Post.getFeed({
    query: {},
    fetchOptions: { cache: "no-store" },
  }).catch(() => null);

  const rawPosts =
    feedRes && feedRes.status === 200
      ? dedupeAndSortByNewest([
          ...feedRes.body.followingUsersPosts,
          ...feedRes.body.followingCommunitiesPosts,
        ])
      : [];

  const posts = await enrichPosts(rawPosts);

  // Trending is excluded here, so that lane is marked exhausted up front.
  const cursors: FeedCursors | null =
    feedRes && feedRes.status === 200
      ? {
          referenceTime: feedRes.body.referenceTime,
          trendingCursorValue: null,
          trendingCursorCreatedAt: null,
          hasMoreTrending: false,
          followingUsersCursorCreatedAt: feedRes.body.nextFollowingUsersCursorCreatedAt,
          hasMoreFollowingUsers: feedRes.body.hasMoreFollowingUsers,
          followingCommunitiesCursorCreatedAt:
            feedRes.body.nextFollowingCommunitiesCursorCreatedAt,
          hasMoreFollowingCommunities: feedRes.body.hasMoreFollowingCommunities,
        }
      : null;

  return (
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/following" />

      <main className="mx-auto flex min-h-0 w-full max-w-[820px] flex-col">
        <div className="shrink-0 border-b border-border-soft px-6 pt-6 pb-3">
          <h1 className="font-display text-[1.3rem] font-extrabold">Following</h1>
          <p className="eyebrow mt-1">posts from people and communities you follow</p>
        </div>

        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 py-4">
          <FeedList
            initialPosts={posts}
            initialCursors={cursors}
            lanes={["followingUsers", "followingCommunities"]}
            emptyHint="Follow some people or communities to fill this feed."
          />
        </div>
      </main>
    </div>
  );
}
