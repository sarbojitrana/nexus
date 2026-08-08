"use client";

import Link from "next/link";
import { useCallback, useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { useApi } from "@/lib/use-api";
import { useUpload } from "@/lib/use-upload";
import { apiErrorMessage } from "@/lib/api-error";
import { RemoteAvatar, RemoteBanner } from "@/components/media/remote-image";
import { ZoomableAvatar, ZoomableBanner } from "@/components/media/zoomable";
import { NewPostButton } from "@/components/new-post-modal";
import { formatTimeAgo } from "@/lib/format";
import type {
  CommunityResponse,
  MiniCommunityUser,
  CommunityReport,
  PopulatedPost,
} from "@nexus/zod";

const TABS = ["Posts", "Members", "Moderation", "Settings"] as const;
type Tab = (typeof TABS)[number];

export function CommunityDetail({ slug }: { slug: string }) {
  const api = useApi();
  const router = useRouter();

  const [community, setCommunity] = useState<CommunityResponse | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [tab, setTab] = useState<Tab>("Posts");
  const [banner, setBanner] = useState<string | null>(null);

  const load = useCallback(async () => {
    const res = await api.Community.getCommunityByIdOrSlug({ params: { idOrSlug: slug } }).catch(
      () => null
    );
    if (!res || res.status !== 200) {
      setNotFound(true);
      return;
    }
    setCommunity(res.body);
  }, [api, slug]);

  useEffect(() => {
    load();
  }, [load]);

  if (notFound) {
    return (
      <div className="px-6 py-8">
        <div className="border border-border bg-surface p-8 text-center text-[0.88rem] text-text-muted">
          This community doesn&apos;t exist.
        </div>
      </div>
    );
  }

  if (!community) {
    return <p className="px-6 py-8 font-mono text-[0.78rem] text-text-faint">loading…</p>;
  }

  const isMember = community.viewerRole !== null;
  const isMod = community.viewerRole === "moderator" || community.viewerRole === "admin";
  const isAdmin = community.viewerRole === "admin";
  const visibleTabs = TABS.filter(
    (t) => (t !== "Moderation" || isMod) && (t !== "Settings" || isAdmin)
  );

  async function toggleMembership() {
    const action = isMember
      ? api.Community.leaveCommunity({ params: { id: community!.id } })
      : api.Community.joinCommunity({ params: { id: community!.id } });
    const res = await action.catch(() => null);
    if (res) load();
    else setBanner("Couldn't update membership.");
  }

  return (
    <div className="mx-auto flex w-full max-w-[880px] flex-col gap-4 px-6 py-6">
      {banner && (
        <div className="border border-accent/40 bg-accent/5 px-4 py-2.5 font-mono text-[0.74rem] text-accent-strong">
          {banner}
        </div>
      )}

      <header className="border border-border bg-surface">
        <ZoomableBanner storageKey={community.bannerKey} className="h-28 border-b border-border" />
        <div className="flex flex-wrap items-start gap-4 p-5">
          <ZoomableAvatar storageKey={community.avatarKey} size={64} className="-mt-11 border-2 border-surface" />
          <div className="min-w-0 flex-1">
            <h1 className="font-display text-[1.35rem] font-extrabold">n/{community.slug}</h1>
            <p className="mt-0.5 text-[0.86rem] text-text-muted">{community.name}</p>
            {community.description && (
              <p className="mt-2 text-[0.86rem] leading-relaxed text-text-muted">
                {community.description}
              </p>
            )}
            <div className="mt-3 flex flex-wrap gap-4 font-mono text-[0.7rem] text-text-faint">
              <span>{community.membersCount} members</span>
              <span>{community.postsCount} posts</span>
              <span>created {formatTimeAgo(community.createdAt)}</span>
              {community.viewerRole && (
                <span className="text-accent-strong">your role: {community.viewerRole}</span>
              )}
            </div>
          </div>
          <button
            onClick={toggleMembership}
            className={`shrink-0 px-5 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] uppercase ${
              isMember
                ? "border border-border text-text-muted hover:border-accent/40"
                : "bg-accent text-accent-text hover:bg-accent-strong"
            }`}
          >
            {isMember ? "Leave" : "Join"}
          </button>
        </div>
      </header>

      <nav className="flex gap-5 border-b border-border-soft">
        {visibleTabs.map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`pb-2.5 font-mono text-[0.72rem] font-bold tracking-[0.08em] uppercase ${
              tab === t
                ? "border-b-2 border-accent text-accent-strong"
                : "text-text-faint hover:text-text-muted"
            }`}
          >
            {t}
          </button>
        ))}
      </nav>

      {tab === "Posts" && (
        <CommunityPosts community={community} isMember={isMember} isMod={isMod} />
      )}
      {tab === "Members" && <CommunityMembers community={community} isAdmin={isAdmin} />}
      {tab === "Moderation" && isMod && <CommunityModeration community={community} />}
      {tab === "Settings" && isAdmin && (
        <CommunitySettings
          community={community}
          onSaved={load}
          onDeleted={() => router.push("/dashboard/communities")}
        />
      )}
    </div>
  );
}

