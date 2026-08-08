"use client";

import { useCallback, useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useAuth } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { useChatSocket, type WsEvent } from "@/hooks/use-chat-socket";
import { formatTimeAgo } from "@/lib/format";
import { Avatar } from "@/components/chat/avatar";
import { MessageThread } from "@/components/chat/message-thread";
import { NewGroupModal } from "@/components/chat/new-group-modal";
import { NewMessageModal } from "@/components/chat/new-message-modal";
import type { ConversationSummary, Message, Invitation } from "@nexus/zod";

export function ChatApp() {
  const api = useApi();
  const { userId } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();

  const [conversations, setConversations] = useState<ConversationSummary[]>([]);
  const [invitations, setInvitations] = useState<Invitation[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(
    searchParams.get("conversation")
  );
  const [messages, setMessages] = useState<Message[]>([]);
  const [showGroupModal, setShowGroupModal] = useState(false);
  const [showMessageModal, setShowMessageModal] = useState(false);
  const [showInvitations, setShowInvitations] = useState(false);

  const loadConversations = useCallback(async () => {
    const res = await api.Chat.listConversations().catch(() => null);
    if (res && res.status === 200) setConversations(res.body);
  }, [api]);

  const loadInvitations = useCallback(async () => {
    const res = await api.Chat.getPendingInvitations().catch(() => null);
    if (res && res.status === 200) setInvitations(res.body);
  }, [api]);

  useEffect(() => {
    loadConversations();
    loadInvitations();
  }, [loadConversations, loadInvitations]);

  useEffect(() => {
    if (!selectedId) {
      setMessages([]);
      return;
    }
    let cancelled = false;
    (async () => {
      const res = await api.Chat.getMessages({ params: { id: selectedId }, query: {} }).catch(
        () => null
      );
      if (!cancelled && res && res.status === 200) {
        setMessages([...res.body.data].reverse());
      }
      await api.Chat.markConversationRead({ params: { id: selectedId } }).catch(() => {});
      setConversations((prev) =>
        prev.map((c) => (c.id === selectedId ? { ...c, unreadCount: 0 } : c))
      );
    })();
    return () => {
      cancelled = true;
    };
  }, [selectedId, api]);

  const handleWsEvent = useCallback(
    (event: WsEvent) => {
      if (event.type === "message") {
        const message = event.payload as Message;
        if (message.conversationId === selectedId) {
          setMessages((prev) => (prev.some((m) => m.id === message.id) ? prev : [...prev, message]));
        }
        loadConversations();
      } else if (event.type === "presence") {
        const { userId: presenceUserId, online } = event.payload as {
          userId: string;
          online: boolean;
        };
        setConversations((prev) =>
          prev.map((c) => ({
            ...c,
            participants: c.participants.map((p) =>
              p.userId === presenceUserId ? { ...p, isOnline: online } : p
            ),
          }))
        );
      } else if (event.type === "notification") {
        loadInvitations();
      }
    },
    [selectedId, loadConversations, loadInvitations]
  );

  useChatSocket(handleWsEvent);

  function selectConversation(id: string) {
    setSelectedId(id);
    router.replace(`/dashboard/messages?conversation=${id}`);
  }

  async function respondToInvitation(id: string, accept: boolean) {
    const res = await (accept
      ? api.Chat.acceptInvitation({ params: { id } })
      : api.Chat.rejectInvitation({ params: { id } })
    ).catch(() => null);
    if (res && res.status === 204) {
      setInvitations((prev) => prev.filter((inv) => inv.id !== id));
      if (accept) loadConversations();
    }
  }

  const selected = conversations.find((c) => c.id === selectedId) ?? null;

  return (
    <div className="grid h-[calc(100dvh-0px)] grid-cols-1 md:grid-cols-[340px_1fr]">
      <div className="flex flex-col border-r border-border-soft">
        <div className="flex flex-col gap-3 border-b border-border-soft px-4 py-3.5">
          <div className="flex items-center justify-between">
            <strong className="font-display text-[1rem] font-extrabold">Messages</strong>
            <button
              onClick={() => setShowInvitations((v) => !v)}
              className="relative border border-border px-3 py-1.5 font-mono text-[0.68rem] font-bold tracking-[0.06em] text-text-muted uppercase hover:border-accent/40 hover:text-text"
            >
              Invites
              {invitations.length > 0 && (
                <span className="absolute -top-2 -right-2 flex h-4 min-w-4 items-center justify-center bg-accent px-1 font-mono text-[0.6rem] font-bold text-accent-text">
                  {invitations.length}
                </span>
              )}
            </button>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <button
              onClick={() => setShowMessageModal(true)}
              className="border border-border px-3 py-2 font-mono text-[0.68rem] font-bold tracking-[0.06em] text-text-muted uppercase whitespace-nowrap hover:border-accent/40 hover:text-text"
            >
              + Message
            </button>
            <button
              onClick={() => setShowGroupModal(true)}
              className="bg-accent px-3 py-2 font-mono text-[0.68rem] font-bold tracking-[0.06em] text-accent-text uppercase whitespace-nowrap hover:bg-accent-strong"
            >
              + Group
            </button>
          </div>
        </div>

        {showInvitations && (
          <div className="flex flex-col gap-2 border-b border-border-soft px-4 py-3">
            {invitations.length === 0 && (
              <span className="text-[0.78rem] text-text-faint">No pending invitations</span>
            )}
            {invitations.map((inv) => (
              <div key={inv.id} className="flex items-center justify-between gap-2">
                <span className="text-[0.8rem] text-text-muted">Group invitation</span>
                <div className="flex gap-1.5">
                  <button
                    onClick={() => respondToInvitation(inv.id, true)}
                    className="rounded-[9px] bg-up/15 px-2.5 py-1 text-[0.72rem] font-bold text-up"
                  >
                    Accept
                  </button>
                  <button
                    onClick={() => respondToInvitation(inv.id, false)}
                    className="rounded-[9px] bg-surface px-2.5 py-1 text-[0.72rem] font-bold text-text-faint"
                  >
                    Reject
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="flex-1 overflow-y-auto">
          {conversations.map((c) => {
            const other = c.participants.find((p) => p.userId !== userId);
            const title = c.name ?? other?.displayName ?? "Conversation";
            const isOnline = !c.isGroup && (other?.isOnline ?? false);
            return (
              <button
                key={c.id}
                onClick={() => selectConversation(c.id)}
                className={`flex w-full items-center gap-2.5 px-4 py-2.5 text-left hover:bg-surface ${
                  selectedId === c.id ? "bg-surface" : ""
                }`}
              >
                <Avatar online={c.isGroup ? undefined : isOnline} rounded={c.isGroup ? "9px" : "full"} />
                <div className="flex min-w-0 flex-1 flex-col gap-px">
                  <div className="flex items-center justify-between gap-2">
                    <strong className="truncate text-[0.84rem] font-bold">{title}</strong>
                    {c.lastMessage && (
                      <span className="shrink-0 text-[0.66rem] text-text-faint">
                        {formatTimeAgo(c.lastMessage.createdAt)}
                      </span>
                    )}
                  </div>
                  <span className="truncate text-[0.76rem] text-text-faint">
                    {c.lastMessage?.content ?? (c.lastMessage?.attachmentKey ? "📎 attachment" : "No messages yet")}
                  </span>
                </div>
                {c.unreadCount > 0 && (
                  <span className="flex h-5 min-w-5 shrink-0 items-center justify-center rounded-full bg-accent px-1 text-[0.68rem] font-bold text-accent-text">
                    {c.unreadCount}
                  </span>
                )}
              </button>
            );
          })}
          {conversations.length === 0 && (
            <div className="px-4 py-6 text-center text-[0.8rem] text-text-faint">
              No conversations yet. Message someone from search.
            </div>
          )}
        </div>
      </div>

      <div className="hidden md:block">
        {selected ? (
          <MessageThread
            conversation={selected}
            messages={messages}
            onMessageSent={(m) => {
              // The server broadcasts every message back over the socket,
              // including to the sender, so this and the WS handler race --
              // whichever loses would otherwise append a duplicate.
              setMessages((prev) => (prev.some((x) => x.id === m.id) ? prev : [...prev, m]));
              loadConversations();
            }}
          />
        ) : (
          <div className="flex h-full items-center justify-center text-[0.86rem] text-text-faint">
            Select a conversation
          </div>
        )}
      </div>

      {showGroupModal && (
        <NewGroupModal
          onClose={() => setShowGroupModal(false)}
          onCreated={(id) => {
            setShowGroupModal(false);
            loadConversations();
            selectConversation(id);
          }}
        />
      )}

      {showMessageModal && (
        <NewMessageModal
          onClose={() => setShowMessageModal(false)}
          onStarted={(id) => {
            setShowMessageModal(false);
            loadConversations();
            selectConversation(id);
          }}
        />
      )}
    </div>
  );
}
