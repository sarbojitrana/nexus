import { UserButton } from "@clerk/nextjs";
import { Logo } from "@/components/logo";
import { PostCard } from "@/components/post-card";

const MY_COMMUNITIES = [
  { name: "n/indiehackers", members: "12.4k" },
  { name: "n/photography", members: "8.1k" },
  { name: "n/climbing", members: "3.9k" },
  { name: "n/lofi", members: "21k" },
];

const TRENDING = [
  { name: "n/lofi", members: "21,048 members" },
  { name: "n/gamedev", members: "9,762 members" },
  { name: "n/plants", members: "6,310 members" },
];

const PEOPLE = [
  { handle: "@theo_writes", meta: "Photography · 1.2k followers" },
  { handle: "@priya.builds", meta: "Indiehackers · 890 followers" },
];

const POSTS = [
  {
    community: "n/photography",
    author: "@theo_writes",
    timeAgo: "1h ago",
    title: "Shot the whole set on a $40 thrifted lens — full breakdown inside",
    snippet:
      "Manual focus only, no stabilization. Happy to answer anything about the setup.",
    votes: 1200,
    comments: 214,
    voted: true,
  },
  {
    community: "n/climbing",
    author: "@nadia.k",
    timeAgo: "3h ago",
    title: "Finally sent my first V6 after 6 months of failing on the same move",
    snippet: "The crimp was never the problem — it was foot placement the entire time.",
    votes: 433,
    comments: 58,
  },
  {
    community: "n/indiehackers",
    author: "@priya.builds",
    timeAgo: "5h ago",
    title: "Ask Nexus: how do you price a B2B tool when you have zero customers yet?",
    snippet: `Everyone says "charge more than feels comfortable" but nobody says a number.`,
    votes: 89,
    comments: 141,
  },
];

const NAV_LINKS = [
  { label: "Home feed", active: true },
  { label: "Communities" },
  { label: "Following" },
  { label: "Notifications" },
  { label: "Profile" },
];

export default function DashboardPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr_268px]">
      <aside className="hidden flex-col gap-6.5 border-r border-border-soft px-4 py-5.5 lg:flex">
        <div className="flex items-center justify-between">
          <Logo />
          <UserButton />
        </div>
        <nav className="flex flex-col gap-0.5">
          {NAV_LINKS.map((l) => (
            <a
              key={l.label}
              href="#"
              className={`flex items-center gap-2.5 rounded-[9px] px-3 py-2.5 text-[0.88rem] font-semibold ${
                l.active
                  ? "bg-accent/10 text-accent-strong"
                  : "text-text-muted hover:bg-surface hover:text-text"
              }`}
            >
              <span
                className={`h-1.5 w-1.5 rounded-full ${l.active ? "bg-accent" : "bg-current opacity-40"}`}
              />
              {l.label}
            </a>
          ))}
        </nav>
        <div>
          <div className="mb-1 px-3 font-mono text-[0.68rem] tracking-[0.1em] text-text-faint uppercase">
            Your communities
          </div>
          <div className="flex flex-col">
            {MY_COMMUNITIES.map((c) => (
              <a
                key={c.name}
                href="#"
                className="flex items-center justify-between rounded-[9px] px-3 py-1.5 text-[0.83rem] text-text-muted hover:bg-surface hover:text-text"
              >
                {c.name}
                <span className="font-mono text-[0.72rem] text-text-faint">{c.members}</span>
              </a>
            ))}
          </div>
        </div>
      </aside>

      <main className="flex max-w-[640px] flex-col gap-4 px-6 py-5.5">
        <div className="flex items-center justify-between">
          <div className="flex gap-1">
            <span className="rounded-full bg-accent/10 px-3 py-1.5 text-[0.8rem] font-bold text-accent-strong">
              For you
            </span>
            <span className="rounded-full px-3 py-1.5 text-[0.8rem] font-bold text-text-muted">
              Trending
            </span>
            <span className="rounded-full px-3 py-1.5 text-[0.8rem] font-bold text-text-muted">
              Following
            </span>
          </div>
          <a
            href="#"
            className="rounded-[9px] bg-accent px-4 py-2 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong"
          >
            + New post
          </a>
        </div>

        {POSTS.map((p) => (
          <PostCard key={p.title} {...p} />
        ))}
      </main>

      <aside className="hidden flex-col gap-5.5 border-l border-border-soft px-4.5 py-5.5 lg:flex">
        <div className="rounded-2xl border border-border bg-surface p-4">
          <h4 className="mb-3 text-[0.82rem] font-bold">Trending communities</h4>
          <div className="flex flex-col gap-3">
            {TRENDING.map((t) => (
              <div key={t.name} className="flex items-center gap-2.5">
                <div className="h-[30px] w-[30px] shrink-0 rounded-[9px] bg-gradient-to-br from-accent to-down" />
                <div className="flex min-w-0 flex-col gap-px">
                  <strong className="text-[0.82rem] font-bold">{t.name}</strong>
                  <span className="text-[0.72rem] text-text-faint">{t.members}</span>
                </div>
                <button className="ml-auto shrink-0 rounded-full border border-border px-2.5 py-1 text-[0.72rem] font-bold text-text-muted">
                  Follow
                </button>
              </div>
            ))}
          </div>
        </div>
        <div className="rounded-2xl border border-border bg-surface p-4">
          <h4 className="mb-3 text-[0.82rem] font-bold">People to follow</h4>
          <div className="flex flex-col gap-3">
            {PEOPLE.map((p) => (
              <div key={p.handle} className="flex items-center gap-2.5">
                <div className="h-[30px] w-[30px] shrink-0 rounded-full bg-gradient-to-br from-accent to-down" />
                <div className="flex min-w-0 flex-col gap-px">
                  <strong className="text-[0.82rem] font-bold">{p.handle}</strong>
                  <span className="text-[0.72rem] text-text-faint">{p.meta}</span>
                </div>
                <button className="ml-auto shrink-0 rounded-full border border-border px-2.5 py-1 text-[0.72rem] font-bold text-text-muted">
                  Follow
                </button>
              </div>
            ))}
          </div>
        </div>
      </aside>
    </div>
  );
}
