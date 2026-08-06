import { z } from "zod";
import { schemaWithCursorPagination } from "@/utils.js";

export const ZNotificationType = z.enum(["follow", "group_invitation", "message"]);

export const ZNotification = z.object({
  id: z.string().uuid(),
  userId: z.string(),
  actorId: z.string().nullable(),
  type: ZNotificationType,
  data: z.record(z.string(), z.unknown()),
  isRead: z.boolean(),
  createdAt: z.string().datetime(),
});
export type Notification = z.infer<typeof ZNotification>;

export const ZGetNotificationsQuery = z.object({
  cursorCreatedAt: z.string().datetime().optional(),
});

export const ZGetNotificationsResponse = schemaWithCursorPagination(ZNotification);
