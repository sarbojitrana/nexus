"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import { RemoteAvatar } from "@/components/media/remote-image";
import { useChatSocket, type WsEvent } from "@/hooks/use-chat-socket";
import { formatTimeAgo } from "@/lib/format";
import type { Notification, User } from "@nexus/zod";

export function NotificationsApp() {
  const api = useApi();
  const router = useRouter();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [actors, setActors] = useState<Map<string, User>>(new Map());
  const [isLoading, setIsLoading] = useState(true);

  const enrichActors = useCallback(
    async (list: Notification[]) => {
      const ids = [...new Set(list.map((n) => n.actorId).filter((id): id is string => !!id))].filter(
        (id) => !actors.has(id)
      );
      if (ids.length === 0) return;
      const results = await Promise.all(
        ids.map((id) => api.User.getUserById({ params: { id } }).catch(() => null))
      );
      setActors((prev) => {
        const next = new Map(prev);
        results.forEach((r) => {
          if (r && r.status === 200) next.set(r.body.id, r.body);
        });
        return next;
      });
    },
    [api, actors]
  );

  const load = useCallback(async () => {
    const res = await api.Notification.list({ query: {} }).catch(() => null);
    if (res && res.status === 200) {
      setNotifications(res.body.data);
      enrichActors(res.body.data);
    }
    setIsLoading(false);
  }, [api]);

  useEffect(() => {
    load();
  }, [load]);

  const handleWsEvent = useCallback(
    (event: WsEvent) => {
      if (event.type !== "notification") return;
      const notification = event.payload as Notification;
      setNotifications((prev) => [notification, ...prev]);
      enrichActors([notification]);
    },
    [enrichActors]
  );

  useChatSocket(handleWsEvent);

  async function handleClick(n: Notification) {
    if (!n.isRead) {
      await api.Notification.markRead({ params: { id: n.id } }).catch(() => {});
      setNotifications((prev) => prev.map((x) => (x.id === n.id ? { ...x, isRead: true } : x)));
    }
    const conversationId = n.data.conversationId;
    if (typeof conversationId === "string") {
      router.push(`/dashboard/messages?conversation=${conversationId}`);
    }
  }

  async function markAllRead() {
    await api.Notification.markAllRead().catch(() => {});
    setNotifications((prev) => prev.map((n) => ({ ...n, isRead: true })));
  }

  return (
    <div className="mx-auto flex min-h-0 w-full max-w-[820px] flex-col gap-4 overflow-y-auto px-6 py-8">
      <div className="flex items-center justify-between">
        <h1 className="font-display text-[1.4rem] font-extrabold">Notifications</h1>
        {notifications.some((n) => !n.isRead) && (
          <button
            onClick={markAllRead}
            className="text-[0.78rem] font-bold text-accent-strong hover:underline"
          >
            Mark all read
          </button>
        )}
      </div>

      {isLoading && <p className="text-[0.84rem] text-text-faint">Loading...</p>}

      <div className="flex flex-col gap-1">
        {notifications.map((n) => (
          <button
            key={n.id}
            onClick={() => handleClick(n)}
            className={`flex items-start gap-3 rounded-2xl border px-4 py-3 text-left ${
              n.isRead ? "border-border bg-surface" : "border-accent/30 bg-accent/5"
            }`}
          >
            <RemoteAvatar
              storageKey={actors.get(n.actorId ?? "")?.avatarKey}
              url={actors.get(n.actorId ?? "")?.avatarUrl}
              size={32}
              className="mt-0.5"
            />
            <div className="flex min-w-0 flex-col gap-0.5">
              <span className="text-[0.84rem] text-text">{describe(n, actors.get(n.actorId ?? ""))}</span>
              <span className="text-[0.7rem] text-text-faint">{formatTimeAgo(n.createdAt)}</span>
            </div>
            {!n.isRead && <span className="mt-1.5 ml-auto h-2 w-2 shrink-0 rounded-full bg-accent" />}
          </button>
        ))}
        {!isLoading && notifications.length === 0 && (
          <div className="rounded-2xl border border-border bg-surface p-6 text-center text-[0.86rem] text-text-muted">
            No notifications yet.
          </div>
        )}
      </div>
    </div>
  );
}

function describe(n: Notification, actor: User | undefined): string {
  const who = actor ? `@${actor.username}` : "someone";
  switch (n.type) {
    case "follow":
      return `${who} followed you`;
    case "group_invitation":
      return `${who} invited you to a group`;
    case "message":
      return `${who} sent you a message`;
    default:
      return "New notification";
  }
}
