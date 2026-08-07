"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useApi } from "@/lib/use-api";
import type { MiniCommunity } from "@nexus/zod";

export function CommunitiesApp() {
  const api = useApi();
  const [query, setQuery] = useState("");
  const [communities, setCommunities] = useState<MiniCommunity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [joinedIds, setJoinedIds] = useState<Set<string>>(new Set());

  const load = useCallback(
    async (name: string) => {
      setIsLoading(true);
      const res = await api.Community.getCommunities({ query: { name: name || undefined } }).catch(
        () => null
      );
      if (res && res.status === 200) setCommunities(res.body.data);
      setIsLoading(false);
    },
    [api]
  );

  useEffect(() => {
    const timer = setTimeout(() => load(query.trim()), 250);
    return () => clearTimeout(timer);
  }, [query, load]);

  async function join(communityId: string) {
    const res = await api.Community.joinCommunity({ params: { id: communityId } }).catch(() => null);
    if (res) setJoinedIds((prev) => new Set(prev).add(communityId));
  }

  return (
    <div className="mx-auto flex w-full max-w-[720px] flex-col gap-4 px-6 py-8">
      <div className="flex items-center justify-between">
        <h1 className="font-display text-[1.4rem] font-extrabold">Communities</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className=" bg-accent px-4 py-2 text-[0.8rem] font-bold text-accent-text hover:bg-accent-strong"
        >
          + Create community
        </button>
      </div>

      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search communities..."
        className=" border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
      />

      <div className="flex flex-col gap-2">
        {isLoading && <p className="text-[0.84rem] text-text-faint">Loading...</p>}

        {!isLoading &&
          communities.map((c) => (
            <div
              key={c.communityId}
              className="flex items-center gap-3 border border-border bg-surface p-4"
            >
              <div className="h-11 w-11 shrink-0 bg-accent" />
              <Link
                href={`/dashboard/communities/${c.slug}`}
                className="flex min-w-0 flex-1 flex-col gap-0.5"
              >
                <strong className="text-[0.9rem] font-bold hover:text-accent-strong">
                  n/{c.slug}
                </strong>
                <span className="font-mono text-[0.7rem] text-text-faint">
                  {c.communityName} · {c.membersCount} members · {c.postsCount} posts
                </span>
              </Link>
              <button
                onClick={() => join(c.communityId)}
                disabled={joinedIds.has(c.communityId)}
                className="shrink-0 border border-border px-4 py-2 font-mono text-[0.7rem] font-bold tracking-[0.05em] text-text-muted uppercase hover:border-accent hover:text-accent-strong disabled:opacity-50"
              >
                {joinedIds.has(c.communityId) ? "Joined" : "Join"}
              </button>
            </div>
          ))}

        {!isLoading && communities.length === 0 && (
          <div className=" border border-border bg-surface p-6 text-center text-[0.86rem] text-text-muted">
            {query.trim() ? "No communities match that search." : "No communities yet — create the first one."}
          </div>
        )}
      </div>

      {showCreateModal && (
        <CreateCommunityModal
          onClose={() => setShowCreateModal(false)}
          onCreated={() => {
            setShowCreateModal(false);
            load(query.trim());
          }}
        />
      )}
    </div>
  );
}

function CreateCommunityModal({
  onClose,
  onCreated,
}: {
  onClose: () => void;
  onCreated: () => void;
}) {
  const api = useApi();
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [slugTouched, setSlugTouched] = useState(false);

  function handleNameChange(value: string) {
    setName(value);
    if (!slugTouched) {
      setSlug(
        value
          .toLowerCase()
          .trim()
          .replace(/[^a-z0-9]+/g, "-")
          .replace(/(^-|-$)/g, "")
      );
    }
  }

  async function submit() {
    if (!name.trim() || !slug.trim()) {
      setError("Name and slug are required.");
      return;
    }
    setIsSubmitting(true);
    setError(null);
    const res = await api.Community.createCommunity({
      body: { name: name.trim(), slug: slug.trim(), description: description.trim() || null },
    }).catch(() => null);
    setIsSubmitting(false);
    if (res && res.status === 201) {
      onCreated();
    } else {
      setError("Couldn't create the community. That slug might already be taken.");
    }
  }

  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 p-4">
      <div className="flex w-full max-w-[440px] flex-col gap-4 border border-border bg-surface-raised p-5">
        <h3 className="text-[1rem] font-bold">Create a community</h3>

        <label className="flex flex-col gap-1.5">
          <span className="text-[0.78rem] font-bold text-text-muted">Name</span>
          <input
            value={name}
            onChange={(e) => handleNameChange(e.target.value)}
            placeholder="Photography"
            className=" border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
        </label>

        <label className="flex flex-col gap-1.5">
          <span className="text-[0.78rem] font-bold text-text-muted">Slug</span>
          <input
            value={slug}
            onChange={(e) => {
              setSlug(e.target.value);
              setSlugTouched(true);
            }}
            placeholder="photography"
            className=" border border-border bg-surface px-3.5 py-2.5 font-mono text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
        </label>

        <label className="flex flex-col gap-1.5">
          <span className="text-[0.78rem] font-bold text-text-muted">Description</span>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            placeholder="What's this community about? (optional)"
            className="resize-none border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
        </label>

        {error && <p className="text-[0.78rem] text-accent-strong">{error}</p>}

        <div className="flex justify-end gap-2">
          <button
            onClick={onClose}
            className=" px-4 py-2 text-[0.82rem] font-bold text-text-muted hover:bg-surface"
          >
            Cancel
          </button>
          <button
            onClick={submit}
            disabled={isSubmitting || !name.trim() || !slug.trim()}
            className=" bg-accent px-4 py-2 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
          >
            {isSubmitting ? "Creating..." : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
