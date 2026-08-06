import Link from "next/link";
import { UserButton } from "@clerk/nextjs";
import { Logo } from "@/components/logo";
import { getServerApi } from "@/lib/api-server";

const NAV_LINKS = [
  { label: "Home feed", href: "/dashboard" },
  { label: "Communities", href: "/dashboard/communities" },
  { label: "Following", href: "/dashboard/following" },
  { label: "Messages", href: "/dashboard/messages" },
  { label: "Notifications", href: "/dashboard/notifications" },
];

export async function DashboardSidebar({ active }: { active: string }) {
  const api = await getServerApi();
  const res = await api.Community.getCommunities({
    query: { sort: "created_at", order: "desc" },
    fetchOptions: { cache: "no-store" },
  });
  const communities = res.status === 200 ? res.body.data.slice(0, 5) : [];

  return (
    <aside className="hidden flex-col gap-6.5 border-r border-border-soft px-4 py-5.5 lg:flex">
      <div className="flex items-center justify-between">
        <Logo />
        <UserButton />
      </div>
      <nav className="flex flex-col gap-0.5">
        {NAV_LINKS.map((l) => {
          const isActive = l.href === active;
          return (
            <Link
              key={l.label}
              href={l.href}
              className={`flex items-center gap-2.5 rounded-[9px] px-3 py-2.5 text-[0.88rem] font-semibold ${
                isActive
                  ? "bg-accent/10 text-accent-strong"
                  : "text-text-muted hover:bg-surface hover:text-text"
              }`}
            >
              <span
                className={`h-1.5 w-1.5 rounded-full ${isActive ? "bg-accent" : "bg-current opacity-40"}`}
              />
              {l.label}
            </Link>
          );
        })}
      </nav>
      <div>
        <div className="mb-1 px-3 font-mono text-[0.68rem] tracking-[0.1em] text-text-faint uppercase">
          Communities
        </div>
        <div className="flex flex-col">
          {communities.map((c) => (
            <div
              key={c.communityId}
              className="flex items-center justify-between rounded-[9px] px-3 py-1.5 text-[0.83rem] text-text-muted hover:bg-surface hover:text-text"
            >
              n/{c.communityName}
              <span className="font-mono text-[0.72rem] text-text-faint">{c.membersCount}</span>
            </div>
          ))}
          {communities.length === 0 && (
            <span className="px-3 py-1.5 text-[0.8rem] text-text-faint">No communities yet</span>
          )}
        </div>
      </div>
    </aside>
  );
}
