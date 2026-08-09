import Link from "next/link";
import { UserButton } from "@clerk/nextjs";
import { Logo } from "@/components/logo";
import { NavUnread } from "@/components/nav-unread";
import { StatusReadout } from "@/components/status-readout";
import { getServerApi } from "@/lib/api-server";

const NAV_LINKS = [
  { label: "Home feed", href: "/dashboard" },
  { label: "Communities", href: "/dashboard/communities" },
  { label: "Following", href: "/dashboard/following" },
  { label: "Messages", href: "/dashboard/messages", unread: "messages" as const },
  { label: "Notifications", href: "/dashboard/notifications", unread: "notifications" as const },
  { label: "Profile", href: "/dashboard/profile" },
];

export async function DashboardSidebar({ active }: { active: string }) {
  const api = await getServerApi();

  // getMe doubles as account bootstrap: it provisions the local users row from
  // Clerk if the user.created webhook never landed. Every write in the app has
  // a foreign key to users, so this has to succeed before anything else works.
  const [meRes, res] = await Promise.all([
    api.User.getMe({ fetchOptions: { cache: "no-store" } }).catch(() => null),
    api.Community.getCommunities({
      query: { sort: "members_count", order: "desc" },
      fetchOptions: { cache: "no-store" },
    }).catch(() => null),
  ]);

  const me = meRes && meRes.status === 200 ? meRes.body : null;
  const communities = res && res.status === 200 ? res.body.data.slice(0, 6) : [];

  return (
    <aside className="hidden min-h-0 flex-col gap-7 overflow-y-auto border-r border-border-soft px-5 py-6 lg:flex">
      <div className="flex items-start justify-between">
        <Logo />
        <UserButton />
      </div>

      {me && (
        <div className="-mt-4 font-mono text-[0.7rem] text-text-faint">
          @{me.username}
        </div>
      )}

      <nav className="flex flex-col">
        {NAV_LINKS.map((l) => {
          const isActive = l.href === active;
          return (
            <Link
              key={l.label}
              href={l.href}
              className={`flex items-center gap-2.5 border-b border-border-soft py-2.5 font-mono text-[0.8rem] ${
                isActive
                  ? "text-accent-strong"
                  : "text-text-muted hover:text-text"
              }`}
            >
              <span className={isActive ? "text-accent" : "opacity-0"}>›</span>
              {l.label}
              {l.unread && <NavUnread kind={l.unread} />}
            </Link>
          );
        })}
      </nav>

      <div>
        <div className="eyebrow mb-2">Communities</div>
        <div className="flex flex-col gap-1.5">
          {communities.map((c) => (
            <Link
              key={c.communityId}
              href={`/dashboard/communities/${c.slug}`}
              className="flex items-center justify-between font-mono text-[0.76rem] text-text-muted hover:text-accent-strong"
            >
              <span className="truncate">n/{c.slug}</span>
              <span className="shrink-0 pl-2 text-text-faint">{c.membersCount}</span>
            </Link>
          ))}
          {communities.length === 0 && (
            <span className="font-mono text-[0.74rem] text-text-faint">Nothing to show</span>
          )}
        </div>
      </div>

      <div className="mt-auto">
        <StatusReadout className="!text-left" />
      </div>
    </aside>
  );
}
