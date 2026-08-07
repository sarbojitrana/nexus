"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import type { MiniUser } from "@nexus/zod";

export function InviteToGroupModal({
  conversationId,
  onClose,
}: {
  conversationId: string;
  onClose: () => void;
}) {
  const api = useApi();
  const [query, setQuery] = useState("");
  const [candidates, setCandidates] = useState<MiniUser[]>([]);
  const [invitedIds, setInvitedIds] = useState<Set<string>>(new Set());
  const [pendingId, setPendingId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const trimmed = query.trim();
    if (!trimmed) {
      setCandidates([]);
      return;
    }
    const timer = setTimeout(async () => {
      const res = await api.User.getUsers({ query: { name: trimmed } }).catch(() => null);
      if (res && res.status === 200) setCandidates(res.body.data);
    }, 250);
    return () => clearTimeout(timer);
  }, [query, api]);

  async function invite(userId: string) {
    setPendingId(userId);
    setError(null);
    const res = await api.Chat.inviteToConversation({
      params: { id: conversationId },
      body: { userId },
    }).catch(() => null);
    setPendingId(null);
    if (res && res.status === 201) {
      setInvitedIds((prev) => new Set(prev).add(userId));
    } else {
      setError("Couldn't invite them — their privacy settings may not allow it.");
    }
  }

  return (
    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/70 p-4">
      <div className="flex w-full max-w-[420px] flex-col gap-4 border border-border bg-surface-raised p-5">
        <div className="flex items-center justify-between">
          <h3 className="font-mono text-[0.78rem] font-bold tracking-[0.08em] uppercase">
            Invite to group
          </h3>
          <button onClick={onClose} className="font-mono text-[0.8rem] text-text-faint hover:text-text">
            ✕
          </button>
        </div>

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search people…"
          autoFocus
          className="border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
        />

        <div className="flex max-h-[260px] flex-col overflow-y-auto">
          {candidates.map((u) => (
            <div
              key={u.id}
              className="flex items-center justify-between px-1 py-2 text-[0.83rem]"
            >
              <span>
                <strong>@{u.username}</strong>{" "}
                <span className="text-text-faint">{u.displayName}</span>
              </span>
              <button
                onClick={() => invite(u.id)}
                disabled={invitedIds.has(u.id) || pendingId === u.id}
                className="border border-border px-2.5 py-1 font-mono text-[0.64rem] tracking-[0.05em] text-text-muted uppercase hover:border-accent hover:text-accent-strong disabled:opacity-50"
              >
                {invitedIds.has(u.id) ? "Invited" : pendingId === u.id ? "…" : "Invite"}
              </button>
            </div>
          ))}
          {query.trim() && candidates.length === 0 && (
            <span className="px-1 py-2 font-mono text-[0.72rem] text-text-faint">
              Nothing to show
            </span>
          )}
        </div>

        {error && <p className="font-mono text-[0.72rem] text-accent-strong">{error}</p>}
      </div>
    </div>
  );
}
