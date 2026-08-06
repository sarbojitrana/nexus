import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
  ZConversationSummary,
  ZStartDirectConversationPayload,
  ZCreateGroupConversationPayload,
  ZInviteToConversationPayload,
  ZInvitation,
  ZSendMessagePayload,
  ZMessage,
  ZGetMessagesQuery,
  ZGetMessagesResponse,
} from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const chatContract = c.router(
  {
    startDirectConversation: {
      summary: "Start (or resume) a 1:1 conversation with another user",
      method: "POST",
      path: "/conversations/direct",
      body: ZStartDirectConversationPayload,
      responses: {
        200: ZConversationSummary,
      },
      metadata: getSecurityMetadata(),
    },
    createGroup: {
      summary: "Create a group conversation and invite members",
      method: "POST",
      path: "/conversations/group",
      body: ZCreateGroupConversationPayload,
      responses: {
        201: ZConversationSummary,
      },
      metadata: getSecurityMetadata(),
    },
    listConversations: {
      summary: "List the authenticated user's conversations",
      method: "GET",
      path: "/conversations",
      responses: {
        200: z.array(ZConversationSummary),
      },
      metadata: getSecurityMetadata(),
    },
    inviteToConversation: {
      summary: "Invite another user to a group conversation",
      method: "POST",
      path: "/conversations/:id/invitations",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZInviteToConversationPayload,
      responses: {
        201: ZInvitation,
      },
      metadata: getSecurityMetadata(),
    },
    getMessages: {
      summary: "Get a conversation's message history",
      method: "GET",
      path: "/conversations/:id/messages",
      pathParams: z.object({ id: z.string().uuid() }),
      query: ZGetMessagesQuery,
      responses: {
        200: ZGetMessagesResponse,
      },
      metadata: getSecurityMetadata(),
    },
    sendMessage: {
      summary: "Send a message in a conversation",
      method: "POST",
      path: "/conversations/:id/messages",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZSendMessagePayload,
      responses: {
        201: ZMessage,
      },
      metadata: getSecurityMetadata(),
    },
    markConversationRead: {
      summary: "Mark a conversation as read",
      method: "POST",
      path: "/conversations/:id/read",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    getPendingInvitations: {
      summary: "List the authenticated user's pending group invitations",
      method: "GET",
      path: "/invitations",
      responses: {
        200: z.array(ZInvitation),
      },
      metadata: getSecurityMetadata(),
    },
    acceptInvitation: {
      summary: "Accept a group invitation",
      method: "POST",
      path: "/invitations/:id/accept",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    rejectInvitation: {
      summary: "Reject a group invitation",
      method: "POST",
      path: "/invitations/:id/reject",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
  },
  { pathPrefix: "/api/v1" }
);
