"use client";

import { useCallback, useEffect, useState } from "react";
import { useClerk } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";
import { useUpload } from "@/lib/use-upload";
import type { User } from "@nexus/zod";

const TABS = ["Profile", "Privacy", "Account"] as const;
type Tab = (typeof TABS)[number];

export function SettingsApp() {
  const api = useApi();
  const [tab, setTab] = useState<Tab>("Profile");
  const [user, setUser] = useState<User | null>(null);
  const [status, setStatus] = useState<"loading" | "ready" | "error">("loading");

  const load = useCallback(async () => {
    setStatus("loading");
    const res = await api.User.getMe().catch(() => null);
    // The previous version left this stuck on "Loading..." forever whenever
    // the request failed -- now a failure is its own visible state.
    if (res && res.status === 200) {
      setUser(res.body);
      setStatus("ready");
    } else {
      setStatus("error");
    }
  }, [api]);

  useEffect(() => {
    load();
  }, [load]);

  return (
    <div className="mx-auto flex w-full max-w-[640px] flex-col gap-5 px-6 py-6">
      <div>
        <h1 className="font-display text-[1.3rem] font-extrabold">Settings</h1>
        <p className="eyebrow mt-1">profile · privacy · account</p>
      </div>

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

      {status === "loading" && (
        <p className="font-mono text-[0.76rem] text-text-faint">loading…</p>
      )}

      {status === "error" && (
        <div className="border border-border bg-surface p-6 text-center">
          <p className="text-[0.88rem] text-text-muted">Couldn&apos;t load your settings.</p>
          <button
            onClick={load}
            className="mt-3 border border-border px-4 py-2 font-mono text-[0.7rem] tracking-[0.06em] text-text-muted uppercase hover:border-accent/40"
          >
            Retry
          </button>
        </div>
      )}

      {status === "ready" && user && tab === "Profile" && (
        <ProfileTab user={user} onSaved={setUser} />
      )}
      {status === "ready" && user && tab === "Privacy" && (
        <PrivacyTab user={user} onSaved={setUser} />
      )}
      {status === "ready" && tab === "Account" && <AccountTab />}
    </div>
  );
}

