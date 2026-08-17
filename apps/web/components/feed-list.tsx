"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { PostCard } from "@/components/post-card";
import { useApi } from "@/lib/use-api";
import { enrichPostsWith, dedupeAndSortByNewest } from "@/lib/enrich-posts-shared";
import type { PostCardData } from "@/lib/post-card-data";
import type { PopulatedPost } from "@nexus/zod";

/** The feed is three independently-paged lanes, so "load more" carries a
 *  cursor per lane and a lane is exhausted on its own. */
export type FeedCursors = {
  referenceTime: string;
  trendingCursorValue: number | null;
  trendingCursorCreatedAt: string | null;
  hasMoreTrending: boolean;
  followingUsersCursorCreatedAt: string | null;
  hasMoreFollowingUsers: boolean;
  followingCommunitiesCursorCreatedAt: string | null;
  hasMoreFollowingCommunities: boolean;
};

export function FeedList({
  initialPosts,
  initialCursors,
  lanes = ["trending", "followingUsers", "followingCommunities"],
  emptyTitle = "Nothing to show yet.",
  emptyHint = "Follow people or communities, or write the first post.",
}: {
  initialPosts: PostCardData[];
  initialCursors: FeedCursors | null;
  /** Which lanes this feed draws from. The endpoint always returns all three,
   *  so a following-only feed has to drop trending from later pages or it
   *  would quietly start mixing it in. */
  lanes?: ("trending" | "followingUsers" | "followingCommunities")[];
  emptyTitle?: string;
  emptyHint?: string;
}) {
  const api = useApi();
  const [posts, setPosts] = useState(initialPosts);
  const [cursors, setCursors] = useState(initialCursors);
  const [isLoading, setIsLoading] = useState(false);
  const sentinelRef = useRef<HTMLDivElement>(null);

  const hasMore =
    !!cursors &&
    ((lanes.includes("trending") && cursors.hasMoreTrending) ||
      (lanes.includes("followingUsers") && cursors.hasMoreFollowingUsers) ||
      (lanes.includes("followingCommunities") && cursors.hasMoreFollowingCommunities));

  const loadMore = useCallback(async () => {
    if (!cursors || isLoading || !hasMore) return;
    setIsLoading(true);

    // Pinning referenceTime keeps paging stable: without it, posts published
    // mid-scroll would shift the window and duplicate or skip rows.
    const res = await api.Post.getFeed({
      query: {
        referenceTime: cursors.referenceTime,
        trendingCursorValue: cursors.trendingCursorValue ?? undefined,
        trendingCursorCreatedAt: cursors.trendingCursorCreatedAt ?? undefined,
        followingUsersCursorCreatedAt: cursors.followingUsersCursorCreatedAt ?? undefined,
        followingCommunitiesCursorCreatedAt:
          cursors.followingCommunitiesCursorCreatedAt ?? undefined,
      },
    }).catch(() => null);

    if (!res || res.status !== 200) {
      setIsLoading(false);
      return;
    }

    const raw: PopulatedPost[] = dedupeAndSortByNewest([
      ...(lanes.includes("trending") ? res.body.trendingPosts : []),
      ...(lanes.includes("followingUsers") ? res.body.followingUsersPosts : []),
      ...(lanes.includes("followingCommunities") ? res.body.followingCommunitiesPosts : []),
    ]);

    const enriched = await enrichPostsWith(api, raw);

    setPosts((prev) => {
      const seen = new Set(prev.map((p) => p.id));
      return [...prev, ...enriched.filter((p) => !seen.has(p.id))];
    });

    setCursors({
      referenceTime: res.body.referenceTime,
      trendingCursorValue: res.body.nextTrendingCursorValue,
      trendingCursorCreatedAt: res.body.nextTrendingCursorCreatedAt,
      hasMoreTrending: res.body.hasMoreTrending,
      followingUsersCursorCreatedAt: res.body.nextFollowingUsersCursorCreatedAt,
      hasMoreFollowingUsers: res.body.hasMoreFollowingUsers,
      followingCommunitiesCursorCreatedAt: res.body.nextFollowingCommunitiesCursorCreatedAt,
      hasMoreFollowingCommunities: res.body.hasMoreFollowingCommunities,
    });
    setIsLoading(false);
  }, [api, cursors, hasMore, isLoading, lanes]);

  // Auto-load as the sentinel scrolls into view, with a manual button as the
  // fallback for when IntersectionObserver never fires (very tall viewports).
  useEffect(() => {
    const el = sentinelRef.current;
    if (!el || !hasMore) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) loadMore();
      },
      { rootMargin: "400px" }
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [hasMore, loadMore]);

  if (posts.length === 0) {
    return (
      <div className="border border-border bg-surface p-8 text-center">
        <p className="text-[0.88rem] text-text-muted">{emptyTitle}</p>
        <p className="mt-1 font-mono text-[0.72rem] text-text-faint">{emptyHint}</p>
      </div>
    );
  }

  return (
    <>
      {posts.map((p) => (
        <PostCard key={p.id} post={p} />
      ))}

      <div ref={sentinelRef} className="py-3">
        {isLoading && (
          <div className="flex flex-col gap-2">
            <div className="feed-scanline" />
            <span className="text-center font-mono text-[0.66rem] tracking-[0.14em] text-text-faint uppercase">
              loading
            </span>
          </div>
        )}

        {!isLoading && hasMore && (
          <button
            onClick={loadMore}
            className="group flex w-full items-center gap-3 font-mono text-[0.66rem] tracking-[0.14em] text-text-faint uppercase"
          >
            <span className="h-px flex-1 bg-border transition-colors group-hover:bg-accent/50" />
            <span className="transition-colors group-hover:text-accent-strong">load more ↓</span>
            <span className="h-px flex-1 bg-border transition-colors group-hover:bg-accent/50" />
          </button>
        )}

        {!hasMore && (
          <div className="flex w-full items-center gap-3 font-mono text-[0.62rem] tracking-[0.16em] text-text-faint uppercase">
            <span className="h-px flex-1 bg-border-soft" />
            <span>end of feed</span>
            <span className="h-px flex-1 bg-border-soft" />
          </div>
        )}
      </div>
    </>
  );
}
