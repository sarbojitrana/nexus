import { initContract } from "@ts-rest/core";
import { z } from "zod";
import { ZGetNotificationsQuery, ZGetNotificationsResponse } from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const notificationContract = c.router(
  {
    list: {
      summary: "List the authenticated user's notifications",
      method: "GET",
      path: "/notifications",
      query: ZGetNotificationsQuery,
      responses: {
        200: ZGetNotificationsResponse,
      },
      metadata: getSecurityMetadata(),
    },
    markRead: {
      summary: "Mark a notification as read",
      method: "POST",
      path: "/notifications/:id/read",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    markAllRead: {
      summary: "Mark all notifications as read",
      method: "POST",
      path: "/notifications/read-all",
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
  },
  { pathPrefix: "/api/v1" }
);
