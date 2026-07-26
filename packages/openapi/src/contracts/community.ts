import { initContract } from "@ts-rest/core";
import { z } from "zod";
import {
  ZCommunity,
  ZCommunityResponse,
  ZCommunityMember,
  ZCommunityReport,
  ZBannedFromCommunityUser,
  ZCreateCommunityPayload,
  ZUpdateCommunitySettingsPayload,
  ZChangeMemberRoleInCommunityPayload,
  ZGetCommunitiesQuery,
  ZGetCommunitiesResponse,
  ZGetCommunityMembersQuery,
  ZGetCommunityMembersResponse,
  ZReportCommunityPostPayload,
  ZResolveCommunityPostReportPayload,
  ZBanCommunityMemberPayload,
  ZGetCommunityReportsQuery,
  ZGetCommunityReportsResponse,
  ZPopulatedPost,
  ZCommunityFollow,
} from "@nexus/zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const communityContract = c.router(
  {
    createCommunity: {
      summary: "Create a community",
      method: "POST",
      path: "/communities",
      body: ZCreateCommunityPayload,
      responses: {
        201: ZCommunity,
      },
      metadata: getSecurityMetadata(),
    },
    getCommunities: {
      summary: "Search / browse communities",
      method: "GET",
      path: "/communities",
      query: ZGetCommunitiesQuery,
      responses: {
        200: ZGetCommunitiesResponse,
      },
    },
    getCommunityByIdOrSlug: {
      summary: "Get a community by id or slug",
      method: "GET",
      path: "/communities/:idOrSlug",
      pathParams: z.object({ idOrSlug: z.string() }),
      responses: {
        200: ZCommunityResponse,
      },
    },
    updateCommunitySettings: {
      summary: "Update community settings",
      method: "PATCH",
      path: "/communities/:id",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZUpdateCommunitySettingsPayload,
      responses: {
        200: ZCommunity,
      },
      metadata: getSecurityMetadata(),
    },
    deleteCommunity: {
      summary: "Delete a community",
      method: "DELETE",
      path: "/communities/:id",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    joinCommunity: {
      summary: "Join a community",
      method: "POST",
      path: "/communities/:id/join",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        201: ZCommunityMember,
      },
      metadata: getSecurityMetadata(),
    },
    leaveCommunity: {
      summary: "Leave a community",
      method: "POST",
      path: "/communities/:id/leave",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    followCommunity: {
      summary: "Follow a community without joining it",
      method: "POST",
      path: "/communities/:id/follow",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        201: ZCommunityFollow,
      },
      metadata: getSecurityMetadata(),
    },
    unfollowCommunity: {
      summary: "Unfollow a community",
      method: "DELETE",
      path: "/communities/:id/follow",
      pathParams: z.object({ id: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    isFollowingCommunity: {
      summary: "Check if the authenticated user follows this community",
      method: "GET",
      path: "/communities/:id/follow",
      pathParams: z.object({ id: z.string().uuid() }),
      responses: {
        204: c.noBody(),
        404: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    getCommunityMembers: {
      summary: "Get a community's members",
      method: "GET",
      path: "/communities/:id/members",
      pathParams: z.object({ id: z.string().uuid() }),
      query: ZGetCommunityMembersQuery,
      responses: {
        200: ZGetCommunityMembersResponse,
      },
    },
    changeMemberRole: {
      summary: "Change a community member's role",
      method: "PATCH",
      path: "/communities/:id/members/:userId/role",
      pathParams: z.object({ id: z.string().uuid(), userId: z.string() }),
      body: ZChangeMemberRoleInCommunityPayload,
      responses: {
        200: ZCommunityMember,
      },
      metadata: getSecurityMetadata(),
    },
    getCommunityPostById: {
      summary: "Get a post within a community",
      method: "GET",
      path: "/communities/:id/posts/:postId",
      pathParams: z.object({ id: z.string().uuid(), postId: z.string().uuid() }),
      responses: {
        200: ZPopulatedPost,
      },
    },
    deleteCommunityPost: {
      summary: "Delete a post within a community (moderator)",
      method: "DELETE",
      path: "/communities/:id/posts/:postId",
      pathParams: z.object({ id: z.string().uuid(), postId: z.string().uuid() }),
      body: c.noBody(),
      responses: {
        204: c.noBody(),
      },
      metadata: getSecurityMetadata(),
    },
    banMember: {
      summary: "Ban a member from a community (moderator)",
      method: "POST",
      path: "/communities/:id/bans",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZBanCommunityMemberPayload,
      responses: {
        201: ZBannedFromCommunityUser,
      },
      metadata: getSecurityMetadata(),
    },
    reportPost: {
      summary: "Report a post within a community",
      method: "POST",
      path: "/communities/:id/reports",
      pathParams: z.object({ id: z.string().uuid() }),
      body: ZReportCommunityPostPayload,
      responses: {
        201: ZCommunityReport,
      },
      metadata: getSecurityMetadata(),
    },
    getCommunityReports: {
      summary: "Get a community's reports (moderator)",
      method: "GET",
      path: "/communities/:id/reports",
      pathParams: z.object({ id: z.string().uuid() }),
      query: ZGetCommunityReportsQuery,
      responses: {
        200: ZGetCommunityReportsResponse,
      },
      metadata: getSecurityMetadata(),
    },
    getReportById: {
      summary: "Get a single report (moderator)",
      method: "GET",
      path: "/communities/:id/reports/:reportId",
      pathParams: z.object({ id: z.string().uuid(), reportId: z.string().uuid() }),
      responses: {
        200: ZCommunityReport,
      },
      metadata: getSecurityMetadata(),
    },
    resolveReport: {
      summary: "Resolve or dismiss a report (moderator)",
      method: "PATCH",
      path: "/communities/:id/reports/:reportId",
      pathParams: z.object({ id: z.string().uuid(), reportId: z.string().uuid() }),
      body: ZResolveCommunityPostReportPayload,
      responses: {
        200: ZCommunityReport,
      },
      metadata: getSecurityMetadata(),
    },
  },
  { pathPrefix: "/api/v1" }
);
