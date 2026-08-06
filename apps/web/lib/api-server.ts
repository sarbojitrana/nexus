import { auth } from "@clerk/nextjs/server";
import { createApiClient } from "@nexus/openapi/client";

const baseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function getServerApi() {
  const { getToken } = await auth();
  return createApiClient(baseUrl, getToken);
}

export async function getServerUserId() {
  const { userId } = await auth();
  return userId;
}
