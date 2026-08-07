"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useApi } from "@/lib/use-api";
import type { SearchResults } from "@nexus/zod";

const FILTERS = ["All", "Posts", "People", "Communities"] as const;
type Filter = (typeof FILTERS)[number];

export function SearchBar() {
  const api = useApi();
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<Filter>("All");
  const [results, setResults] = useState<SearchResults | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const trimmed = query.trim();
    if (trimmed.length === 0) {
      setResults(null);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    const timer = setTimeout(async () => {
      const res = await api.Search.search({ query: { q: trimmed } }).catch(() => null);
      if (res && res.status === 200) setResults(res.body);
      setIsLoading(false);
    }, 300);

    return () => clearTimeout(timer);
  }, [query, api]);

  useEffect(() => {
    function onClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  }, []);

  const showPosts = filter === "All" || filter === "Posts";
  const showPeople = filter === "All" || filter === "People";
  const showCommunities = filter === "All" || filter === "Communities";

  const visibleCount =
    (showPosts ? (results?.posts.length ?? 0) : 0) +
    (showPeople ? (results?.users.length ?? 0) : 0) +
    (showCommunities ? (results?.communities.length ?? 0) : 0);

  return (
    <div ref={containerRef} className="relative">
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onFocus={() => setIsOpen(true)}
        placeholder="Search posts, people, communities…"
        className="w-full border border-border bg-surface px-3.5 py-2.5 font-mono text-[0.8rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
      />

      {isOpen && query.trim().length > 0 && (
        <div className="absolute top-[calc(100%+4px)] left-0 z-20 max-h-[65vh] w-full overflow-y-auto border border-border bg-surface-raised">
          <div className="flex gap-1.5 border-b border-border-soft p-2.5">
            {FILTERS.map((f) => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`border px-2.5 py-1 font-mono text-[0.64rem] tracking-[0.06em] uppercase ${
                  filter === f
                    ? "border-accent text-accent-strong"
                    : "border-border text-text-faint hover:text-text-muted"
                }`}
              >
                {f}
              </button>
            ))}
          </div>

          <div className="p-2">
            {isLoading && (
              <div className="px-3 py-2 font-mono text-[0.72rem] text-text-faint">searching…</div>
            )}

            {!isLoading && visibleCount === 0 && (
              <div className="px-3 py-2 font-mono text-[0.72rem] text-text-faint">
                Nothing to show
              </div>
            )}

            {showPeople && results && results.users.length > 0 && (
              <Group label="People">
                {results.users.map((u) => (
                  <Link
                    key={u.id}
                    href={`/dashboard/profile/${u.id}`}
                    onClick={() => setIsOpen(false)}
                    className="flex items-center justify-between px-3 py-2 text-[0.83rem] hover:bg-surface"
                  >
                    <span>
                      <strong className="font-bold">@{u.username}</strong>{" "}
                      <span className="text-text-faint">{u.displayName}</span>
                    </span>
                    <span className="font-mono text-[0.64rem] text-text-faint">
                      {u.followerCount} followers
                    </span>
                  </Link>
                ))}
              </Group>
            )}

            {showCommunities && results && results.communities.length > 0 && (
              <Group label="Communities">
                {results.communities.map((c) => (
                  <Link
                    key={c.id}
                    href={`/dashboard/communities/${c.slug}`}
                    onClick={() => setIsOpen(false)}
                    className="flex items-center justify-between px-3 py-2 text-[0.83rem] hover:bg-surface"
                  >
                    <strong className="font-bold">n/{c.slug}</strong>
                    <span className="font-mono text-[0.64rem] text-text-faint">
                      {c.membersCount} members
                    </span>
                  </Link>
                ))}
              </Group>
            )}

            {showPosts && results && results.posts.length > 0 && (
              <Group label="Posts">
                {results.posts.map((p) => (
                  <Link
                    key={p.id}
                    href={`/dashboard/posts/${p.id}`}
                    onClick={() => setIsOpen(false)}
                    className="block px-3 py-2 text-[0.83rem] hover:bg-surface"
                  >
                    {p.title ?? p.content?.slice(0, 80) ?? "(untitled)"}
                  </Link>
                ))}
              </Group>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

function Group({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="mb-2 last:mb-0">
      <div className="eyebrow px-3 py-1">{label}</div>
      <div className="flex flex-col">{children}</div>
    </div>
  );
}
