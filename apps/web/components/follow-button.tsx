"use client";

import { useState } from "react";
import { useApi } from "@/lib/use-api";

type FollowButtonProps = {
  kind: "user" | "community";
  id: string;
};

export function FollowButton({ kind, id }: FollowButtonProps) {
  const api = useApi();
  const [state, setState] = useState<"idle" | "loading" | "following">("idle");

  async function handleClick() {
    if (state !== "idle") return;
    setState("loading");
    try {
      if (kind === "user") {
        await api.User.followUser({ params: { id } });
      } else {
        await api.Community.followCommunity({ params: { id } });
      }
      setState("following");
    } catch {
      setState("following");
    }
  }

  return (
    <button
      onClick={handleClick}
      disabled={state !== "idle"}
      className="ml-auto shrink-0 rounded-full border border-border px-2.5 py-1 text-[0.72rem] font-bold text-text-muted hover:border-accent hover:text-accent-strong disabled:opacity-60"
    >
      {state === "following" ? "Following" : "Follow"}
    </button>
  );
}
