import { z } from "zod";

export const ZOrder = z.enum(["asc", "desc"]);

export type PaginatedResponse<T> = {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
};

export const schemaWithPagination = <T>(
  schema: z.ZodSchema<T>
): z.ZodSchema<PaginatedResponse<T>> =>
  z.object({
    data: z.array(schema),
    total: z.number(),
    page: z.number(),
    limit: z.number(),
    totalPages: z.number(),
  });

export type CursorPaginatedResponse<T> = {
  data: T[];
  cursorSortValue: string;
  cursorCreatedAt: string;
  hasMore: boolean;
};

export const schemaWithCursorPagination = <T>(
  schema: z.ZodSchema<T>
): z.ZodSchema<CursorPaginatedResponse<T>> =>
  z.object({
    data: z.array(schema),
    cursorSortValue: z.string(),
    cursorCreatedAt: z.string().datetime(),
    hasMore: z.boolean(),
  });