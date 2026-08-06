import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

extendZodWithOpenApi(z);

export * from "./utils.js";
export * from "./health.js";
export * from "./user.js";
export * from "./post.js";
export * from "./community.js";
export * from "./follow.js";
export * from "./chat.js";
export * from "./notification.js";
export * from "./search.js";
export * from "./storage.js";