"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import type { MiniUser } from "@nexus/zod";

export function NewGroupModal({
  onClose,
  onCreated,
}: {
  onClose: () => void;
  onCreated: (conversationId: string) => void;
}) {
  const api = useApi();
  const [name, setName] = useState("");
  const [query, setQuery] = useState("");
  const [candidates, setCandidates] = useState<MiniUser[]>([]);
  const [selected, setSelected] = useState<Map<string, MiniUser>>(new Map());
  const [isSubmitting, setIsSubmitting] = useState(false);
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

  function toggle(u: MiniUser) {
    setSelected((prev) => {
      const next = new Map(prev);
      if (next.has(u.id)) next.delete(u.id);
      else next.set(u.id, u);
      return next;
    });
  }

  async function submit() {
    if (!name.trim() || selected.size === 0) return;
    setIsSubmitting(true);
    setError(null);
    const res = await api.Chat.createGroup({
      body: { name: name.trim(), inviteeIds: [...selected.keys()] },
    }).catch(() => null);
    setIsSubmitting(false);
    if (res && res.status === 201) {
      onCreated(res.body.id);
    } else {
      setError("Couldn't create the group. Try again.");
    }
  }

  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 p-4">
      <div className="flex w-full max-w-[420px] flex-col gap-4 rounded-2xl border border-border bg-surface-raised p-5">
        <h3 className="text-[1rem] font-bold">New group</h3>

        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Group name"
          className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />

        <div className="flex flex-wrap gap-1.5">
          {[...selected.values()].map((u) => (
            <span
              key={u.id}
              className="flex items-center gap-1.5 rounded-full bg-accent/10 px-2.5 py-1 text-[0.72rem] font-bold text-accent-strong"
            >
              @{u.username}
              <button onClick={() => toggle(u)} className="opacity-70 hover:opacity-100">
                ×
              </button>
            </span>
          ))}
        </div>

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search people to invite..."
          className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />

        <div className="flex max-h-[200px] flex-col gap-0.5 overflow-y-auto">
          {candidates.map((u) => (
            <button
              key={u.id}
              onClick={() => toggle(u)}
              className={`flex items-center justify-between rounded-[9px] px-3 py-2 text-left text-[0.83rem] ${
                selected.has(u.id) ? "bg-accent/10 text-accent-strong" : "text-text hover:bg-surface"
              }`}
            >
              <span>
                <strong>@{u.username}</strong>{" "}
                <span className="text-text-faint">{u.displayName}</span>
              </span>
              {selected.has(u.id) && <span>✓</span>}
            </button>
          ))}
        </div>

        {error && <p className="text-[0.78rem] text-accent-strong">{error}</p>}

        <div className="flex justify-end gap-2">
          <button
            onClick={onClose}
            className="rounded-[9px] px-4 py-2 text-[0.82rem] font-bold text-text-muted hover:bg-surface"
          >
            Cancel
          </button>
          <button
            onClick={submit}
            disabled={isSubmitting || !name.trim() || selected.size === 0}
            className="rounded-[9px] bg-accent px-4 py-2 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
          >
            {isSubmitting ? "Creating..." : "Create group"}
          </button>
        </div>
      </div>
    </div>
  );
}
