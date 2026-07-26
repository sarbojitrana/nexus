import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { userContract } from "./user.js";
import { postContract } from "./post.js";
import { communityContract } from "./community.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  User: userContract,
  Post: postContract,
  Community: communityContract,
});