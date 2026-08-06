import { initClient, tsRestFetchApi } from "@ts-rest/core";
import { apiContract } from "./contracts/index.js";

export type GetToken = () => Promise<string | null | undefined>;

export function createApiClient(baseUrl: string, getToken?: GetToken) {
  return initClient(apiContract, {
    baseUrl,
    ...(getToken && {
      api: async (args) => {
        const token = await getToken();
        return tsRestFetchApi({
          ...args,
          headers: {
            ...args.headers,
            ...(token && { Authorization: `Bearer ${token}` }),
          },
        });
      },
    }),
  });
}
