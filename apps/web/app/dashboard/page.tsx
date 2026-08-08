import Link from "next/link";
import { PostCard } from "@/components/post-card";
import { FollowButton } from "@/components/follow-button";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { SearchBar } from "@/components/search-bar";
import { NewPostButton } from "@/components/new-post-modal";
import { RemoteAvatar } from "@/components/media/remote-image";
import { HudCorners } from "@/components/logo";
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
        ]).slice(0, 20)
      : [];

  const posts = await enrichPosts(rawPosts);

  const trending = trendingRes && trendingRes.status === 200 ? trendingRes.body.data.slice(0, 4) : [];
  const people =
    peopleRes && peopleRes.status === 200
      ? peopleRes.body.data.filter((u) => u.id !== viewerId).slice(0, 4)
      : [];

  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard" />

      <div className="mx-auto grid w-full max-w-[1240px] grid-cols-1 lg:grid-cols-[minmax(0,1fr)_330px]">
      <main className="flex w-full max-w-[820px] flex-col gap-4 px-6 py-6">
        <SearchBar />

        <div className="flex items-center justify-between border-b border-border-soft pb-3">
          <span className="font-mono text-[0.74rem] font-bold tracking-[0.08em] text-accent-strong uppercase">
            For you
          </span>
          <NewPostButton />
        </div>

        {posts.map((p) => (
          <PostCard key={p.id} post={p} />
        ))}

        {posts.length === 0 && (
          <div className="border border-border bg-surface p-8 text-center">
            <p className="text-[0.88rem] text-text-muted">Nothing to show yet.</p>
            <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
              Follow people or communities, or write the first post.
            </p>
          </div>
        )}
      </main>

      <aside className="hidden flex-col gap-6 border-l border-border-soft px-5 py-6 lg:flex">
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
                <RemoteAvatar storageKey={p.avatarKey} size={28} />
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
