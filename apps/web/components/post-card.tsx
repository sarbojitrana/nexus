"use client";

import Link from "next/link";
import { useState } from "react";
import { useApi } from "@/lib/use-api";
import { MediaGallery } from "@/components/media/media-gallery";
import { formatCount, formatTimeAgo } from "@/lib/format";
import type { PopulatedPost, PostMedia } from "@nexus/zod";

export type PostCardData = {
  id: string;
  communitySlug: string | null;
  communityId: string | null;
  authorUsername: string;
  authorId: string;
  createdAt: string;
  title: string | null;
  content: string | null;
  upvotes: number;
  downvotes: number;
  commentCount: number;
  media: PostMedia[];
};

export function toPostCardData(
  post: PopulatedPost,
  authorUsername: string,
  communitySlug: string | null
): PostCardData {
  return {
    id: post.id,
    communitySlug,
    communityId: post.communityId,
    authorUsername,
    authorId: post.authorId,
    createdAt: post.createdAt,
    title: post.title,
    content: post.content,
    upvotes: post.upvotes,
    downvotes: post.downvotes,
    commentCount: post.commentCount,
    media: post.postMedia ?? [],
  };
}

export function PostCard({ post, hideCommunity }: { post: PostCardData; hideCommunity?: boolean }) {
  const api = useApi();
  // Optimistic vote state -- the API returns 204 with no body, and there's no
  // "my vote" field on the post, so the UI tracks the delta locally.
  const [vote, setVote] = useState<"upvote" | "downvote" | null>(null);
  const [delta, setDelta] = useState(0);

  async function react(reaction: "upvote" | "downvote") {
    const previous = vote;
    const next = previous === reaction ? null : reaction;

    const base = previous === "upvote" ? -1 : previous === "downvote" ? 1 : 0;
    const add = next === "upvote" ? 1 : next === "downvote" ? -1 : 0;

    setVote(next);
    setDelta(delta + base + add);

    const res = await api.Post.reactToPost({
      params: { id: post.id },
      body: { reaction },
    }).catch(() => null);

    if (!res || (res.status !== 204 && res.status !== 200)) {
      setVote(previous);
      setDelta(delta);
    }
  }

  const score = post.upvotes - post.downvotes + delta;

  return (
    <article className="flex gap-4 border border-border bg-surface p-4">
      <div className="flex min-w-[34px] flex-col items-center gap-1">
        <button
          onClick={() => react("upvote")}
          aria-label="Upvote"
          className={`font-mono text-[0.9rem] leading-none ${
            vote === "upvote" ? "text-up" : "text-text-faint hover:text-up"
          }`}
        >
          ▲
        </button>
        <span className="font-mono text-[0.78rem] font-bold tabular-nums">
          {formatCount(score)}
        </span>
        <button
          onClick={() => react("downvote")}
          aria-label="Downvote"
          className={`font-mono text-[0.9rem] leading-none ${
            vote === "downvote" ? "text-down" : "text-text-faint hover:text-down"
          }`}
        >
          ▼
        </button>
      </div>

      <div className="flex min-w-0 flex-1 flex-col">
        <div className="flex flex-wrap items-center gap-2 font-mono text-[0.7rem] text-text-faint">
          {!hideCommunity && post.communitySlug && (
            <Link
              href={`/dashboard/communities/${post.communitySlug}`}
              className="text-accent-strong hover:underline"
            >
              n/{post.communitySlug}
            </Link>
          )}
          <Link href={`/dashboard/profile/${post.authorId}`} className="hover:text-text-muted">
            @{post.authorUsername}
          </Link>
          <span>·</span>
          <span>{formatTimeAgo(post.createdAt)}</span>
        </div>

        <Link href={`/dashboard/posts/${post.id}`} className="group">
          <h3 className="mt-1.5 font-display text-[1rem] leading-snug font-bold group-hover:text-accent-strong">
            {post.title ?? post.content?.slice(0, 90)}
          </h3>
        </Link>

        {post.title && post.content && (
          <p className="mt-1 line-clamp-3 text-[0.86rem] leading-relaxed text-text-muted">
            {post.content}
          </p>
        )}

        <MediaGallery media={post.media} />

        <div className="mt-3 flex gap-4 font-mono text-[0.7rem] text-text-faint">
          <Link href={`/dashboard/posts/${post.id}`} className="hover:text-text-muted">
            {post.commentCount} comments
          </Link>
        </div>
      </div>
    </article>
  );
}
