"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { apiErrorMessage } from "@/lib/api-error";
import { ShareButton } from "@/components/share-button";
import { MediaGallery } from "@/components/media/media-gallery";
import { formatCount, formatTimeAgo } from "@/lib/format";
import type { PopulatedPost } from "@nexus/zod";

type Meta = { usernames: Record<string, string>; communitySlug: string | null };

export function PostDetail({ postId }: { postId: string }) {
  const api = useApi();
  const router = useRouter();
  const { userId } = useAuth();

  const [post, setPost] = useState<PopulatedPost | null>(null);
  const [comments, setComments] = useState<PopulatedPost[]>([]);
  const [meta, setMeta] = useState<Meta>({ usernames: {}, communitySlug: null });
  const [isLoading, setIsLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  const [commentText, setCommentText] = useState("");
  const [isCommenting, setIsCommenting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");
  const [showReport, setShowReport] = useState(false);
  const [reportReason, setReportReason] = useState("");
  const [banner, setBanner] = useState<string | null>(null);

  const resolveNames = useCallback(
    async (ids: string[]) => {
      const unique = [...new Set(ids)];
      const results = await Promise.all(
        unique.map((id) => api.User.getUserById({ params: { id } }).catch(() => null))
      );
      const map: Record<string, string> = {};
      results.forEach((r) => {
        if (r && r.status === 200) map[r.body.id] = r.body.username;
      });
      setMeta((prev) => ({ ...prev, usernames: { ...prev.usernames, ...map } }));
    },
    [api]
  );

  const load = useCallback(async () => {
    const res = await api.Post.getPostById({ params: { id: postId } }).catch(() => null);
    if (!res || res.status !== 200) {
      setNotFound(true);
      setIsLoading(false);
      return;
    }
    setPost(res.body);
    setEditTitle(res.body.title ?? "");
    setEditContent(res.body.content ?? "");

    const commentsRes = await api.Post.getPostComments({
      params: { id: postId },
      query: {},
    }).catch(() => null);
    const loadedComments = commentsRes && commentsRes.status === 200 ? commentsRes.body.data : [];
    setComments(loadedComments);

    resolveNames([res.body.authorId, ...loadedComments.map((c) => c.authorId)]);

    if (res.body.communityId) {
      const communityRes = await api.Community.getCommunityByIdOrSlug({
        params: { idOrSlug: res.body.communityId },
      }).catch(() => null);
      if (communityRes && communityRes.status === 200) {
        setMeta((prev) => ({ ...prev, communitySlug: communityRes.body.slug }));
      }
    }

    setIsLoading(false);
  }, [api, postId, resolveNames]);

  useEffect(() => {
    load();
  }, [load]);

  async function submitComment() {
    if (!commentText.trim()) return;
    setIsCommenting(true);
    const res = await api.Post.createPost({
      body: {
        postType: "comment",
        parentPostId: postId,
        content: commentText.trim(),
        title: null,
      },
    }).catch(() => null);
    setIsCommenting(false);
    if (res && res.status === 201) {
      setCommentText("");
      load();
    } else {
      setBanner("Couldn't post that comment.");
    }
  }

  async function react(reaction: "upvote" | "downvote") {
    await api.Post.reactToPost({ params: { id: postId }, body: { reaction } }).catch(() => null);
    load();
  }

  async function saveEdit() {
    const res = await api.Post.updatePost({
      params: { id: postId },
      body: { title: editTitle.trim() || null, content: editContent.trim() || null },
    }).catch(() => null);
    if (res && res.status === 200) {
      setIsEditing(false);
      load();
    } else {
      setBanner("Couldn't save changes.");
    }
  }

  async function deletePost() {
    const res = await api.Post.deletePost({ params: { id: postId } }).catch(() => null);
    if (res && res.status === 204) router.push("/dashboard");
    else setBanner("Couldn't delete this post.");
  }

  async function submitReport() {
    if (!post?.communityId || !reportReason.trim()) return;
    const res = await api.Community.reportPost({
      params: { id: post.communityId },
      body: { postId: postId, reason: reportReason.trim() },
    }).catch(() => null);
    setShowReport(false);
    setReportReason("");
    setBanner(
      res && res.status === 201 ? "Report submitted to the moderators." : "Couldn't submit report."
    );
  }

  if (isLoading) {
    return <p className="px-6 py-8 font-mono text-[0.78rem] text-text-faint">loading…</p>;
  }

  if (notFound || !post) {
    return (
      <div className="px-6 py-8">
        <div className="border border-border bg-surface p-8 text-center text-[0.88rem] text-text-muted">
          This post doesn&apos;t exist or was removed.
        </div>
      </div>
    );
  }

  const isAuthor = post.authorId === userId;
  const score = post.upvotes - post.downvotes;

  return (
    <div className="mx-auto flex min-h-0 w-full max-w-[820px] flex-col gap-4 overflow-y-auto px-6 py-6">
      {banner && (
        <div className="border border-accent/40 bg-accent/5 px-4 py-2.5 font-mono text-[0.74rem] text-accent-strong">
          {banner}
        </div>
      )}

      <article className="border border-border bg-surface p-5">
        <div className="flex flex-wrap items-center gap-2 font-mono text-[0.7rem] text-text-faint">
          {meta.communitySlug && (
            <Link
              href={`/dashboard/communities/${meta.communitySlug}`}
              className="text-accent-strong hover:underline"
            >
              n/{meta.communitySlug}
            </Link>
          )}
          <Link href={`/dashboard/profile/${post.authorId}`} className="hover:text-text-muted">
            @{meta.usernames[post.authorId] ?? "unknown"}
          </Link>
          <span>·</span>
          <span>{formatTimeAgo(post.createdAt)}</span>
        </div>

        {isEditing ? (
          <div className="mt-3 flex flex-col gap-2">
            <input
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              className="border border-border bg-bg px-3 py-2 text-[0.9rem] focus:border-accent focus:outline-none"
            />
            <textarea
              value={editContent}
              onChange={(e) => setEditContent(e.target.value)}
              rows={5}
              className="resize-none border border-border bg-bg px-3 py-2 text-[0.86rem] focus:border-accent focus:outline-none"
            />
            <div className="flex gap-2">
              <button
                onClick={saveEdit}
                className="bg-accent px-4 py-1.5 font-mono text-[0.7rem] font-bold tracking-[0.06em] text-accent-text uppercase"
              >
                Save
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="px-4 py-1.5 font-mono text-[0.7rem] tracking-[0.06em] text-text-muted uppercase"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <>
            {post.title && (
              <h1 className="mt-2 font-display text-[1.35rem] leading-tight font-extrabold">
                {post.title}
              </h1>
            )}
            {post.content && (
              <p className="mt-2 text-[0.92rem] leading-relaxed whitespace-pre-wrap text-text-muted">
                {post.content}
              </p>
            )}
            <MediaGallery media={post.postMedia ?? []} />
          </>
        )}

        <div className="mt-4 flex flex-wrap items-center gap-4 border-t border-border-soft pt-3 font-mono text-[0.72rem] text-text-faint">
          <div className="flex items-center gap-2">
            <button onClick={() => react("upvote")} aria-label="Upvote" className="hover:text-up">
              ▲
            </button>
            <span className="font-bold tabular-nums text-text">{formatCount(score)}</span>
            <button
              onClick={() => react("downvote")}
              aria-label="Downvote"
              className="hover:text-down"
            >
              ▼
            </button>
          </div>
          <span>{post.commentCount} comments</span>

          {isAuthor && !isEditing && (
            <>
              <button onClick={() => setIsEditing(true)} className="hover:text-text-muted">
                edit
              </button>
              <button onClick={deletePost} className="hover:text-accent-strong">
                delete
              </button>
            </>
          )}
          {!isAuthor && post.communityId && (
            <button onClick={() => setShowReport(true)} className="hover:text-accent-strong">
              report
            </button>
          )}
        </div>
      </article>

      <section className="flex flex-col gap-2">
        <div className="flex gap-2">
          <input
            value={commentText}
            onChange={(e) => setCommentText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                submitComment();
              }
            }}
            placeholder="Add a comment…"
            className="flex-1 border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
          <button
            onClick={submitComment}
            disabled={isCommenting || !commentText.trim()}
            className="bg-accent px-5 py-2.5 font-mono text-[0.7rem] font-bold tracking-[0.06em] text-accent-text uppercase disabled:opacity-50"
          >
            Reply
          </button>
        </div>

        <div className="eyebrow mt-2">{comments.length} comments</div>

        {comments.map((c) => (
          <CommentRow
            key={c.id}
            comment={c}
            username={meta.usernames[c.authorId] ?? "unknown"}
            onChanged={load}
          />
        ))}

        {comments.length === 0 && (
          <div className="border border-border bg-surface p-6 text-center font-mono text-[0.74rem] text-text-faint">
            Nothing to show — be the first to comment.
          </div>
        )}
      </section>

      {showReport && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/70 p-4">
          <div className="flex w-full max-w-[420px] flex-col gap-3 border border-border bg-surface-raised p-5">
            <h3 className="font-mono text-[0.78rem] font-bold tracking-[0.08em] uppercase">
              Report post
            </h3>
            <textarea
              value={reportReason}
              onChange={(e) => setReportReason(e.target.value)}
              rows={4}
              placeholder="Why are you reporting this?"
              className="resize-none border border-border bg-surface px-3 py-2 text-[0.86rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
            />
            <div className="flex justify-end gap-2">
              <button
                onClick={() => setShowReport(false)}
                className="px-4 py-2 font-mono text-[0.7rem] tracking-[0.06em] text-text-muted uppercase"
              >
                Cancel
              </button>
              <button
                onClick={submitReport}
                disabled={!reportReason.trim()}
                className="bg-accent px-4 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] text-accent-text uppercase disabled:opacity-50"
              >
                Submit
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function CommentRow({
  comment,
  username,
  onChanged,
  depth = 0,
}: {
  comment: PopulatedPost;
  username: string;
  onChanged: () => void;
  depth?: number;
}) {
  const api = useApi();
  const { userId } = useAuth();

  const [replies, setReplies] = useState<PopulatedPost[]>([]);
  const [names, setNames] = useState<Record<string, string>>({});
  const [showReplies, setShowReplies] = useState(false);
  const [page, setPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [isLoading, setIsLoading] = useState(false);

  const [isReplying, setIsReplying] = useState(false);
  const [replyText, setReplyText] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Replies use offset pagination (page/limit) rather than a cursor, so
  // "load more" walks pages and appends.
  const loadPage = useCallback(
    async (next: number) => {
      setIsLoading(true);
      const res = await api.Post.getPostReplies({
        params: { id: comment.id },
        query: { page: next, limit: 5 },
      }).catch(() => null);
      setIsLoading(false);

      if (!res || res.status !== 200) return;

      setReplies((prev) => {
        const seen = new Set(prev.map((r) => r.id));
        return next === 1 ? res.body.data : [...prev, ...res.body.data.filter((r) => !seen.has(r.id))];
      });
      setPage(res.body.page);
      setTotalPages(res.body.totalPages);

      const ids = [...new Set(res.body.data.map((r) => r.authorId))];
      const resolved = await Promise.all(
        ids.map((id) =>
          api.User.getUserById({ params: { id } })
            .then((r) => (r.status === 200 ? ([id, r.body.username] as const) : null))
            .catch(() => null)
        )
      );
      setNames((prev) => {
        const nextNames = { ...prev };
        resolved.forEach((entry) => {
          if (entry) nextNames[entry[0]] = entry[1];
        });
        return nextNames;
      });
    },
    [api, comment.id]
  );

  async function toggleReplies() {
    if (showReplies) {
      setShowReplies(false);
      return;
    }
    setShowReplies(true);
    if (replies.length === 0) await loadPage(1);
  }

  async function submitReply() {
    if (!replyText.trim()) return;
    setIsSending(true);
    setError(null);

    const res = await api.Post.createPost({
      body: {
        postType: "comment",
        parentPostId: comment.id,
        content: replyText.trim(),
        title: null,
      },
    }).catch(() => null);
    setIsSending(false);

    if (res && res.status === 201) {
      setReplyText("");
      setIsReplying(false);
      setShowReplies(true);
      await loadPage(1);
      onChanged();
      return;
    }
    setError(apiErrorMessage(res, "Couldn't post that reply."));
  }

  async function remove() {
    const res = await api.Post.deletePost({ params: { id: comment.id } }).catch(() => null);
    if (res && res.status === 204) onChanged();
  }

  return (
    <div className="border border-border bg-surface p-3.5">
      <div className="flex flex-wrap items-center gap-2 font-mono text-[0.68rem] text-text-faint">
        <Link href={`/dashboard/profile/${comment.authorId}`} className="hover:text-text-muted">
          @{username}
        </Link>
        <span>·</span>
        <span>{formatTimeAgo(comment.createdAt)}</span>
      </div>

      <p className="mt-1.5 text-[0.86rem] leading-relaxed whitespace-pre-wrap text-text">
        {comment.content}
      </p>

      <div className="mt-2 flex flex-wrap gap-3 font-mono text-[0.68rem] text-text-faint">
        <button onClick={() => setIsReplying((v) => !v)} className="hover:text-accent-strong">
          reply
        </button>
        {comment.commentCount > 0 && (
          <button onClick={toggleReplies} className="hover:text-text-muted">
            {showReplies ? "hide replies" : `${comment.commentCount} replies`}
          </button>
        )}
        <ShareButton path={`/dashboard/posts/${comment.id}`} />
        {comment.authorId === userId && (
          <button onClick={remove} className="hover:text-accent-strong">
            delete
          </button>
        )}
      </div>

      {isReplying && (
        <div className="mt-2.5 flex flex-col gap-2">
          <textarea
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            rows={2}
            autoFocus
            placeholder={`Reply to @${username}…`}
            className="resize-none border border-border bg-bg px-3 py-2 text-[0.84rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
          {error && <p className="font-mono text-[0.68rem] text-accent-strong">{error}</p>}
          <div className="flex gap-2">
            <button
              onClick={submitReply}
              disabled={isSending || !replyText.trim()}
              className="bg-accent px-3 py-1.5 font-mono text-[0.66rem] font-bold tracking-[0.05em] text-accent-text uppercase disabled:opacity-50"
            >
              {isSending ? "Posting…" : "Reply"}
            </button>
            <button
              onClick={() => {
                setIsReplying(false);
                setError(null);
              }}
              className="px-3 py-1.5 font-mono text-[0.66rem] tracking-[0.05em] text-text-muted uppercase"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {showReplies && (
        <div className="mt-2.5 flex flex-col gap-2 border-l border-border-soft pl-3">
          {replies.map((r) => (
            <CommentRow
              key={r.id}
              comment={r}
              username={names[r.authorId] ?? "unknown"}
              onChanged={() => {
                loadPage(1);
                onChanged();
              }}
              depth={depth + 1}
            />
          ))}

          {page < totalPages && (
            <button
              onClick={() => loadPage(page + 1)}
              disabled={isLoading}
              className="self-start font-mono text-[0.68rem] text-accent-strong hover:underline disabled:opacity-50"
            >
              {isLoading ? "loading…" : `load more replies (${page}/${totalPages})`}
            </button>
          )}
        </div>
      )}
    </div>
  );
}
