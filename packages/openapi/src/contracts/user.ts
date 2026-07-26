import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
  ZUser,
  ZUpdateUserPayload,
  ZGetUsersQuery,
  ZGetPostsByUserIDQuery,
  ZPostCommentsResponse,
  schemaWithCursorPagination,
  ZMiniUser,
  ZUserFollow,
  ZGetFollowersQuery,
  ZGetFollowersResponse,
} from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const userContract = c.router(
  {
    getUsers: {
      summary: "Search / browse users",
      method: "GET",
      path: "/users",
      query: ZGetUsersQuery,
      responses: {
        200: schemaWithCursorPagination(ZMiniUser),
      },
    },
    getUserById: {
      summary: "Get a user's profile",
      method: "GET",
      path: "/users/:id",
      pathParams: z.object({ id: z.string() }),
      responses: {
        200: ZUser,
      },
    },
    getUserPosts: {
      summary: "Get a user's posts",
      method: "GET",
      path: "/users/:id/posts",
      pathParams: z.object({ id: z.string() }),
      query: ZGetPostsByUserIDQuery,
      responses: {
        200: ZPostCommentsResponse,
      },
    },
    getMe: {
      summary: "Get the authenticated user's own profile",
      method: "GET",
      path: "/users/me",
      responses: {
        200: ZUser,
      },
      metadata: getSecurityMetadata(),
    },
    updateMe: {
      summary: "Update the authenticated user's profile",
      method: "PATCH",
      path: "/users/me",
      body: ZUpdateUserPayload,
      responses: {
        200: ZUser,
      },
      metadata: getSecurityMetadata(),
    },
    deleteMe: {
      summary: "Delete the authenticated user's account",
      method: "DELETE",
      path: "/users/me",
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    followUser: {
      summary: "Follow a user",
      method: "POST",
      path: "/users/:id/follow",
      pathParams: z.object({ id: z.string() }),
      body: c.noBody(),
      responses: {
        201: ZUserFollow,
      },
      metadata: getSecurityMetadata(),
    },
    unfollowUser: {
      summary: "Unfollow a user",
      method: "DELETE",
      path: "/users/:id/follow",
      pathParams: z.object({ id: z.string() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    isFollowingUser: {
      summary: "Check if the authenticated user follows this user",
      method: "GET",
      path: "/users/:id/follow",
      pathParams: z.object({ id: z.string() }),
      responses: {
        204: c.noBody(),
        404: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    getFollowers: {
      summary: "Get a user's followers",
      method: "GET",
      path: "/users/:id/followers",
      pathParams: z.object({ id: z.string() }),
      query: ZGetFollowersQuery,
      responses: {
        200: ZGetFollowersResponse,
      },
    },
    getFollowing: {
      summary: "Get who a user is following",
      method: "GET",
      path: "/users/:id/following",
      pathParams: z.object({ id: z.string() }),
      query: ZGetFollowersQuery,
      responses: {
        200: ZGetFollowersResponse,
      },
    },
  },
  { pathPrefix: "/api/v1" }
);
