"use client";

import { useEffect, useRef, useState } from "react";
import { useClerk, useUser } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { useUpload } from "@/lib/use-upload";
import { RemoteAvatar, RemoteBanner } from "@/components/media/remote-image";
import type { User, MiniUser } from "@nexus/zod";

// Everything Nexus owns about a user lives here -- bio, username, banner, and
// the privacy settings. Clerk's own "Manage account" covers name, email,
// password and sessions, so none of that is duplicated.
export function EditProfileModal({
  user,
  onClose,
  onSaved,
}: {
  user: User;
  onClose: () => void;
  onSaved: (u: User) => void;
}) {
  const api = useApi();
  const upload = useUpload();
  const { signOut, openUserProfile } = useClerk();
  const { user: clerkUser } = useUser();
  const avatarInput = useRef<HTMLInputElement>(null);
  const bannerInput = useRef<HTMLInputElement>(null);

  const [username, setUsername] = useState(user.username);
  const [displayName, setDisplayName] = useState(user.displayName);
  const [bio, setBio] = useState(user.bio ?? "");
  const [avatarKey, setAvatarKey] = useState(user.avatarKey);
  const [bannerKey, setBannerKey] = useState(user.bannerKey);

  const [profileVisibility, setProfileVisibility] = useState(user.profileVisibility);
  const [showOnlineStatus, setShowOnlineStatus] = useState(user.showOnlineStatus);
  const [groupInvitePermission, setGroupInvitePermission] = useState(user.groupInvitePermission);
  const [shareReadReceipts, setShareReadReceipts] = useState(user.shareReadReceipts);

  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [blocked, setBlocked] = useState<MiniUser[]>([]);

  useEffect(() => {
    api.User.getBlockedUsers()
      .then((res) => {
        if (res.status === 200) setBlocked(res.body);
      })
      .catch(() => {});
  }, [api]);

  async function unblock(id: string) {
    const res = await api.User.unblockUser({ params: { id } }).catch(() => null);
    if (res && res.status === 204) setBlocked((prev) => prev.filter((u) => u.id !== id));
  }

  async function handleImage(file: File | undefined, target: "avatar" | "banner") {
    if (!file) return;
    setIsUploading(true);
    const result = await upload(file);
    setIsUploading(false);
    if (!result.ok) {
      setMessage(result.error);
      return;
    }
    if (target === "avatar") {
      setAvatarKey(result.media.storageKey);
      // Clerk is where the picture is shown in its own account widget, and
      // where most people change it. Pushing it there too keeps the two from
      // drifting apart depending on which side you uploaded from.
      try {
        await clerkUser?.setProfileImage({ file });
      } catch {
        setMessage("Uploaded, but couldn't sync the picture to your account.");
        return;
      }
    } else {
      setBannerKey(result.media.storageKey);
    }
    setMessage("Image uploaded — save to apply.");
  }

  async function save() {
    setIsSaving(true);
    setMessage(null);

    const [profileRes, settingsRes] = await Promise.all([
      api.User.updateMe({
        body: { username, displayName, bio: bio.trim() || null, avatarKey, bannerKey },
      }).catch(() => null),
      api.User.updateMySettings({
        body: { profileVisibility, showOnlineStatus, groupInvitePermission, shareReadReceipts },
      }).catch(() => null),
    ]);

    setIsSaving(false);

    if (profileRes?.status === 200 && settingsRes?.status === 200) {
      // The settings response is the later write, so it carries both changes.
      onSaved(settingsRes.body);
      onClose();
      return;
    }
    if (profileRes?.status !== 200) {
      setMessage("Couldn't save your profile — that username may already be taken.");
    } else {
      setMessage("Profile saved, but privacy settings didn't update.");
      onSaved(profileRes.body);
    }
  }

  async function deleteAccount() {
    const res = await api.User.deleteMe().catch(() => null);
    if (res && res.status === 204) await signOut({ redirectUrl: "/" });
    else setMessage("Couldn't delete your account.");
  }

  return (
    <div className="fixed inset-0 z-30 flex items-start justify-center overflow-y-auto bg-black/70 p-4 py-10">
      <div className="flex w-full max-w-[520px] flex-col gap-5 border border-border bg-surface-raised p-5">
        <div className="flex items-center justify-between">
          <h3 className="font-mono text-[0.8rem] font-bold tracking-[0.08em] uppercase">
            Update profile
          </h3>
          <button
            onClick={onClose}
            className="font-mono text-[0.8rem] text-text-faint hover:text-text"
          >
            ✕
          </button>
        </div>

        <div className="flex flex-wrap items-center gap-5">
          <div className="flex flex-col items-center gap-2">
            <RemoteAvatar storageKey={avatarKey} url={clerkUser?.imageUrl} size={56} />
            <button
              onClick={() => avatarInput.current?.click()}
              disabled={isUploading}
              className="font-mono text-[0.64rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
            >
              Avatar
            </button>
            <input
              ref={avatarInput}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={(e) => handleImage(e.target.files?.[0], "avatar")}
            />
          </div>
          <div className="flex flex-1 flex-col items-center gap-2">
            <RemoteBanner storageKey={bannerKey} className="h-14" />
            <button
              onClick={() => bannerInput.current?.click()}
              disabled={isUploading}
              className="font-mono text-[0.64rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
            >
              Banner
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

        <Field label="Username">
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
          />
        </Field>

        <Field label="Display name">
          <input
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
          />
        </Field>

        <Field label="Bio">
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            rows={3}
            maxLength={1000}
            placeholder="Tell people what you're into…"
            className="resize-none border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] placeholder:text-text-faint focus:border-accent focus:outline-none"
          />
        </Field>

        <div className="border-t border-border-soft pt-4">
          <div className="eyebrow mb-3">Privacy</div>

          <div className="flex flex-col gap-3">
            <Field label="Who can see your profile">
              <select
                value={profileVisibility}
                onChange={(e) => setProfileVisibility(e.target.value as typeof profileVisibility)}
                className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
              >
                <option value="public">Everyone</option>
                <option value="followers_only">Followers only</option>
                <option value="private">Only me</option>
              </select>
            </Field>

            <Field label="Who can add you to a group">
              <select
                value={groupInvitePermission}
                onChange={(e) =>
                  setGroupInvitePermission(e.target.value as typeof groupInvitePermission)
                }
                className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
              >
                <option value="everyone">Everyone</option>
                <option value="followers_only">People who follow me</option>
                <option value="no_one">No one</option>
              </select>
            </Field>

            <Toggle
              label="Show when you're online"
              checked={showOnlineStatus}
              onChange={setShowOnlineStatus}
            />
            <Toggle
              label="Share read receipts"
              checked={shareReadReceipts}
              onChange={setShareReadReceipts}
            />
          </div>
        </div>

        {message && <p className="font-mono text-[0.72rem] text-accent-strong">{message}</p>}

        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border-soft pt-4">
          <button
            onClick={() => openUserProfile()}
            className="font-mono text-[0.68rem] tracking-[0.05em] text-text-faint uppercase hover:text-text-muted"
          >
            Manage account ↗
          </button>

          <div className="flex gap-2">
            <button
              onClick={onClose}
              className="px-4 py-2 font-mono text-[0.72rem] tracking-[0.06em] text-text-muted uppercase"
            >
              Cancel
            </button>
            <button
              onClick={save}
              disabled={isSaving || isUploading}
              className="bg-accent px-5 py-2 font-mono text-[0.72rem] font-bold tracking-[0.06em] text-accent-text uppercase disabled:opacity-50"
            >
              {isSaving ? "Saving…" : "Save"}
            </button>
          </div>
        </div>

        <div className="border-t border-border-soft pt-4">
          <div className="eyebrow mb-3">Blocked accounts</div>
          {blocked.length === 0 ? (
            <p className="font-mono text-[0.72rem] text-text-faint">
              You haven&apos;t blocked anyone.
            </p>
          ) : (
            <div className="flex flex-col gap-2">
              {blocked.map((u) => (
                <div
                  key={u.id}
                  className="flex items-center justify-between gap-3 border border-border bg-surface px-3 py-2"
                >
                  <span className="min-w-0 truncate text-[0.84rem]">
                    @{u.username}{" "}
                    <span className="text-text-faint">{u.displayName}</span>
                  </span>
                  <button
                    onClick={() => unblock(u.id)}
                    className="shrink-0 border border-border px-3 py-1 font-mono text-[0.64rem] font-bold tracking-[0.05em] text-text-muted uppercase hover:border-accent hover:text-accent-strong"
                  >
                    Unblock
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="flex flex-wrap items-center justify-between gap-3 border border-accent/30 bg-accent/5 p-3.5">
          <div>
            <div className="text-[0.82rem] font-bold text-accent-strong">Delete account</div>
            <div className="mt-0.5 font-mono text-[0.68rem] text-text-faint">
              Permanently deletes your account and all of its data.
            </div>
          </div>
          {confirmDelete ? (
            <div className="flex gap-2">
              <button
                onClick={deleteAccount}
                className="bg-accent px-3 py-1.5 font-mono text-[0.66rem] font-bold tracking-[0.05em] text-accent-text uppercase"
              >
                Confirm
              </button>
              <button
                onClick={() => setConfirmDelete(false)}
                className="px-3 py-1.5 font-mono text-[0.66rem] tracking-[0.05em] text-text-muted uppercase"
              >
                Cancel
              </button>
            </div>
          ) : (
            <button
              onClick={() => setConfirmDelete(true)}
              className="border border-accent px-3.5 py-1.5 font-mono text-[0.66rem] font-bold tracking-[0.05em] text-accent-strong uppercase hover:bg-accent/10"
            >
              Delete
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="eyebrow">{label}</span>
      {children}
    </label>
  );
}

function Toggle({
  label,
  checked,
  onChange,
}: {
  label: string;
  checked: boolean;
  onChange: (v: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-between gap-4 border border-border bg-surface px-3.5 py-3">
      <span className="text-[0.84rem]">{label}</span>
      <button
        onClick={() => onChange(!checked)}
        aria-pressed={checked}
        aria-label={label}
        className={`relative h-6 w-11 shrink-0 ${checked ? "bg-accent" : "bg-border"}`}
      >
        <span
          className={`absolute top-1 h-4 w-4 bg-white transition-transform ${
            checked ? "translate-x-6" : "translate-x-1"
          }`}
        />
      </button>
    </div>
  );
}
