"use client";

import { useEffect, useRef, useState } from "react";
import { useAuth } from "@clerk/nextjs";

export type WsEvent =
  | { type: "message"; payload: unknown }
  | { type: "presence"; payload: { userId: string; online: boolean } }
  | { type: "notification"; payload: unknown };

function wsBaseUrl() {
  const httpUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
  return httpUrl.replace(/^http/, "ws");
}

export function useChatSocket(onEvent: (event: WsEvent) => void) {
  const { getToken, isSignedIn } = useAuth();
  const [isConnected, setIsConnected] = useState(false);
  const onEventRef = useRef(onEvent);
  onEventRef.current = onEvent;

  useEffect(() => {
    if (!isSignedIn) return;

    let socket: WebSocket | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    let cancelled = false;
    let attempt = 0;

    async function connect() {
      const token = await getToken();
      if (cancelled || !token) return;

      socket = new WebSocket(`${wsBaseUrl()}/ws/chat?token=${encodeURIComponent(token)}`);

      socket.onopen = () => {
        attempt = 0;
        setIsConnected(true);
      };
      socket.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data) as WsEvent;
          onEventRef.current(parsed);
        } catch {
        }
      };
      socket.onclose = () => {
        setIsConnected(false);
        if (cancelled) return;
        const delay = Math.min(1000 * 2 ** attempt, 15000);
        attempt += 1;
        reconnectTimer = setTimeout(connect, delay);
      };
      socket.onerror = () => {
        socket?.close();
      };
    }

    connect();

    return () => {
      cancelled = true;
      if (reconnectTimer) clearTimeout(reconnectTimer);
      socket?.close();
    };
  }, [isSignedIn, getToken]);

  return { isConnected };
}
