"use client";

import { useCallback, useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import { useChatSocket, type WsEvent } from "@/hooks/use-chat-socket";

/** Pink square beside Notifications / Messages when either has something
 *  unread. Both counts come from endpoints the app already calls, and the
 *  existing chat socket keeps them live without extra polling. */
export function NavUnread({ kind }: { kind: "notifications" | "messages" }) {
  const api = useApi();
  const [hasUnread, setHasUnread] = useState(false);

  const check = useCallback(async () => {
    if (kind === "notifications") {
      const res = await api.Notification.list({ query: {} }).catch(() => null);
      if (res && res.status === 200) setHasUnread(res.body.data.some((n) => !n.isRead));
      return;
    }
    const res = await api.Chat.listConversations().catch(() => null);
    if (res && res.status === 200) setHasUnread(res.body.some((c) => c.unreadCount > 0));
  }, [api, kind]);

  useEffect(() => {
    check();
  }, [check]);

  const onEvent = useCallback(
    (event: WsEvent) => {
      if (kind === "notifications" && event.type === "notification") setHasUnread(true);
      if (kind === "messages" && event.type === "message") check();
    },
    [kind, check]
  );

  useChatSocket(onEvent);

  if (!hasUnread) return null;

  return (
    <span
      aria-label="unread"
      className="ml-auto h-2 w-2 shrink-0 bg-accent"
      title={kind === "messages" ? "Unread messages" : "Unread notifications"}
    />
  );
}
