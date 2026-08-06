"use client";

import { useRef, useState } from "react";
import { useApi } from "@/lib/use-api";
import type { User } from "@nexus/zod";

export function ProfileTab({
  user,
  onSaved,
}: {
  user: User;
  onSaved: (user: User) => void;
}) {
  const api = useApi();
  const [username, setUsername] = useState(user.username);
  const [displayName, setDisplayName] = useState(user.displayName);
  const [bio, setBio] = useState(user.bio ?? "");
  const [avatarKey, setAvatarKey] = useState(user.avatarKey);
  const [bannerKey, setBannerKey] = useState(user.bannerKey);
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const avatarInputRef = useRef<HTMLInputElement>(null);
  const bannerInputRef = useRef<HTMLInputElement>(null);

  async function uploadImage(file: File): Promise<string | null> {
    setIsUploading(true);
    const presign = await api.Storage.presignUpload({
      body: { mimeType: file.type || "application/octet-stream" },
    }).catch(() => null);
    if (!presign || presign.status !== 200) {
      setIsUploading(false);
      return null;
    }
    const ok = await fetch(presign.body.uploadUrl, {
      method: "PUT",
      headers: { "Content-Type": file.type || "application/octet-stream" },
      body: file,
    })
      .then((r) => r.ok)
      .catch(() => false);
    setIsUploading(false);
    return ok ? presign.body.key : null;
  }

  async function handleAvatarChange(file: File | undefined) {
    if (!file) return;
    const key = await uploadImage(file);
    if (key) setAvatarKey(key);
  }

  async function handleBannerChange(file: File | undefined) {
    if (!file) return;
    const key = await uploadImage(file);
    if (key) setBannerKey(key);
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
      setMessage("Couldn't save. Check your username isn't taken and try again.");
    }
  }

  return (
    <div className="flex flex-col gap-5">
      <div className="flex items-center gap-4">
        <div className="flex flex-col items-center gap-1.5">
          <div className="h-16 w-16 rounded-full bg-gradient-to-br from-accent to-down" />
          <button
            onClick={() => avatarInputRef.current?.click()}
            disabled={isUploading}
            className="text-[0.72rem] font-bold text-accent-strong hover:underline"
          >
            Change avatar
          </button>
          <input
            ref={avatarInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(e) => handleAvatarChange(e.target.files?.[0])}
          />
        </div>
        <div className="flex flex-1 flex-col items-center gap-1.5">
          <div className="h-16 w-full rounded-[9px] bg-gradient-to-br from-accent to-down" />
          <button
            onClick={() => bannerInputRef.current?.click()}
            disabled={isUploading}
            className="text-[0.72rem] font-bold text-accent-strong hover:underline"
          >
            Change banner
          </button>
          <input
            ref={bannerInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(e) => handleBannerChange(e.target.files?.[0])}
          />
        </div>
      </div>

      <Field label="Username">
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
        />
      </Field>

      <Field label="Display name">
        <input
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
        />
      </Field>

      <Field label="Bio">
        <textarea
          value={bio}
          onChange={(e) => setBio(e.target.value)}
          rows={3}
          className="resize-none rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
        />
      </Field>

      <div className="flex items-center gap-3">
        <button
          onClick={save}
          disabled={isSaving || isUploading}
          className="rounded-[9px] bg-accent px-4 py-2.5 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
        >
          {isSaving ? "Saving..." : "Save changes"}
        </button>
        {message && <span className="text-[0.78rem] text-text-muted">{message}</span>}
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="text-[0.78rem] font-bold text-text-muted">{label}</span>
      {children}
    </label>
  );
}