function CommunityPosts({
  community,
  isMember,
  isMod,
}: {
  community: CommunityResponse;
  isMember: boolean;
  isMod: boolean;
}) {
  const api = useApi();
  const [posts, setPosts] = useState<PopulatedPost[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const load = useCallback(async () => {
    const res = await api.Community.getCommunityPosts({
      params: { id: community.id },
      query: {},
    }).catch(() => null);
    if (res && res.status === 200) setPosts(res.body.data);
    setIsLoading(false);
  }, [api, community.id]);

  useEffect(() => {
    load();
  }, [load]);

  async function modDelete(postId: string) {
    const res = await api.Community.deleteCommunityPost({
      params: { id: community.id, postId },
    }).catch(() => null);
    if (res && res.status === 204) load();
  }

  return (
    <div className="flex flex-col gap-3">
      {isMember && (
        <div className="flex justify-end">
          <NewPostButton communityId={community.id} />
        </div>
      )}

      {isLoading && <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>}

      {posts.map((p) => (
        <div key={p.id} className="border border-border bg-surface p-4">
          <Link href={`/dashboard/posts/${p.id}`}>
            <h3 className="font-display text-[0.98rem] font-bold hover:text-accent-strong">
              {p.title ?? p.content?.slice(0, 80)}
            </h3>
          </Link>
          <div className="mt-2 flex gap-4 font-mono text-[0.68rem] text-text-faint">
            <span>{p.upvotes - p.downvotes} points</span>
            <span>{p.commentCount} comments</span>
            <span>{formatTimeAgo(p.createdAt)}</span>
            {isMod && (
              <button onClick={() => modDelete(p.id)} className="hover:text-accent-strong">
                remove
              </button>
            )}
          </div>
        </div>
      ))}

      {!isLoading && posts.length === 0 && (
        <div className="border border-border bg-surface p-8 text-center">
          <p className="text-[0.88rem] text-text-muted">Nothing to show yet.</p>
          <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
            {isMember ? "Write the first post here." : "Join to start posting."}
          </p>
        </div>
      )}
    </div>
  );
}

function CommunityMembers({
  community,
  isAdmin,
}: {
  community: CommunityResponse;
  isAdmin: boolean;
}) {
  const api = useApi();
  const [members, setMembers] = useState<MiniCommunityUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const load = useCallback(async () => {
    const res = await api.Community.getCommunityMembers({
      params: { id: community.id },
      query: {},
    }).catch(() => null);
    if (res && res.status === 200) setMembers(res.body.data);
    setIsLoading(false);
  }, [api, community.id]);

  useEffect(() => {
    load();
  }, [load]);

  async function changeRole(userId: string, newRole: "member" | "moderator" | "admin") {
    const res = await api.Community.changeMemberRole({
      params: { id: community.id, userId },
      body: { newRole },
    }).catch(() => null);
    if (res) load();
  }

  async function ban(userId: string) {
    const res = await api.Community.banMember({
      params: { id: community.id },
      body: { userIdToBan: userId },
    }).catch(() => null);
    if (res) load();
  }

  return (
    <div className="flex flex-col gap-2">
      {isLoading && <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>}

      {members.map((m) => (
        <div
          key={m.userId}
          className="flex flex-wrap items-center gap-3 border border-border bg-surface p-3.5"
        >
          <RemoteAvatar storageKey={m.avatarKey} url={m.avatarUrl} size={32} />
          <Link href={`/dashboard/profile/${m.userId}`} className="min-w-0 flex-1">
            <strong className="block truncate text-[0.86rem] font-bold hover:text-accent-strong">
              {m.name}
            </strong>
            <span className="block truncate font-mono text-[0.7rem] text-text-muted">
              @{m.username}
            </span>
            <span className="font-mono text-[0.66rem] text-text-faint">
              {m.role} · joined {formatTimeAgo(m.joinedAt)}
            </span>
          </Link>

          {isAdmin && (
            <div className="flex shrink-0 items-center gap-2">
              <select
                value={m.role}
                onChange={(e) =>
                  changeRole(m.userId, e.target.value as "member" | "moderator" | "admin")
                }
                className="border border-border bg-bg px-2 py-1 font-mono text-[0.68rem] text-text focus:border-accent focus:outline-none"
              >
                <option value="member">member</option>
                <option value="moderator">moderator</option>
                <option value="admin">admin</option>
              </select>
              <button
                onClick={() => ban(m.userId)}
                className="border border-border px-2.5 py-1 font-mono text-[0.66rem] tracking-[0.05em] text-text-faint uppercase hover:border-accent hover:text-accent-strong"
              >
                Ban
              </button>
            </div>
          )}
        </div>
      ))}

      {!isLoading && members.length === 0 && (
        <div className="border border-border bg-surface p-8 text-center text-[0.86rem] text-text-muted">
          Nothing to show.
        </div>
      )}
    </div>
  );
}

function CommunityModeration({ community }: { community: CommunityResponse }) {
  const api = useApi();
  const [reports, setReports] = useState<CommunityReport[]>([]);
  const [statusFilter, setStatusFilter] = useState<"pending" | "resolved" | "dismissed">("pending");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const load = useCallback(async () => {
    setIsLoading(true);
    const res = await api.Community.getCommunityReports({
      params: { id: community.id },
      query: { status: statusFilter },
    }).catch(() => null);
    if (res && res.status === 200) setReports(res.body.data);
    setIsLoading(false);
  }, [api, community.id, statusFilter]);

  useEffect(() => {
    load();
  }, [load]);

  async function resolve(reportId: string, updatedStatus: "resolved" | "dismissed") {
    setError(null);
    const res = await api.Community.resolveReport({
      params: { id: community.id, reportId },
      body: { updatedStatus },
    }).catch(() => null);

    if (res && res.status === 200) {
      load();
      return;
    }
    setError(apiErrorMessage(res, "Couldn't update that report."));
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex gap-2">
        {(["pending", "resolved", "dismissed"] as const).map((s) => (
          <button
            key={s}
            onClick={() => setStatusFilter(s)}
            className={`border px-3 py-1.5 font-mono text-[0.68rem] tracking-[0.06em] uppercase ${
              statusFilter === s
                ? "border-accent text-accent-strong"
                : "border-border text-text-faint hover:text-text-muted"
            }`}
          >
            {s}
          </button>
        ))}
      </div>

      {error && (
        <p className="border border-accent/30 bg-accent/5 px-3 py-2 font-mono text-[0.72rem] text-accent-strong">
          {error}
        </p>
      )}

      {isLoading && <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>}

      {reports.map((r) => (
        <div key={r.id} className="border border-border bg-surface p-4">
          <div className="flex flex-wrap items-center gap-3 font-mono text-[0.68rem] text-text-faint">
            <span className="text-accent-strong">{r.status}</span>
            <span>reported {formatTimeAgo(r.createdAt)}</span>
            <Link href={`/dashboard/posts/${r.postId}`} className="hover:text-text-muted">
              view post ↗
            </Link>
          </div>
          <p className="mt-2 text-[0.86rem] leading-relaxed text-text">{r.reason}</p>

          {r.status === "pending" && (
            <div className="mt-3 flex gap-2">
              <button
                onClick={() => resolve(r.id, "resolved")}
                className="bg-up/15 px-3 py-1.5 font-mono text-[0.68rem] font-bold tracking-[0.05em] text-up uppercase"
              >
                Resolve
              </button>
              <button
                onClick={() => resolve(r.id, "dismissed")}
                className="border border-border px-3 py-1.5 font-mono text-[0.68rem] tracking-[0.05em] text-text-muted uppercase hover:border-accent/40"
              >
                Dismiss
              </button>
            </div>
          )}
        </div>
      ))}

      {!isLoading && reports.length === 0 && (
        <div className="border border-border bg-surface p-8 text-center">
          <p className="text-[0.88rem] text-text-muted">Nothing to show.</p>
          <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
            No {statusFilter} reports in this community.
          </p>
        </div>
      )}
    </div>
  );
}

function CommunitySettings({
  community,
  onSaved,
  onDeleted,
}: {
  community: CommunityResponse;
  onSaved: () => void;
  onDeleted: () => void;
}) {
  const api = useApi();
  const upload = useUpload();
  const avatarInput = useRef<HTMLInputElement>(null);
  const bannerInput = useRef<HTMLInputElement>(null);

  const [name, setName] = useState(community.name);
  const [description, setDescription] = useState(community.description ?? "");
  const [avatarKey, setAvatarKey] = useState(community.avatarKey);
  const [bannerKey, setBannerKey] = useState(community.bannerKey);
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);

  async function handleImage(file: File | undefined, target: "avatar" | "banner") {
    if (!file) return;
    setIsUploading(true);
    setMessage(null);
    const result = await upload(file);
    setIsUploading(false);
    if (!result.ok) {
      setMessage(result.error);
      return;
    }
    if (target === "avatar") setAvatarKey(result.media.storageKey);
    else setBannerKey(result.media.storageKey);
    setMessage("Image uploaded — save to apply.");
  }

  async function save() {
    setIsSaving(true);
    setMessage(null);
    const res = await api.Community.updateCommunitySettings({
      params: { id: community.id },
      body: {
        name: name.trim(),
        description: description.trim() || null,
        avatarKey,
        bannerKey,
      },
    }).catch(() => null);
    setIsSaving(false);
    if (res && res.status === 200) {
      setMessage("Saved.");
      onSaved();
    } else {
      setMessage(apiErrorMessage(res, "Couldn't save changes."));
    }
  }

  async function remove() {
    const res = await api.Community.deleteCommunity({ params: { id: community.id } }).catch(
      () => null
    );
    if (res && res.status === 204) onDeleted();
    else setMessage("Couldn't delete this community.");
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-5 border border-border bg-surface p-4">
        <div className="flex flex-col items-center gap-2">
          <RemoteAvatar storageKey={avatarKey} size={56} />
          <button
            onClick={() => avatarInput.current?.click()}
            disabled={isUploading}
            className="font-mono text-[0.64rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
          >
            Icon
          </button>
          <input
            ref={avatarInput}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(e) => handleImage(e.target.files?.[0], "avatar")}
          />
        </div>
        <div className="flex min-w-[160px] flex-1 flex-col items-center gap-2">
          <RemoteBanner storageKey={bannerKey} className="h-14" />
          <button
            onClick={() => bannerInput.current?.click()}
            disabled={isUploading}
            className="font-mono text-[0.64rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
          >
            Cover
          </button>
          <input
            ref={bannerInput}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(e) => handleImage(e.target.files?.[0], "banner")}
          />
        </div>
      </div>

      <label className="flex flex-col gap-1.5">
        <span className="eyebrow">Name</span>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
        />
      </label>

      <label className="flex flex-col gap-1.5">
        <span className="eyebrow">Description</span>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={4}
          className="resize-none border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
        />
      </label>

      <div className="flex items-center gap-3">
        <button
          onClick={save}
          disabled={isSaving}
          className="bg-accent px-5 py-2.5 font-mono text-[0.72rem] font-bold tracking-[0.06em] text-accent-text uppercase disabled:opacity-50"
        >
          {isSaving ? "Saving…" : "Save changes"}
        </button>
        {message && <span className="font-mono text-[0.72rem] text-text-muted">{message}</span>}
      </div>

      <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border border-accent/30 bg-accent/5 p-4">
        <div>
          <div className="text-[0.86rem] font-bold text-accent-strong">Delete community</div>
          <div className="mt-0.5 font-mono text-[0.7rem] text-text-faint">
            Permanently removes the community and its posts.
          </div>
        </div>
        {confirmDelete ? (
          <div className="flex gap-2">
            <button
              onClick={remove}
              className="bg-accent px-3 py-2 font-mono text-[0.68rem] font-bold tracking-[0.05em] text-accent-text uppercase"
            >
              Confirm
            </button>
            <button
              onClick={() => setConfirmDelete(false)}
              className="px-3 py-2 font-mono text-[0.68rem] tracking-[0.05em] text-text-muted uppercase"
            >
              Cancel
            </button>
          </div>
        ) : (
          <button
            onClick={() => setConfirmDelete(true)}
            className="border border-accent px-4 py-2 font-mono text-[0.68rem] font-bold tracking-[0.05em] text-accent-strong uppercase hover:bg-accent/10"
          >
            Delete
          </button>
        )}
      </div>
    </div>
  );
}
