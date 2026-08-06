"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import type { SearchResults } from "@nexus/zod";

export function SearchBar() {
  const api = useApi();
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResults | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const trimmed = query.trim();
    if (trimmed.length === 0) {
      setResults(null);
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
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  async function messageUser(userId: string) {
    const res = await api.Chat.startDirectConversation({ body: { userId } }).catch(() => null);
    setIsOpen(false);
    if (res && res.status === 200) {
      router.push(`/dashboard/messages?conversation=${res.body.id}`);
    } else {
      router.push("/dashboard/messages");
    }
  }

  const hasResults =
    results && (results.posts.length > 0 || results.users.length > 0 || results.communities.length > 0);

  return (
    <div ref={containerRef} className="relative">
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onFocus={() => setIsOpen(true)}
        placeholder="Search posts, people, communities..."
        className="w-full rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
      />

      {isOpen && query.trim().length > 0 && (
        <div className="absolute top-[calc(100%+6px)] left-0 z-10 max-h-[60vh] w-full overflow-y-auto rounded-2xl border border-border bg-surface-raised p-2 shadow-xl">
          {isLoading && (
            <div className="px-3 py-2 text-[0.78rem] text-text-faint">Searching...</div>
          )}
          {!isLoading && !hasResults && (
            <div className="px-3 py-2 text-[0.78rem] text-text-faint">No results</div>
          )}

          {results && results.users.length > 0 && (
            <SearchGroup label="People">
              {results.users.map((u) => (
                <button
                  key={u.id}
                  onClick={() => messageUser(u.id)}
                  className="flex w-full items-center justify-between rounded-[9px] px-3 py-2 text-left text-[0.83rem] text-text hover:bg-surface"
                >
                  <span className="flex flex-col">
                    <strong className="font-bold">@{u.username}</strong>
                    <span className="text-[0.72rem] text-text-faint">{u.displayName}</span>
                  </span>
                  <span className="font-mono text-[0.68rem] text-text-faint">Message</span>
                </button>
              ))}
            </SearchGroup>
          )}

          {results && results.communities.length > 0 && (
            <SearchGroup label="Communities">
              {results.communities.map((c) => (
                <div
                  key={c.id}
                  className="flex items-center justify-between rounded-[9px] px-3 py-2 text-[0.83rem] text-text"
                >
                  <strong className="font-bold">n/{c.slug}</strong>
                  <span className="font-mono text-[0.68rem] text-text-faint">
                    {c.membersCount} members
                  </span>
                </div>
              ))}
            </SearchGroup>
          )}

          {results && results.posts.length > 0 && (
            <SearchGroup label="Posts">
              {results.posts.map((p) => (
                <div key={p.id} className="rounded-[9px] px-3 py-2 text-[0.83rem] text-text">
                  {p.title ?? p.content?.slice(0, 80) ?? "(untitled)"}
                </div>
              ))}
            </SearchGroup>
          )}
        </div>
      )}
    </div>
  );
}

function SearchGroup({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="mb-1 last:mb-0">
      <div className="px-3 py-1 font-mono text-[0.66rem] tracking-[0.1em] text-text-faint uppercase">
        {label}
      </div>
      <div className="flex flex-col">{children}</div>
    </div>
  );
}
