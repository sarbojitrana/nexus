"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import type { MiniCommunity } from "@nexus/zod";

export function NewPostButton() {
  const api = useApi();
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [communityId, setCommunityId] = useState("");
  const [communities, setCommunities] = useState<MiniCommunity[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen) return;
    api.Community.getCommunities({ query: {} }).then((res) => {
      if (res.status === 200) setCommunities(res.body.data);
    });
  }, [isOpen, api]);

  async function submit() {
    if (!title.trim()) {
      setError("Title is required.");
      return;
    }
    setIsSubmitting(true);
    setError(null);
    const res = await api.Post.createPost({
      body: {
        postType: "post",
        title: title.trim(),
        content: content.trim() || null,
        communityId: communityId || null,
      },
    }).catch(() => null);
    setIsSubmitting(false);
    if (res && res.status === 201) {
      setIsOpen(false);
      setTitle("");
      setContent("");
      setCommunityId("");
      router.refresh();
    } else {
      setError("Couldn't create the post. Try again.");
    }
  }

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="rounded-[9px] bg-accent px-4 py-2 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong"
      >
        + New post
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 p-4">
          <div className="flex w-full max-w-[480px] flex-col gap-4 rounded-2xl border border-border bg-surface-raised p-5">
            <h3 className="text-[1rem] font-bold">New post</h3>

            <select
              value={communityId}
              onChange={(e) => setCommunityId(e.target.value)}
              className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
            >
              <option value="">General (no community)</option>
              {communities.map((c) => (
                <option key={c.communityId} value={c.communityId}>
                  n/{c.slug}
                </option>
              ))}
            </select>

            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Title"
              className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
            />

            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="What's on your mind? (optional)"
              rows={5}
              className="resize-none rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
            />

            {error && <p className="text-[0.78rem] text-accent-strong">{error}</p>}

            <div className="flex justify-end gap-2">
              <button
                onClick={() => setIsOpen(false)}
                className="rounded-[9px] px-4 py-2 text-[0.82rem] font-bold text-text-muted hover:bg-surface"
              >
                Cancel
              </button>
              <button
                onClick={submit}
                disabled={isSubmitting || !title.trim()}
                className="rounded-[9px] bg-accent px-4 py-2 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
              >
                {isSubmitting ? "Posting..." : "Post"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
