import { initClient } from "@ts-rest/core";
import { apiContract } from "./contracts/index.js";

export function createApiClient(baseUrl: string) {
  return initClient(apiContract, { baseUrl });
}
