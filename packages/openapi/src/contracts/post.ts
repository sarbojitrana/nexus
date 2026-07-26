import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
  ZPost,
  ZPopulatedPost,
  ZCreatePostPayload,
  ZUpdatePostByIDPayload,
  ZGetCommentsByPostIDQuery,
  ZGetRepliesByCommentIDQuery,
  ZReactToPostPayload,
  ZGetPostsQuery,
  ZGetPostsQueryResponse,
  ZPostCommentsResponse,
  ZPostRepliesResponse,
} from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const postContract = c.router(
  {
    createPost: {
      summary: "Create a post or comment",
      method: "POST",
      path: "/posts",
      body: ZCreatePostPayload,
      responses: {
        201: ZPost,
      },
      metadata: getSecurityMetadata(),
    },
    updatePost: {
      summary: "Update a post or comment",
      method: "PATCH",
      path: "/posts/:id",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZUpdatePostByIDPayload,
      responses: {
        200: ZPost,
      },
      metadata: getSecurityMetadata(),
    },
    deletePost: {
      summary: "Delete a post or comment",
      method: "DELETE",
      path: "/posts/:id",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    getPostById: {
      summary: "Get a post or comment",
      method: "GET",
      path: "/posts/:id",
      pathParams: z.object({ id: z.string().uuid() }),
      responses: {
        200: ZPopulatedPost,
      },
    },
    getPostComments: {
      summary: "Get top-level comments on a post",
      method: "GET",
      path: "/posts/:id/comments",
      pathParams: z.object({ id: z.string().uuid() }),
      query: ZGetCommentsByPostIDQuery,
      responses: {
        200: ZPostCommentsResponse,
      },
    },
    getPostReplies: {
      summary: "Get replies to a comment",
      method: "GET",
      path: "/posts/:id/replies",
      pathParams: z.object({ id: z.string().uuid() }),
      query: ZGetRepliesByCommentIDQuery,
      responses: {
        200: ZPostRepliesResponse,
      },
    },
    reactToPost: {
      summary: "Upvote, downvote, or un-vote a post",
      method: "POST",
      path: "/posts/:id/react",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZReactToPostPayload,
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    getFeed: {
      summary: "Get the home feed (trending + following)",
      method: "GET",
      path: "/feed",
      query: ZGetPostsQuery,
      responses: {
        200: ZGetPostsQueryResponse,
      },
    },
  },
  { pathPrefix: "/api/v1" }
);
