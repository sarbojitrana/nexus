import { z } from "zod";
import { schemaWithCursorPagination } from "@/utils.js";

export const ZConversation = z.object({
  id: z.string().uuid(),
  isGroup: z.boolean(),
  name: z.string().nullable(),
  createdBy: z.string(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const ZInvitationStatus = z.enum(["pending", "accepted", "rejected"]);

export const ZInvitation = z.object({
  id: z.string().uuid(),
  conversationId: z.string().uuid(),
  inviterId: z.string(),
  inviteeId: z.string(),
  status: ZInvitationStatus,
  createdAt: z.string().datetime(),
  respondedAt: z.string().datetime().nullable(),
});
export type Invitation = z.infer<typeof ZInvitation>;

export const ZMessage = z.object({
  id: z.string().uuid(),
  conversationId: z.string().uuid(),
  senderId: z.string(),
  content: z.string().nullable(),
  attachmentKey: z.string().nullable(),
  attachmentMimeType: z.string().nullable(),
  attachmentFileSize: z.number().nullable(),
  createdAt: z.string().datetime(),
});
export type Message = z.infer<typeof ZMessage>;

export const ZParticipantView = z.object({
  userId: z.string(),
  username: z.string(),
  displayName: z.string(),
  avatarKey: z.string().nullable(),
  avatarUrl: z.string().nullable(),
  isOnline: z.boolean(),
  lastReadAt: z.string().datetime().optional(),
});
export type ParticipantView = z.infer<typeof ZParticipantView>;

export const ZConversationSummary = ZConversation.extend({
  lastMessage: ZMessage.nullable(),
  unreadCount: z.number(),
  participants: z.array(ZParticipantView),
});
export type ConversationSummary = z.infer<typeof ZConversationSummary>;

export const ZStartDirectConversationPayload = z.object({
  userId: z.string(),
});

export const ZCreateGroupConversationPayload = z.object({
  name: z.string().max(100),
  inviteeIds: z.array(z.string()).min(1).max(50),
});

export const ZInviteToConversationPayload = z.object({
  userId: z.string(),
});

export const ZSendMessagePayload = z
  .object({
    content: z.string().max(4000).nullable().optional(),
    attachmentKey: z.string().nullable().optional(),
    attachmentMimeType: z.string().nullable().optional(),
    attachmentFileSize: z.number().nullable().optional(),
  })
  .refine((v) => !!v.content || !!v.attachmentKey, {
    message: "a message needs content, an attachment, or both",
  });

export const ZGetMessagesQuery = z.object({
  cursorCreatedAt: z.string().datetime().optional(),
});

export const ZGetMessagesResponse = schemaWithCursorPagination(ZMessage);

export const ZWsEvent = z.object({
  type: z.enum(["message", "presence", "notification"]),
  payload: z.unknown(),
});

export const ZWsPresencePayload = z.object({
  userId: z.string(),
  online: z.boolean(),
});
