"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import { RemoteAvatar } from "@/components/media/remote-image";
import { ShareButton } from "@/components/share-button";
import { formatTimeAgo } from "@/lib/format";
import type { SearchResults as Results } from "@nexus/zod";

const TABS = ["All", "Posts", "People", "Communities"] as const;
type Tab = (typeof TABS)[number];

export function SearchResults({ initialQuery }: { initialQuery: string }) {
  const api = useApi();
  const router = useRouter();

  const [query, setQuery] = useState(initialQuery);
  const [submitted, setSubmitted] = useState(initialQuery);
  const [tab, setTab] = useState<Tab>("All");
  const [results, setResults] = useState<Results | null>(null);
  const [isLoading, setIsLoading] = useState(!!initialQuery);

  const run = useCallback(
    async (q: string) => {
      if (!q.trim()) {
        setResults(null);
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      const res = await api.Search.search({ query: { q: q.trim() } }).catch(() => null);
      if (res && res.status === 200) setResults(res.body);
      setIsLoading(false);
    },
    [api]
  );

  useEffect(() => {
    run(submitted);
  }, [submitted, run]);

  function submit() {
    const q = query.trim();
    if (!q) return;
    setSubmitted(q);
    // Keep the URL shareable and back-button friendly.
    router.replace(`/dashboard/search?q=${encodeURIComponent(q)}`);
  }

  const showPosts = tab === "All" || tab === "Posts";
  const showPeople = tab === "All" || tab === "People";
  const showCommunities = tab === "All" || tab === "Communities";

  const total =
    (results?.posts.length ?? 0) + (results?.users.length ?? 0) + (results?.communities.length ?? 0);

  return (
    <div className="flex min-h-0 flex-col">
      <div className="shrink-0 space-y-3 px-6 pt-6 pb-3">
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && submit()}
          autoFocus
          placeholder="Search posts, people, communities…"
          className="w-full border border-border bg-surface px-3.5 py-2.5 font-mono text-[0.8rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />

        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border-soft pb-3">
          <div className="flex gap-2">
            {TABS.map((t) => (
              <button
                key={t}
                onClick={() => setTab(t)}
                className={`border px-3 py-1.5 font-mono text-[0.68rem] font-bold tracking-[0.06em] uppercase ${
                  tab === t
                    ? "border-accent text-accent-strong"
                    : "border-border text-text-faint hover:text-text-muted"
                }`}
              >
                {t}
              </button>
            ))}
          </div>
          {submitted && !isLoading && (
            <span className="font-mono text-[0.68rem] text-text-faint">
              {total} result{total === 1 ? "" : "s"} for “{submitted}”
            </span>
          )}
        </div>
      </div>

      <div className="flex min-h-0 flex-1 flex-col gap-5 overflow-y-auto px-6 pb-6">
        {isLoading && <p className="font-mono text-[0.76rem] text-text-faint">searching…</p>}

        {!isLoading && submitted && total === 0 && (
          <div className="border border-border bg-surface p-8 text-center">
            <p className="text-[0.88rem] text-text-muted">No results for “{submitted}”.</p>
            <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
              Try a shorter or differently spelled term.
            </p>
          </div>
        )}

        {!isLoading && showPeople && results && results.users.length > 0 && (
          <section>
            <h2 className="eyebrow mb-2">People</h2>
            <div className="flex flex-col gap-2">
              {results.users.map((u) => (
                <div
                  key={u.id}
                  className="flex items-center gap-3 border border-border bg-surface p-3.5"
                >
                  <RemoteAvatar size={36} />
                  <Link href={`/dashboard/profile/${u.id}`} className="min-w-0 flex-1">
                    <strong className="block truncate text-[0.88rem] font-bold hover:text-accent-strong">
                      @{u.username}
                    </strong>
                    <span className="font-mono text-[0.68rem] text-text-faint">
                      {u.displayName} · {u.followerCount} followers
                    </span>
                  </Link>
                  <ShareButton
                    path={`/dashboard/profile/${u.id}`}
                    className="shrink-0 font-mono text-[0.68rem] text-text-faint"
                  />
                </div>
              ))}
            </div>
          </section>
        )}

        {!isLoading && showCommunities && results && results.communities.length > 0 && (
          <section>
            <h2 className="eyebrow mb-2">Communities</h2>
            <div className="flex flex-col gap-2">
              {results.communities.map((c) => (
                <div
                  key={c.id}
                  className="flex items-center gap-3 border border-border bg-surface p-3.5"
                >
                  <RemoteAvatar size={36} />
                  <Link href={`/dashboard/communities/${c.slug}`} className="min-w-0 flex-1">
                    <strong className="block truncate text-[0.88rem] font-bold hover:text-accent-strong">
                      n/{c.slug}
                    </strong>
                    <span className="font-mono text-[0.68rem] text-text-faint">
                      {c.name} · {c.membersCount} members
                    </span>
                  </Link>
                  <ShareButton
                    path={`/dashboard/communities/${c.slug}`}
                    className="shrink-0 font-mono text-[0.68rem] text-text-faint"
                  />
                </div>
              ))}
            </div>
          </section>
        )}

        {!isLoading && showPosts && results && results.posts.length > 0 && (
          <section>
            <h2 className="eyebrow mb-2">Posts</h2>
            <div className="flex flex-col gap-2">
              {results.posts.map((p) => (
                <div key={p.id} className="border border-border bg-surface p-3.5">
                  <Link href={`/dashboard/posts/${p.id}`}>
                    <h3 className="font-display text-[0.96rem] font-bold hover:text-accent-strong">
                      {p.title ?? p.content?.slice(0, 90) ?? "(untitled)"}
                    </h3>
                  </Link>
                  {p.title && p.content && (
                    <p className="mt-1 line-clamp-2 text-[0.84rem] text-text-muted">{p.content}</p>
                  )}
                  <div className="mt-2 flex gap-4 font-mono text-[0.68rem] text-text-faint">
                    <span>{p.upvotes} points</span>
                    <span>{p.commentCount} comments</span>
                    <span>{formatTimeAgo(p.createdAt)}</span>
                    <ShareButton path={`/dashboard/posts/${p.id}`} />
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  );
}
