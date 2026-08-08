"use client";

import { useEffect, useRef, useState } from "react";
import { useAuth } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { useUpload } from "@/lib/use-upload";
import { formatTimeAgo } from "@/lib/format";
import { InviteToGroupModal } from "@/components/chat/invite-modal";
import { ChatAttachment } from "@/components/chat/chat-attachment";
import type { ConversationSummary, Message } from "@nexus/zod";

export function MessageThread({
  conversation,
  messages,
  onMessageSent,
}: {
  conversation: ConversationSummary;
  messages: Message[];
  onMessageSent: (message: Message) => void;
}) {
  const api = useApi();
  const upload = useUpload();
  const { userId } = useAuth();
  const [text, setText] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [isSending, setIsSending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showInvite, setShowInvite] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: "end" });
  }, [messages.length, conversation.id]);

  const title =
    conversation.name ??
    conversation.participants.find((p) => p.userId !== userId)?.displayName ??
    "Conversation";

  async function send() {
    const trimmed = text.trim();
    if (!trimmed && !file) return;
    setIsSending(true);
    setError(null);

    let attachmentKey: string | undefined;
    let attachmentMimeType: string | undefined;
    let attachmentFileSize: number | undefined;

    if (file) {
      const result = await upload(file);
      if (!result.ok) {
        setError(result.error);
        setIsSending(false);
        return;
      }
      attachmentKey = result.media.storageKey;
      attachmentMimeType = result.media.mimeType;
      attachmentFileSize = result.media.fileSize;
    }

    const res = await api.Chat.sendMessage({
      params: { id: conversation.id },
      body: {
        content: trimmed || null,
        attachmentKey: attachmentKey ?? null,
        attachmentMimeType: attachmentMimeType ?? null,
        attachmentFileSize: attachmentFileSize ?? null,
      },
    }).catch(() => null);

    setIsSending(false);
    if (res && res.status === 201) {
      onMessageSent(res.body);
      setText("");
      setFile(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  }

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center gap-2 border-b border-border-soft px-4 py-3">
        <strong className="text-[0.9rem] font-bold">{title}</strong>
        {conversation.isGroup && (
          <>
            <span className="font-mono text-[0.68rem] text-text-faint">
              {conversation.participants.length} members
            </span>
            <button
              onClick={() => setShowInvite(true)}
              className="ml-auto border border-border px-2.5 py-1 font-mono text-[0.64rem] tracking-[0.05em] text-text-muted uppercase hover:border-accent hover:text-accent-strong"
            >
              + Invite
            </button>
          </>
        )}
      </div>

      <div className="flex flex-1 flex-col gap-2 overflow-y-auto px-4 py-4">
        {messages.map((m) => {
          const mine = m.senderId === userId;
          const sender = conversation.participants.find((p) => p.userId === m.senderId);
          return (
            <div key={m.id} className={`flex flex-col ${mine ? "items-end" : "items-start"}`}>
              {conversation.isGroup && !mine && (
                <span className="mb-0.5 px-1 text-[0.68rem] text-text-faint">
                  {sender?.displayName ?? "unknown"}
                </span>
              )}
              <div
                className={`max-w-[70%] rounded-2xl px-3.5 py-2 text-[0.86rem] ${
                  mine ? "bg-accent text-accent-text" : "bg-surface text-text"
                }`}
              >
                {m.content && <p className="whitespace-pre-wrap">{m.content}</p>}
                {m.attachmentKey && (
                  <ChatAttachment
                    storageKey={m.attachmentKey}
                    mimeType={m.attachmentMimeType}
                    fileSize={m.attachmentFileSize}
                  />
                )}
              </div>
              <span className="mt-0.5 px-1 text-[0.66rem] text-text-faint">
                {formatTimeAgo(m.createdAt)}
              </span>
            </div>
          );
        })}
        {messages.length === 0 && (
          <div className="flex flex-1 items-center justify-center text-[0.82rem] text-text-faint">
            No messages yet. Say hi.
          </div>
        )}
        <div ref={bottomRef} />
      </div>
      {error && (
        <p className="border-t border-border-soft px-3 pt-2 font-mono text-[0.7rem] text-accent-strong">
          {error}
        </p>
      )}

      <div className="flex items-center gap-2 border-t border-border-soft px-3 py-3">
        <button
          onClick={() => fileInputRef.current?.click()}
          className="shrink-0 rounded-[9px] border border-border px-2.5 py-2 text-[0.8rem] text-text-muted hover:bg-surface"
          aria-label="Attach file"
        >
          📎
        </button>
        <input
          ref={fileInputRef}
          type="file"
          className="hidden"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        />
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              send();
            }
          }}
          placeholder={file ? `${file.name} attached · add a note...` : "Message..."}
          className="flex-1 rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />
        <button
          onClick={send}
          disabled={isSending || (!text.trim() && !file)}
          className="shrink-0 rounded-[9px] bg-accent px-4 py-2.5 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
        >
          Send
        </button>
      </div>

      {showInvite && (
        <InviteToGroupModal
          conversationId={conversation.id}
          onClose={() => setShowInvite(false)}
        />
      )}
    </div>
  );
}
