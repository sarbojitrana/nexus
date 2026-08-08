"use client";

import { useCallback, useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";

type FollowButtonProps = {
  kind: "user" | "community";
  id: string;
};

// The follow-check endpoints answer 204 when following and 404 when not, so
// the initial state comes from the server rather than assuming "not following"
// and showing Follow to someone who already does.
export function FollowButton({ kind, id }: FollowButtonProps) {
  const api = useApi();
  const [isFollowing, setIsFollowing] = useState<boolean | null>(null);
  const [isBusy, setIsBusy] = useState(false);

  const check = useCallback(async () => {
    const res =
      kind === "user"
        ? await api.User.isFollowingUser({ params: { id } }).catch(() => null)
        : await api.Community.isFollowingCommunity({ params: { id } }).catch(() => null);
    setIsFollowing(res?.status === 204);
  }, [api, id, kind]);

  useEffect(() => {
    check();
  }, [check]);

  async function toggle() {
    if (isBusy || isFollowing === null) return;
    setIsBusy(true);

    const action = isFollowing
      ? kind === "user"
        ? api.User.unfollowUser({ params: { id } })
        : api.Community.unfollowCommunity({ params: { id } })
      : kind === "user"
        ? api.User.followUser({ params: { id } })
        : api.Community.followCommunity({ params: { id } });

    const res = await action.catch(() => null);
    setIsBusy(false);

    if (res && (res.status === 200 || res.status === 201 || res.status === 204)) {
      setIsFollowing(!isFollowing);
    } else {
      // Re-check rather than guess -- a failed follow and an already-following
      // conflict look the same from here.
      check();
    }
  }

  const label = isFollowing === null ? "…" : isFollowing ? "Following" : "Follow";

  return (
    <button
      onClick={toggle}
      disabled={isBusy || isFollowing === null}
      className={`ml-auto shrink-0 border px-2.5 py-1 font-mono text-[0.64rem] font-bold tracking-[0.05em] uppercase disabled:opacity-60 ${
        isFollowing
          ? "border-accent/50 text-accent-strong hover:bg-accent/10"
          : "border-border text-text-muted hover:border-accent hover:text-accent-strong"
      }`}
    >
      {label}
    </button>
  );
}
