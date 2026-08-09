import Link from "next/link";
import { FollowButton } from "@/components/follow-button";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { SearchBar } from "@/components/search-bar";
import { NewPostButton } from "@/components/new-post-modal";
import { RemoteAvatar } from "@/components/media/remote-image";
import { HudCorners } from "@/components/logo";
import { FeedList, type FeedCursors } from "@/components/feed-list";
import { getServerApi, getServerUserId } from "@/lib/api-server";
import { enrichPosts, dedupeAndSortByNewest } from "@/lib/enrich-posts";

export default async function DashboardPage() {
  const [api, viewerId] = await Promise.all([getServerApi(), getServerUserId()]);
  const fetchOptions = { cache: "no-store" as const };

  const [feedRes, trendingRes, peopleRes] = await Promise.all([
    api.Post.getFeed({ query: {}, fetchOptions }).catch(() => null),
    api.Community.getCommunities({
      query: { sort: "members_count", order: "desc" },
      fetchOptions,
    }).catch(() => null),
    api.User.getUsers({ query: { sort: "follower_count", order: "desc" }, fetchOptions }).catch(
      () => null
    ),
  ]);

  const rawPosts =
    feedRes && feedRes.status === 200
      ? dedupeAndSortByNewest([
          ...feedRes.body.trendingPosts,
          ...feedRes.body.followingUsersPosts,
          ...feedRes.body.followingCommunitiesPosts,
        ])
      : [];

  const posts = await enrichPosts(rawPosts);

  const cursors: FeedCursors | null =
    feedRes && feedRes.status === 200
      ? {
          referenceTime: feedRes.body.referenceTime,
          trendingCursorValue: feedRes.body.nextTrendingCursorValue,
          trendingCursorCreatedAt: feedRes.body.nextTrendingCursorCreatedAt,
          hasMoreTrending: feedRes.body.hasMoreTrending,
          followingUsersCursorCreatedAt: feedRes.body.nextFollowingUsersCursorCreatedAt,
          hasMoreFollowingUsers: feedRes.body.hasMoreFollowingUsers,
          followingCommunitiesCursorCreatedAt:
            feedRes.body.nextFollowingCommunitiesCursorCreatedAt,
          hasMoreFollowingCommunities: feedRes.body.hasMoreFollowingCommunities,
        }
      : null;

  const trending = trendingRes && trendingRes.status === 200 ? trendingRes.body.data.slice(0, 6) : [];
  const people =
    peopleRes && peopleRes.status === 200
      ? peopleRes.body.data.filter((u) => u.id !== viewerId).slice(0, 6)
      : [];

  return (
    // h-dvh + overflow-hidden pins the shell to the viewport so only the feed
    // column scrolls -- the nav, search and composer stay put.
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard" />

      <div className="mx-auto grid w-full max-w-[1240px] grid-cols-1 overflow-hidden lg:grid-cols-[minmax(0,1fr)_330px]">
        <main className="flex min-h-0 w-full max-w-[820px] flex-col">
          <div className="shrink-0 space-y-4 px-6 pt-6 pb-3">
            <SearchBar />
            <div className="flex items-center justify-between border-b border-border-soft pb-3">
              <span className="font-mono text-[0.74rem] font-bold tracking-[0.08em] text-accent-strong uppercase">
                For you
              </span>
              <NewPostButton />
            </div>
          </div>

          <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pb-6">
            <FeedList initialPosts={posts} initialCursors={cursors} />
          </div>
        </main>

        <aside className="hidden min-h-0 flex-col gap-6 overflow-y-auto border-l border-border-soft px-5 py-6 lg:flex">
          <section>
            <h4 className="eyebrow mb-3">Trending communities</h4>
            <div className="flex flex-col gap-3">
              {trending.map((t) => (
                <div key={t.communityId} className="flex items-center gap-2.5">
                  <RemoteAvatar storageKey={t.communityAvatarKey} size={28} />
                  <Link href={`/dashboard/communities/${t.slug}`} className="min-w-0 flex-1">
                    <strong className="block truncate text-[0.82rem] font-bold hover:text-accent-strong">
                      n/{t.slug}
                    </strong>
                    <span className="font-mono text-[0.66rem] text-text-faint">
                      {t.membersCount} members
                    </span>
                  </Link>
                  <FollowButton kind="community" id={t.communityId} />
                </div>
              ))}
              {trending.length === 0 && (
                <span className="font-mono text-[0.72rem] text-text-faint">Nothing to show</span>
              )}
            </div>
          </section>

          <section>
            <h4 className="eyebrow mb-3">People to follow</h4>
            <div className="flex flex-col gap-3">
              {people.map((p) => (
                <div key={p.id} className="flex items-center gap-2.5">
                  <RemoteAvatar storageKey={p.avatarKey} url={p.avatarUrl} size={28} />
                  <Link href={`/dashboard/profile/${p.id}`} className="min-w-0 flex-1">
                    <strong className="block truncate text-[0.82rem] font-bold hover:text-accent-strong">
                      @{p.username}
                    </strong>
                    <span className="font-mono text-[0.66rem] text-text-faint">
                      {p.followerCount} followers
                    </span>
                  </Link>
                  <FollowButton kind="user" id={p.id} />
                </div>
              ))}
              {people.length === 0 && (
                <span className="font-mono text-[0.72rem] text-text-faint">Nothing to show</span>
              )}
            </div>
          </section>
        </aside>
      </div>
    </div>
  );
}
