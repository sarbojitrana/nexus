"use client";

import { useMemo } from "react";
import { useAuth } from "@clerk/nextjs";
import { createApiClient } from "@nexus/openapi/client";

const baseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export function useApi() {
  const { getToken } = useAuth();
  return useMemo(() => createApiClient(baseUrl, getToken), [getToken]);
}
