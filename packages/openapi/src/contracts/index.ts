import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { userContract } from "./user.js";
import { postContract } from "./post.js";
import { communityContract } from "./community.js";
import { chatContract } from "./chat.js";
import { notificationContract } from "./notification.js";
import { searchContract } from "./search.js";
import { storageContract } from "./storage.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  User: userContract,
  Post: postContract,
  Community: communityContract,
  Chat: chatContract,
  Notification: notificationContract,
  Search: searchContract,
  Storage: storageContract,
});