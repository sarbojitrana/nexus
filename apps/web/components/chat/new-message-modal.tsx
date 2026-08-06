"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import type { MiniUser } from "@nexus/zod";

export function NewMessageModal({
  onClose,
  onStarted,
}: {
  onClose: () => void;
  onStarted: (conversationId: string) => void;
}) {
  const api = useApi();
  const [query, setQuery] = useState("");
  const [candidates, setCandidates] = useState<MiniUser[]>([]);
  const [startingId, setStartingId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const trimmed = query.trim();
    if (trimmed.length === 0) {
      setCandidates([]);
      return;
    }
    const timer = setTimeout(async () => {
      const res = await api.User.getUsers({ query: { name: trimmed } }).catch(() => null);
      if (res && res.status === 200) setCandidates(res.body.data);
    }, 250);
    return () => clearTimeout(timer);
  }, [query, api]);

  async function start(userId: string) {
    setStartingId(userId);
    setError(null);
    const res = await api.Chat.startDirectConversation({ body: { userId } }).catch(() => null);
    setStartingId(null);
    if (res && res.status === 200) {
      onStarted(res.body.id);
    } else {
      setError("Couldn't start the conversation. Try again.");
    }
  }

  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 p-4">
      <div className="flex w-full max-w-[420px] flex-col gap-4 rounded-2xl border border-border bg-surface-raised p-5">
        <h3 className="text-[1rem] font-bold">New message</h3>

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search people..."
          autoFocus
          className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />

        <div className="flex max-h-[280px] flex-col gap-0.5 overflow-y-auto">
          {candidates.map((u) => (
            <button
              key={u.id}
              onClick={() => start(u.id)}
              disabled={startingId !== null}
              className="flex items-center justify-between rounded-[9px] px-3 py-2 text-left text-[0.83rem] text-text hover:bg-surface disabled:opacity-50"
            >
              <span>
                <strong>@{u.username}</strong>{" "}
                <span className="text-text-faint">{u.displayName}</span>
              </span>
              {startingId === u.id && <span className="text-text-faint">...</span>}
            </button>
          ))}
          {query.trim().length > 0 && candidates.length === 0 && (
            <span className="px-3 py-2 text-[0.8rem] text-text-faint">No one found</span>
          )}
        </div>

        {error && <p className="text-[0.78rem] text-accent-strong">{error}</p>}

        <div className="flex justify-end">
          <button
            onClick={onClose}
            className="rounded-[9px] px-4 py-2 text-[0.82rem] font-bold text-text-muted hover:bg-surface"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
