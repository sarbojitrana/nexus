import Link from "next/link";
import { Show } from "@clerk/nextjs";
import { Logo, HudCorners } from "@/components/logo";
import { StatusReadout } from "@/components/status-readout";
import { ConstellationCanvas } from "@/components/constellation-canvas";

const FEATURES = [
  {
    num: "01",
    label: "communities",
    title: "Not feeds. Rooms.",
    body: "Every post lives inside a community you chose — n/gamedev never leaks into n/climbing.",
  },
  {
    num: "02",
    label: "signal",
    title: "Follow people, not platforms",
    body: "Your feed blends trending posts with the people and communities you actually follow.",
  },
  {
    num: "03",
    label: "moderation",
    title: "Every room, self-governed",
    body: "Mods, roles, and a report queue exist from the moment a community is created.",
  },
];

export default function JoinPage() {
  return (
    <main className="min-h-dvh">
      <HudCorners />

      <header className="flex flex-wrap items-start justify-between gap-4 px-8 py-5">
        <div className="flex items-start gap-8">
          <Logo />
          <div className="hidden flex-col gap-1 pt-0.5 sm:flex">
            <span className="eyebrow">community protocol // est. 2026</span>
            <span className="eyebrow">
              field status: <span className="text-text-muted">accepting connections</span>
            </span>
          </div>
        </div>

        <div className="flex items-start gap-6">
          <StatusReadout />
          <Show when="signed-out">
            <Link
              href="/sign-up"
              className="border border-accent px-5 py-2.5 font-mono text-[0.7rem] font-bold tracking-[0.1em] text-accent-strong uppercase hover:bg-accent/10"
            >
              Join Nexus
            </Link>
          </Show>
          <Show when="signed-in">
            <Link
              href="/dashboard"
              className="bg-accent px-5 py-2.5 font-mono text-[0.7rem] font-bold tracking-[0.1em] text-accent-text uppercase hover:bg-accent-strong"
            >
              Open dashboard
            </Link>
          </Show>
        </div>
      </header>

      <section className="relative grid grid-cols-1 gap-10 border-t border-border-soft px-8 py-16 lg:grid-cols-[1fr_320px]">
        <ConstellationCanvas />

        <div className="relative">
          <span className="eyebrow text-accent-strong">now in open beta</span>
          <h1 className="mt-3 text-[clamp(2.2rem,4.4vw,3.4rem)] leading-[1.06] font-extrabold">
            Every interest
            <br />
            finds <em className="text-accent not-italic">its node.</em>
          </h1>

          <p className="mt-6 max-w-[52ch] text-[1rem] leading-relaxed text-text-muted">
            Nexus routes you to the communities and people worth your time — post, vote,
            follow — then gets out of the way.
          </p>

          <div className="mt-8 flex flex-wrap gap-3">
            <Show when="signed-out">
              <Link
                href="/sign-up"
                className="border border-dashed border-accent/60 px-6 py-3 font-mono text-[0.74rem] font-bold tracking-[0.08em] text-accent-strong uppercase hover:bg-accent/10"
              >
                Create account ↗
              </Link>
              <Link
                href="/sign-in"
                className="border border-border px-6 py-3 font-mono text-[0.74rem] font-bold tracking-[0.08em] text-text uppercase hover:border-accent/40 hover:bg-accent/5"
              >
                Sign in
              </Link>
            </Show>
            <Show when="signed-in">
              <Link
                href="/dashboard"
                className="border border-dashed border-accent/60 px-6 py-3 font-mono text-[0.74rem] font-bold tracking-[0.08em] text-accent-strong uppercase hover:bg-accent/10"
              >
                Go to your feed ↗
              </Link>
              <Link
                href="/dashboard/communities"
                className="border border-border px-6 py-3 font-mono text-[0.74rem] font-bold tracking-[0.08em] text-text uppercase hover:border-accent/40 hover:bg-accent/5"
              >
                Browse communities
              </Link>
            </Show>
          </div>
        </div>

        <div className="relative">
          <div className="border border-border p-5">
            <div className="mb-3 border-b border-dashed border-border pb-3 font-mono text-[0.64rem] tracking-[0.1em] text-text-faint">
              — nexus // community os —
            </div>
            <ReadoutRow label="protocol" value="open" />
            <ReadoutRow label="communities" value="self-governed" />
            <ReadoutRow label="algorithm" value="yours" />
          </div>
        </div>
      </section>

      <section className="grid grid-cols-1 border-t border-border-soft sm:grid-cols-3">
        {FEATURES.map((f) => (
          <div key={f.num} className="border-r border-border-soft px-7 py-9 last:border-r-0">
            <span className="font-mono text-[0.7rem] tracking-[0.06em] text-accent-strong">
              {f.num} / {f.label}
            </span>
            <h3 className="mt-2 text-[1.08rem] font-bold">{f.title}</h3>
            <p className="mt-1.5 text-[0.88rem] leading-relaxed text-text-muted">{f.body}</p>
          </div>
        ))}
      </section>

      <footer className="flex flex-wrap items-center justify-between gap-3 border-t border-border-soft px-8 py-5 font-mono text-[0.66rem] tracking-[0.04em] text-text-faint">
        <span>Nexus Network — all systems own</span>
        <span>built on communities, not algorithms</span>
      </footer>
    </main>
  );
}

function ReadoutRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="mb-2 flex justify-between font-mono text-[0.76rem] last:mb-0">
      <span className="tracking-[0.04em] text-text-faint">{label}</span>
      <span className="font-bold text-accent-strong">{value}</span>
    </div>
  );
}
