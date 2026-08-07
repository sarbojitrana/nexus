"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import { useUpload, type UploadedMedia } from "@/lib/use-upload";
import type { MiniCommunity } from "@nexus/zod";

type Staged = UploadedMedia & { name: string; previewUrl: string };

export function NewPostButton({ communityId: fixedCommunityId }: { communityId?: string }) {
  const api = useApi();
  const upload = useUpload();
  const router = useRouter();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [communityId, setCommunityId] = useState(fixedCommunityId ?? "");
  const [communities, setCommunities] = useState<MiniCommunity[]>([]);
  const [media, setMedia] = useState<Staged[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen || fixedCommunityId) return;
    api.Community.getCommunities({ query: {} }).then((res) => {
      if (res.status === 200) setCommunities(res.body.data);
    });
  }, [isOpen, api, fixedCommunityId]);

  // Object URLs are revoked on unmount so previews don't leak memory.
  useEffect(() => {
    return () => {
      media.forEach((m) => URL.revokeObjectURL(m.previewUrl));
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleFiles(files: FileList | null) {
    if (!files || files.length === 0) return;
    setIsUploading(true);
    setError(null);

    for (const file of Array.from(files).slice(0, 8 - media.length)) {
      const uploaded = await upload(file);
      if (uploaded) {
        setMedia((prev) => [
          ...prev,
          { ...uploaded, name: file.name, previewUrl: URL.createObjectURL(file) },
        ]);
      } else {
        setError("Upload failed — is storage configured?");
      }
    }

    setIsUploading(false);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  function removeMedia(key: string) {
    setMedia((prev) => {
      const target = prev.find((m) => m.storageKey === key);
      if (target) URL.revokeObjectURL(target.previewUrl);
      return prev.filter((m) => m.storageKey !== key);
    });
  }

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
        media: media.map(({ storageKey, fileSize, mimeType }) => ({
          storageKey,
          fileSize,
          mimeType,
        })),
      },
    }).catch(() => null);

    setIsSubmitting(false);

    if (res && res.status === 201) {
      media.forEach((m) => URL.revokeObjectURL(m.previewUrl));
      setIsOpen(false);
      setTitle("");
      setContent("");
      setMedia([]);
      if (!fixedCommunityId) setCommunityId("");
      router.refresh();
    } else {
      setError(
        "Couldn't create the post. If it's in a community, you need to be a member first."
      );
    }
  }

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="border border-dashed border-accent/60 px-4 py-2 font-mono text-[0.72rem] font-bold tracking-[0.06em] text-accent-strong uppercase hover:bg-accent/10"
      >
        + New post
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-30 flex items-start justify-center overflow-y-auto bg-black/70 p-4 py-10">
          <div className="flex w-full max-w-[520px] flex-col gap-4 border border-border bg-surface-raised p-5">
            <div className="flex items-center justify-between">
              <h3 className="font-mono text-[0.8rem] font-bold tracking-[0.08em] uppercase">
                New post
              </h3>
              <button
                onClick={() => setIsOpen(false)}
                className="font-mono text-[0.8rem] text-text-faint hover:text-text"
              >
                ✕
              </button>
            </div>

            {!fixedCommunityId && (
              <select
                value={communityId}
                onChange={(e) => setCommunityId(e.target.value)}
                className="border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
              >
                <option value="">General (no community)</option>
                {communities.map((c) => (
                  <option key={c.communityId} value={c.communityId}>
                    n/{c.slug}
                  </option>
                ))}
              </select>
            )}

            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Title"
              className="border border-border bg-surface px-3.5 py-2.5 text-[0.9rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
            />

            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="What's on your mind? (optional)"
              rows={5}
              className="resize-none border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
            />

            {media.length > 0 && (
              <div className="grid grid-cols-3 gap-2">
                {media.map((m) => (
                  <div key={m.storageKey} className="relative border border-border">
                    {m.mimeType.startsWith("video/") ? (
                      <video src={m.previewUrl} className="h-20 w-full object-cover" muted />
                    ) : (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img src={m.previewUrl} alt="" className="h-20 w-full object-cover" />
                    )}
                    <button
                      onClick={() => removeMedia(m.storageKey)}
                      aria-label={`Remove ${m.name}`}
                      className="absolute top-0 right-0 bg-black/75 px-1.5 py-0.5 font-mono text-[0.7rem] text-white hover:text-accent-strong"
                    >
                      ✕
                    </button>
                  </div>
                ))}
              </div>
            )}

            <div className="flex items-center gap-3">
              <button
                onClick={() => fileInputRef.current?.click()}
                disabled={isUploading || media.length >= 8}
                className="border border-border px-3 py-2 font-mono text-[0.7rem] tracking-[0.06em] text-text-muted uppercase hover:border-accent/40 hover:text-text disabled:opacity-50"
              >
                {isUploading ? "Uploading…" : "+ Photo / video"}
              </button>
              <span className="font-mono text-[0.66rem] text-text-faint">{media.length}/8</span>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*,video/*"
                multiple
                className="hidden"
                onChange={(e) => handleFiles(e.target.files)}
              />
            </div>

            {error && <p className="font-mono text-[0.72rem] text-accent-strong">{error}</p>}

            <div className="flex justify-end gap-2">
              <button
                onClick={() => setIsOpen(false)}
                className="px-4 py-2 font-mono text-[0.72rem] tracking-[0.06em] text-text-muted uppercase hover:text-text"
              >
                Cancel
              </button>
              <button
                onClick={submit}
                disabled={isSubmitting || isUploading || !title.trim()}
                className="bg-accent px-5 py-2 font-mono text-[0.72rem] font-bold tracking-[0.06em] text-accent-text uppercase hover:bg-accent-strong disabled:opacity-50"
              >
                {isSubmitting ? "Posting…" : "Post"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
