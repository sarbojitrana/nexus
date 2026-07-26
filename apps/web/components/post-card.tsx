type PostCardProps = {
  community: string;
  author: string;
  timeAgo: string;
  title: string;
  snippet: string;
  votes: number;
  comments: number;
  voted?: boolean;
};

export function PostCard({
  community,
  author,
  timeAgo,
  title,
  snippet,
  votes,
  comments,
  voted,
}: PostCardProps) {
  return (
    <article className="flex gap-3.5 rounded-2xl border border-border bg-surface p-4">
      <div className="flex min-w-[30px] flex-col items-center gap-1">
        <button
          aria-label="Upvote"
          className={`p-0.5 text-base leading-none ${voted ? "text-up" : "text-text-faint"}`}
        >
          ▲
        </button>
        <span className="font-mono text-[0.82rem] font-semibold tabular-nums">
          {votes >= 1000 ? `${(votes / 1000).toFixed(1)}k` : votes}
        </span>
        <button aria-label="Downvote" className="p-0.5 text-base leading-none text-text-faint">
          ▼
        </button>
      </div>
      <div className="flex min-w-0 flex-col gap-1.5">
        <div className="flex flex-wrap items-center gap-1.5 text-[0.76rem] text-text-faint">
          <span className="rounded-full bg-accent/10 px-2 py-0.5 font-mono text-[0.72rem] text-accent-strong">
            {community}
          </span>
          <span>{author}</span>
          <span>·</span>
          <span>{timeAgo}</span>
        </div>
        <h3 className="font-body text-[1rem] leading-snug font-bold">{title}</h3>
        <p className="text-[0.86rem] leading-relaxed text-text-muted">{snippet}</p>
        <div className="mt-0.5 flex gap-4 text-[0.78rem] text-text-faint">
          <span>💬 {comments} comments</span>
          <span>↗ share</span>
          <span>⚑ report</span>
        </div>
      </div>
    </article>
  );
}
