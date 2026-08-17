"use client";

import Link from "next/link";
import { useState } from "react";
import { useApi } from "@/lib/use-api";
import { MediaGallery } from "@/components/media/media-gallery";
import { formatTimeAgo } from "@/lib/format";
import { notify } from "@/lib/notify";
import { ShareButton } from "@/components/share-button";
import { apiErrorMessage } from "@/lib/api-error";
import { VoteButtons } from "@/components/vote-buttons";
import type { PostCardData } from "@/lib/post-card-data";

export function PostCard({ post, hideCommunity }: { post: PostCardData; hideCommunity?: boolean }) {
  const api = useApi();
  const [isCommenting, setIsCommenting] = useState(false);
  const [commentText, setCommentText] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [commentCount, setCommentCount] = useState(post.commentCount);
  const [error, setError] = useState<string | null>(null);

  async function submitComment() {
    if (!commentText.trim()) return;
    setIsSending(true);
    setError(null);

    const res = await api.Post.createPost({
      body: {
        postType: "comment",
        parentPostId: post.id,
        content: commentText.trim(),
        title: null,
      },
    }).catch(() => null);
    setIsSending(false);

    if (res && res.status === 201) {
      setCommentText("");
      setIsCommenting(false);
      setCommentCount((c) => c + 1);
      notify("Comment added", post.title ?? "Your comment was posted.");
      return;
    }
    setError(apiErrorMessage(res, "Couldn't post that comment."));
  }

  return (
    <article className="flex gap-4 border border-border bg-surface p-4">
      <VoteButtons postId={post.id} upvotes={post.upvotes} downvotes={post.downvotes} />

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

        <div className="mt-3 flex flex-wrap gap-4 font-mono text-[0.7rem] text-text-faint">
          <button
            onClick={() => setIsCommenting((v) => !v)}
            className="hover:text-accent-strong"
          >
            comment
          </button>
          <Link href={`/dashboard/posts/${post.id}`} className="hover:text-text-muted">
            {commentCount} comments
          </Link>
          <ShareButton path={`/dashboard/posts/${post.id}`} title={post.title ?? undefined} />
        </div>

        {isCommenting && (
          <div className="mt-2.5 flex flex-col gap-2">
            <textarea
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
              rows={2}
              autoFocus
              placeholder="Write a comment…"
              className="resize-none border border-border bg-bg px-3 py-2 text-[0.84rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
            />
            {error && <p className="font-mono text-[0.68rem] text-accent-strong">{error}</p>}
            <div className="flex gap-2">
              <button
                onClick={submitComment}
                disabled={isSending || !commentText.trim()}
                className="bg-accent px-3 py-1.5 font-mono text-[0.66rem] font-bold tracking-[0.05em] text-accent-text uppercase disabled:opacity-50"
              >
                {isSending ? "Posting…" : "Comment"}
              </button>
              <button
                onClick={() => {
                  setIsCommenting(false);
                  setError(null);
                }}
                className="px-3 py-1.5 font-mono text-[0.66rem] tracking-[0.05em] text-text-muted uppercase"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>
    </article>
  );
}

export type { PostCardData };