function ProfileTab({ user, onSaved }: { user: User; onSaved: (u: User) => void }) {
  const api = useApi();
  const upload = useUpload();
  const [username, setUsername] = useState(user.username);
  const [displayName, setDisplayName] = useState(user.displayName);
  const [bio, setBio] = useState(user.bio ?? "");
  const [avatarKey, setAvatarKey] = useState(user.avatarKey);
  const [bannerKey, setBannerKey] = useState(user.bannerKey);
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  async function pickImage(target: "avatar" | "banner") {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = "image/*";
    input.onchange = async () => {
      const file = input.files?.[0];
      if (!file) return;
      setIsUploading(true);
      const uploaded = await upload(file);
      setIsUploading(false);
      if (!uploaded) {
        setMessage("Upload failed — is storage configured?");
        return;
      }
      if (target === "avatar") setAvatarKey(uploaded.storageKey);
      else setBannerKey(uploaded.storageKey);
      setMessage("Image uploaded — save to apply.");
    };
    input.click();
  }

  async function save() {
    setIsSaving(true);
    setMessage(null);
    const res = await api.User.updateMe({
      body: { username, displayName, bio: bio || null, avatarKey, bannerKey },
    }).catch(() => null);
    setIsSaving(false);
    if (res && res.status === 200) {
      onSaved(res.body);
      setMessage("Saved.");
    } else {
      setMessage("Couldn't save — that username may already be taken.");
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-5">
        <div className="flex flex-col items-center gap-2">
          <span className="h-16 w-16 bg-accent" />
          <button
            onClick={() => pickImage("avatar")}
            disabled={isUploading}
            className="font-mono text-[0.66rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
          >
            Change avatar
          </button>
        </div>
        <div className="flex flex-1 flex-col items-center gap-2">
          <span className="h-16 w-full bg-gradient-to-r from-accent to-down" />
          <button
            onClick={() => pickImage("banner")}
            disabled={isUploading}
            className="font-mono text-[0.66rem] tracking-[0.05em] text-accent-strong uppercase hover:underline disabled:opacity-50"
          >
            Change banner
          </button>
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
          className="resize-none border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
        />
      </Field>

      <div className="flex items-center gap-3">
        <button
          onClick={save}
          disabled={isSaving || isUploading}
          className="bg-accent px-5 py-2.5 font-mono text-[0.72rem] font-bold tracking-[0.06em] text-accent-text uppercase disabled:opacity-50"
        >
          {isSaving ? "Saving…" : "Save changes"}
        </button>
        {message && <span className="font-mono text-[0.72rem] text-text-muted">{message}</span>}
      </div>
    </div>
  );
}

function PrivacyTab({ user, onSaved }: { user: User; onSaved: (u: User) => void }) {
  const api = useApi();
  const [profileVisibility, setProfileVisibility] = useState(user.profileVisibility);
  const [showOnlineStatus, setShowOnlineStatus] = useState(user.showOnlineStatus);
  const [groupInvitePermission, setGroupInvitePermission] = useState(user.groupInvitePermission);
  const [shareReadReceipts, setShareReadReceipts] = useState(user.shareReadReceipts);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  async function save() {
    setIsSaving(true);
    setMessage(null);
    const res = await api.User.updateMySettings({
      body: { profileVisibility, showOnlineStatus, groupInvitePermission, shareReadReceipts },
    }).catch(() => null);
    setIsSaving(false);
    if (res && res.status === 200) {
      onSaved(res.body);
      setMessage("Saved.");
    } else {
      setMessage("Couldn't save changes.");
    }
  }

  return (
    <div className="flex flex-col gap-4">
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

      <Toggle
        label="Show when you're online"
        description="People you chat with can see your online status and read receipts."
        checked={showOnlineStatus}
        onChange={setShowOnlineStatus}
      />

      <Field label="Who can add you to a group">
        <select
          value={groupInvitePermission}
          onChange={(e) => setGroupInvitePermission(e.target.value as typeof groupInvitePermission)}
          className="border border-border bg-surface px-3.5 py-2.5 text-[0.88rem] focus:border-accent focus:outline-none"
        >
          <option value="everyone">Everyone</option>
          <option value="followers_only">People who follow me</option>
          <option value="no_one">No one</option>
        </select>
      </Field>

      <Toggle
        label="Share read receipts"
        description="Let people see when you've read their messages."
        checked={shareReadReceipts}
        onChange={setShareReadReceipts}
      />

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
    </div>
  );
}

function AccountTab() {
  const api = useApi();
  const { signOut, openUserProfile } = useClerk();
  const [confirming, setConfirming] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  async function deleteAccount() {
    setIsDeleting(true);
    const res = await api.User.deleteMe().catch(() => null);
    if (res && res.status === 204) await signOut({ redirectUrl: "/" });
    else setIsDeleting(false);
  }

  return (
    <div className="flex flex-col gap-4">
      <Row
        title="Password & security"
        description="Change your password, manage sessions and connected accounts."
        action={
          <button
            onClick={() => openUserProfile()}
            className="border border-border px-4 py-2 font-mono text-[0.7rem] font-bold tracking-[0.05em] text-text-muted uppercase hover:border-accent/40"
          >
            Manage
          </button>
        }
      />

      <Row
        title="Log out"
        description="Sign out of Nexus on this device."
        action={
          <button
            onClick={() => signOut({ redirectUrl: "/" })}
            className="border border-border px-4 py-2 font-mono text-[0.7rem] font-bold tracking-[0.05em] text-text-muted uppercase hover:border-accent/40"
          >
            Log out
          </button>
        }
      />

      <div className="flex flex-wrap items-center justify-between gap-3 border border-accent/30 bg-accent/5 p-4">
        <div>
          <div className="text-[0.86rem] font-bold text-accent-strong">Delete account</div>
          <div className="mt-0.5 font-mono text-[0.7rem] text-text-faint">
            Permanently deletes your account and all of its data.
          </div>
        </div>
        {confirming ? (
          <div className="flex gap-2">
            <button
              onClick={deleteAccount}
              disabled={isDeleting}
              className="bg-accent px-3 py-2 font-mono text-[0.68rem] font-bold tracking-[0.05em] text-accent-text uppercase disabled:opacity-50"
            >
              {isDeleting ? "Deleting…" : "Confirm"}
            </button>
            <button
              onClick={() => setConfirming(false)}
              className="px-3 py-2 font-mono text-[0.68rem] tracking-[0.05em] text-text-muted uppercase"
            >
              Cancel
            </button>
          </div>
        ) : (
          <button
            onClick={() => setConfirming(true)}
            className="border border-accent px-4 py-2 font-mono text-[0.68rem] font-bold tracking-[0.05em] text-accent-strong uppercase hover:bg-accent/10"
          >
            Delete
          </button>
        )}
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

function Row({
  title,
  description,
  action,
}: {
  title: string;
  description: string;
  action: React.ReactNode;
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 border border-border bg-surface p-4">
      <div>
        <div className="text-[0.86rem] font-bold">{title}</div>
        <div className="mt-0.5 font-mono text-[0.7rem] text-text-faint">{description}</div>
      </div>
      {action}
    </div>
  );
}

function Toggle({
  label,
  description,
  checked,
  onChange,
}: {
  label: string;
  description: string;
  checked: boolean;
  onChange: (v: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-between gap-4 border border-border bg-surface p-4">
      <div>
        <div className="text-[0.86rem] font-bold">{label}</div>
        <div className="mt-0.5 font-mono text-[0.7rem] text-text-faint">{description}</div>
      </div>
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
