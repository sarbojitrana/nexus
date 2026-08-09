"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { useAuth } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { MediaGallery } from "@/components/media/media-gallery";
import { EditProfileModal } from "@/components/profile/edit-profile-modal";
import { RemoteAvatar } from "@/components/media/remote-image";
import { ZoomableAvatar, ZoomableBanner } from "@/components/media/zoomable";
import { ShareButton } from "@/components/share-button";
import { formatTimeAgo } from "@/lib/format";
import type { User, MiniUser, PopulatedPost } from "@nexus/zod";

const TABS = ["Posts", "Followers", "Following"] as const;
type Tab = (typeof TABS)[number];

export function ProfileView({ userId }: { userId: string }) {
  const api = useApi();
  const { userId: viewerId } = useAuth();

  const [user, setUser] = useState<User | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [tab, setTab] = useState<Tab>("Posts");
  const [isFollowing, setIsFollowing] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [isBlocked, setIsBlocked] = useState(false);

  const load = useCallback(async () => {
    const isSelf = viewerId === userId;
    const res = isSelf
      ? await api.User.getMe().catch(() => null)
      : await api.User.getUserById({ params: { id: userId } }).catch(() => null);

    if (!res || res.status !== 200) {
      setNotFound(true);
      return;
    }
    setUser(res.body);

    if (viewerId && viewerId !== userId) {
      const [followRes, blockRes] = await Promise.all([
        api.User.isFollowingUser({ params: { id: userId } }).catch(() => null),
        api.User.isBlockingUser({ params: { id: userId } }).catch(() => null),
      ]);
      setIsFollowing(followRes?.status === 204);
      setIsBlocked(blockRes?.status === 204);
    }
  }, [api, userId, viewerId]);

  useEffect(() => {
    load();
  }, [load]);

  async function toggleBlock() {
    const action = isBlocked
      ? api.User.unblockUser({ params: { id: userId } })
      : api.User.blockUser({ params: { id: userId } });
    const res = await action.catch(() => null);
    if (res && res.status === 204) {
      setIsBlocked(!isBlocked);
      // Blocking drops the follow server-side, so refetch rather than guess.
      load();
    }
  }

  async function toggleFollow() {
    const action = isFollowing
      ? api.User.unfollowUser({ params: { id: userId } })
      : api.User.followUser({ params: { id: userId } });
    const res = await action.catch(() => null);
    if (res) {
      setIsFollowing(!isFollowing);
      load();
    }
  }

  if (notFound) {
    return (
      <div className="px-6 py-8">
        <div className="border border-border bg-surface p-8 text-center">
          <p className="text-[0.88rem] text-text-muted">
            This profile doesn&apos;t exist, or its owner keeps it private.
          </p>
        </div>
      </div>
    );
  }

  if (!user) {
    return <p className="px-6 py-8 font-mono text-[0.78rem] text-text-faint">loading…</p>;
  }

  const isSelf = viewerId === user.id;

  return (
    <div className="mx-auto flex min-h-0 w-full max-w-[820px] flex-col gap-4 overflow-y-auto px-6 py-6">
      <header className="border border-border bg-surface">
        <ZoomableBanner storageKey={user.bannerKey} className="h-32 border-b border-border" />
        <div className="flex flex-wrap items-start gap-4 p-5">
          <ZoomableAvatar storageKey={user.avatarKey} url={user.avatarUrl} size={72} className="-mt-12 border-2 border-surface" />
          <div className="min-w-0 flex-1">
            <h1 className="font-display text-[1.3rem] font-extrabold">@{user.username}</h1>
            <p className="mt-0.5 text-[0.88rem] text-text-muted">{user.displayName}</p>
            {user.bio && (
              <p className="mt-2 text-[0.86rem] leading-relaxed text-text-muted">{user.bio}</p>
            )}
            <div className="mt-3 flex flex-wrap gap-4 font-mono text-[0.7rem] text-text-faint">
              <span>{user.followerCount} followers</span>
              <span>{user.followingCount} following</span>
              <span>{user.postsCount} posts</span>
              <span>joined {formatTimeAgo(user.createdAt)}</span>
              <ShareButton
                path={`/dashboard/profile/${user.id}`}
                title={`@${user.username} on Nexus`}
              />
            </div>
          </div>

          {!isSelf && (
            <div className="flex shrink-0 flex-wrap gap-2">
              <button
                onClick={toggleBlock}
                className={`px-4 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] uppercase ${
                  isBlocked
                    ? "border border-accent text-accent-strong hover:bg-accent/10"
                    : "border border-border text-text-faint hover:border-accent/40 hover:text-accent-strong"
                }`}
              >
                {isBlocked ? "Unblock" : "Block"}
              </button>
              <button
                onClick={toggleFollow}
                disabled={isBlocked}
                className={`px-5 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] uppercase ${
                  isFollowing
                    ? "border border-border text-text-muted hover:border-accent/40"
                    : "bg-accent text-accent-text hover:bg-accent-strong"
                }`}
              >
                {isFollowing ? "Following" : "Follow"}
              </button>
              {!isBlocked && <StartChatButton userId={user.id} />}
            </div>
          )}
          {isSelf && (
            <button
              onClick={() => setShowEdit(true)}
              className="shrink-0 border border-border px-5 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] text-text-muted uppercase hover:border-accent/40 hover:text-accent-strong"
            >
              Update profile
            </button>
          )}
        </div>
      </header>

      <nav className="flex gap-5 border-b border-border-soft">
        {TABS.map((t) => (
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

      {tab === "Posts" && <ProfilePosts userId={user.id} username={user.username} />}
      {tab === "Followers" && <FollowList userId={user.id} kind="followers" />}
      {tab === "Following" && <FollowList userId={user.id} kind="following" />}

      {showEdit && (
        <EditProfileModal
          user={user}
          onClose={() => setShowEdit(false)}
          onSaved={setUser}
        />
      )}
    </div>
  );
}

function StartChatButton({ userId }: { userId: string }) {
  const api = useApi();
  const [isStarting, setIsStarting] = useState(false);

  async function start() {
    setIsStarting(true);
    const res = await api.Chat.startDirectConversation({ body: { userId } }).catch(() => null);
    setIsStarting(false);
    if (res && res.status === 200) {
      window.location.href = `/dashboard/messages?conversation=${res.body.id}`;
    }
  }

  return (
    <button
      onClick={start}
      disabled={isStarting}
      className="border border-border px-5 py-2 font-mono text-[0.7rem] font-bold tracking-[0.06em] text-text-muted uppercase hover:border-accent/40 disabled:opacity-50"
    >
      {isStarting ? "…" : "Message"}
    </button>
  );
}

function ProfilePosts({ userId, username }: { userId: string; username: string }) {
  const api = useApi();
  const [posts, setPosts] = useState<PopulatedPost[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    api.User
      .getUserPosts({ params: { id: userId }, query: {} })
      .then((res) => {
        if (res.status === 200) setPosts(res.body.data);
        setIsLoading(false);
      })
      .catch(() => setIsLoading(false));
  }, [api, userId]);

  if (isLoading) return <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>;

  if (posts.length === 0) {
    return (
      <div className="border border-border bg-surface p-8 text-center">
        <p className="text-[0.88rem] text-text-muted">Nothing to show.</p>
        <p className="mt-1 font-mono text-[0.72rem] text-text-faint">
          @{username} hasn&apos;t posted yet.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-3">
      {posts.map((p) => (
        <div key={p.id} className="border border-border bg-surface p-4">
          <Link href={`/dashboard/posts/${p.id}`}>
            <h3 className="font-display text-[0.98rem] font-bold hover:text-accent-strong">
              {p.title ?? p.content?.slice(0, 80)}
            </h3>
          </Link>
          {p.title && p.content && (
            <p className="mt-1 line-clamp-2 text-[0.85rem] text-text-muted">{p.content}</p>
          )}
          <MediaGallery media={p.postMedia ?? []} />
          <div className="mt-2 flex gap-4 font-mono text-[0.68rem] text-text-faint">
            <span>{p.upvotes - p.downvotes} points</span>
            <span>{p.commentCount} comments</span>
            <span>{formatTimeAgo(p.createdAt)}</span>
          </div>
        </div>
      ))}
    </div>
  );
}

function FollowList({ userId, kind }: { userId: string; kind: "followers" | "following" }) {
  const api = useApi();
  const [users, setUsers] = useState<MiniUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const request =
      kind === "followers"
        ? api.User.getFollowers({ params: { id: userId }, query: {} })
        : api.User.getFollowing({ params: { id: userId }, query: {} });

    request
      .then((res) => {
        if (res.status === 200) setUsers(res.body.data);
        setIsLoading(false);
      })
      .catch(() => setIsLoading(false));
  }, [api, userId, kind]);

  if (isLoading) return <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>;

  if (users.length === 0) {
    return (
      <div className="border border-border bg-surface p-8 text-center text-[0.88rem] text-text-muted">
        Nothing to show.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      {users.map((u) => (
        <Link
          key={u.id}
          href={`/dashboard/profile/${u.id}`}
          className="flex items-center gap-3 border border-border bg-surface p-3.5 hover:border-accent/40"
        >
          <RemoteAvatar storageKey={u.avatarKey} url={u.avatarUrl} size={32} />
          <div className="min-w-0">
            <strong className="block truncate text-[0.86rem] font-bold">@{u.username}</strong>
            <span className="font-mono text-[0.66rem] text-text-faint">
              {u.followerCount} followers
            </span>
          </div>
        </Link>
      ))}
    </div>
  );
}
