"use client";

import { useState } from "react";
import { useApi } from "@/lib/use-api";
import { formatCount } from "@/lib/format";

// Comments are posts, so the same /posts/:id/react endpoint serves both --
// this is shared so a comment votes exactly like a post does.
//
// Vote state is optimistic: the endpoint returns 204 with no body and posts
// carry no "my vote" field, so the delta is tracked locally and rolled back
// if the request fails.
export function VoteButtons({
  postId,
  upvotes,
  downvotes,
  layout = "column",
}: {
  postId: string;
  upvotes: number;
  downvotes: number;
  layout?: "column" | "row";
}) {
  const api = useApi();
  const [vote, setVote] = useState<"upvote" | "downvote" | null>(null);
  const [delta, setDelta] = useState(0);

  async function react(reaction: "upvote" | "downvote") {
    const previous = vote;
    const next = previous === reaction ? null : reaction;

    const base = previous === "upvote" ? -1 : previous === "downvote" ? 1 : 0;
    const add = next === "upvote" ? 1 : next === "downvote" ? -1 : 0;
    const previousDelta = delta;

    setVote(next);
    setDelta(delta + base + add);

    const res = await api.Post.reactToPost({
      params: { id: postId },
      body: { reaction },
    }).catch(() => null);

    if (!res || (res.status !== 204 && res.status !== 200)) {
      setVote(previous);
      setDelta(previousDelta);
    }
  }

  const score = upvotes - downvotes + delta;
  const isRow = layout === "row";

  return (
    <div
      className={
        isRow
          ? "flex items-center gap-1.5"
          : "flex min-w-[34px] flex-col items-center gap-1"
      }
    >
      <button
        onClick={() => react("upvote")}
        aria-label="Upvote"
        aria-pressed={vote === "upvote"}
        className={`font-mono leading-none ${isRow ? "text-[0.8rem]" : "text-[0.9rem]"} ${
          vote === "upvote" ? "text-up" : "text-text-faint hover:text-up"
        }`}
      >
        ▲
      </button>

      <span
        className={`font-mono font-bold tabular-nums ${isRow ? "text-[0.7rem]" : "text-[0.78rem]"} ${
          vote === "upvote" ? "text-up" : vote === "downvote" ? "text-down" : ""
        }`}
      >
        {formatCount(score)}
      </span>

      <button
        onClick={() => react("downvote")}
        aria-label="Downvote"
        aria-pressed={vote === "downvote"}
        className={`font-mono leading-none ${isRow ? "text-[0.8rem]" : "text-[0.9rem]"} ${
          vote === "downvote" ? "text-down" : "text-text-faint hover:text-down"
        }`}
      >
        ▼
      </button>
    </div>
  );
}
