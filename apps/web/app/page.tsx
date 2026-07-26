import Link from "next/link";
import { Logo } from "@/components/logo";
import { ConstellationCanvas } from "@/components/constellation-canvas";

const COMMUNITIES = ["n/indiehackers", "n/photography", "n/climbing", "n/lofi", "n/gamedev"];

const FEATURES = [
  {
    num: "01",
    title: "Communities, not feeds",
    body: "Every post lives in a community you chose to join — n/gamedev doesn't leak into n/climbing.",
  },
  {
    num: "02",
    title: "Follow people, not just topics",
    body: "Your home feed blends trending posts with ones from the people and communities you actually follow.",
  },
  {
    num: "03",
    title: "Moderation that scales with you",
    body: "Every community gets its own mods, roles, and report queue from the moment it's created.",
  },
];

export default function JoinPage() {
  return (
    <main>
      <nav className="flex items-center justify-between px-[6vw] py-5">
        <Logo />
        <div className="flex gap-2.5">
          <Link
            href="/sign-in"
            className="rounded-[9px] border border-border px-5 py-2.5 text-sm font-bold text-text hover:border-accent/35 hover:bg-accent/10"
          >
            Sign in
          </Link>
          <Link
            href="/sign-up"
            className="rounded-[9px] bg-accent px-5 py-2.5 text-sm font-bold text-accent-text hover:bg-accent-strong"
          >
            Join Nexus
          </Link>
        </div>
      </nav>

      <section className="relative grid place-items-center overflow-hidden px-[6vw] py-[9vh] text-center">
        <ConstellationCanvas />
        <div className="relative flex max-w-[680px] flex-col items-center gap-5">
          <span className="font-mono text-[0.72rem] tracking-[0.14em] text-accent-strong uppercase">
            now in open beta
          </span>
          <h1 className="text-[clamp(2.4rem,5.4vw,4.1rem)] leading-[1.04]">
            Every interest has
            <br />
            <em className="text-accent not-italic">a room here.</em>
          </h1>
          <p className="max-w-[46ch] text-[1.08rem] leading-relaxed text-text-muted">
            Nexus is where communities form around what you&apos;re actually into —
            post, vote, follow the people and spaces that get it, and watch a feed
            that&apos;s shaped by them instead of a stranger&apos;s algorithm.
          </p>
          <div className="mt-1.5 flex flex-wrap justify-center gap-3">
            <Link
              href="/sign-up"
              className="rounded-[9px] bg-accent px-6.5 py-3.5 text-[0.98rem] font-bold text-accent-text hover:bg-accent-strong"
            >
              Create your account
            </Link>
            <Link
              href="/dashboard"
              className="rounded-[9px] border border-border px-6.5 py-3.5 text-[0.98rem] font-bold text-text hover:border-accent/35 hover:bg-accent/10"
            >
              Browse communities
            </Link>
          </div>
          <div className="mt-2 flex flex-wrap justify-center gap-2.5">
            {COMMUNITIES.map((c) => (
              <span
                key={c}
                className="rounded-full border border-border bg-surface px-3 py-1.5 font-mono text-[0.78rem] text-text-muted"
              >
                {c}
              </span>
            ))}
          </div>
        </div>
      </section>

      <section className="grid grid-cols-1 gap-px border-y border-border-soft bg-border-soft sm:grid-cols-3">
        {FEATURES.map((f) => (
          <div key={f.num} className="flex flex-col gap-2.5 bg-bg px-[6vw] py-10">
            <span className="font-mono text-[0.78rem] text-accent-strong">{f.num}</span>
            <h3 className="text-[1.15rem] font-bold">{f.title}</h3>
            <p className="text-[0.92rem] leading-relaxed text-text-muted">{f.body}</p>
          </div>
        ))}
      </section>

      <section className="grid place-items-center px-[6vw] py-[8vh]">
        <div className="w-full max-w-[560px] overflow-hidden rounded-[20px] border border-border bg-surface shadow-[0_24px_60px_-20px_rgba(0,0,0,0.6)]">
          <div className="flex items-center justify-between border-b border-border-soft px-[18px] py-3.5">
            <Logo className="text-[0.95rem]" />
            <span className="font-mono text-[0.72rem] tracking-[0.14em] text-text-faint uppercase">
              live preview
            </span>
          </div>
          <div className="flex gap-3.5 p-4">
            <div className="flex min-w-[30px] flex-col items-center gap-1">
              <span className="text-up">▲</span>
              <span className="font-mono text-[0.82rem] font-semibold">812</span>
              <span className="text-text-faint">▼</span>
            </div>
            <div className="flex min-w-0 flex-col gap-1.5">
              <div className="flex flex-wrap items-center gap-1.5 text-[0.76rem] text-text-faint">
                <span className="rounded-full bg-accent/10 px-2 py-0.5 font-mono text-[0.72rem] text-accent-strong">
                  n/indiehackers
                </span>
                <span>posted by @mira.codes</span>
                <span>·</span>
                <span>4h ago</span>
              </div>
              <h3 className="text-[1rem] font-bold">
                Shipped my first paid feature after 8 months of building in public
              </h3>
              <p className="text-[0.86rem] leading-relaxed text-text-muted">
                Small win but it&apos;s real revenue — writing up the whole
                pricing experiment in the comments if anyone wants the numbers.
              </p>
              <div className="mt-0.5 flex gap-4 text-[0.78rem] text-text-faint">
                <span>💬 96 comments</span>
                <span>↗ share</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}
