import { PostCard } from "@/components/post-card";
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

  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/following" />

      <main className="mx-auto flex w-full max-w-[820px] flex-col gap-4 px-6 py-6">
        <div className="border-b border-border-soft pb-3">
          <h1 className="font-display text-[1.3rem] font-extrabold">Following</h1>
          <p className="eyebrow mt-1">posts from people and communities you follow</p>
        </div>

        {posts.map((p) => (
          <PostCard key={p.id} post={p} />
        ))}

        {posts.length === 0 && (
          <div className="border border-border bg-surface p-8 text-center">
            <p className="text-[0.88rem] text-text-muted">Nothing to show yet.</p>
            <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
              Follow some people or communities to fill this feed.
            </p>
          </div>
        )}
      </main>
    </div>
  );
}
