import { initContract } from "@ts-rest/core";
import { ZSearchQuery, ZSearchResults } from "@nexus/zod";

const c = initContract();

export const searchContract = c.router(
  {
    search: {
      summary: "Search posts, people, and communities",
      method: "GET",
      path: "/search",
      query: ZSearchQuery,
      responses: {
        200: ZSearchResults,
      },
    },
  },
  { pathPrefix: "/api/v1" }
);
