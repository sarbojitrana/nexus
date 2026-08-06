"use client";

import { useState } from "react";
import { useApi } from "@/lib/use-api";
import type { User } from "@nexus/zod";

export function PrivacyTab({
  user,
  onSaved,
}: {
  user: User;
  onSaved: (user: User) => void;
}) {
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
      setMessage("Couldn't save. Try again.");
    }
  }

  return (
    <div className="flex flex-col gap-5">
      <SelectField
        label="Who can see your profile"
        value={profileVisibility}
        onChange={setProfileVisibility}
        options={[
          { value: "public", label: "Everyone" },
          { value: "followers_only", label: "Followers only" },
          { value: "private", label: "Only me" },
        ]}
      />

      <ToggleField
        label="Show when you're online"
        description="Other people you chat with can see your online status and read receipts."
        checked={showOnlineStatus}
        onChange={setShowOnlineStatus}
      />

      <SelectField
        label="Who can add you to a group"
        value={groupInvitePermission}
        onChange={setGroupInvitePermission}
        options={[
          { value: "everyone", label: "Everyone" },
          { value: "followers_only", label: "People I follow me" },
          { value: "no_one", label: "No one" },
        ]}
      />

      <ToggleField
        label="Share read receipts"
        description="Let people you message see when you've read their messages."
        checked={shareReadReceipts}
        onChange={setShareReadReceipts}
      />

      <div className="flex items-center gap-3">
        <button
          onClick={save}
          disabled={isSaving}
          className="rounded-[9px] bg-accent px-4 py-2.5 text-[0.82rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
        >
          {isSaving ? "Saving..." : "Save changes"}
        </button>
        {message && <span className="text-[0.78rem] text-text-muted">{message}</span>}
      </div>
    </div>
  );
}

function SelectField<T extends string>({
  label,
  value,
  onChange,
  options,
}: {
  label: string;
  value: T;
  onChange: (v: T) => void;
  options: { value: T; label: string }[];
}) {
  return (
    <label className="flex flex-col gap-1.5">
      <span className="text-[0.78rem] font-bold text-text-muted">{label}</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value as T)}
        className="rounded-[9px] border border-border bg-surface px-3.5 py-2.5 text-[0.86rem] text-text focus:border-accent focus:outline-none"
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function ToggleField({
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
    <div className="flex items-center justify-between gap-4 rounded-2xl border border-border bg-surface p-4">
      <div className="flex flex-col gap-0.5">
        <span className="text-[0.84rem] font-bold">{label}</span>
        <span className="text-[0.74rem] text-text-faint">{description}</span>
      </div>
      <button
        onClick={() => onChange(!checked)}
        className={`relative h-6 w-11 shrink-0 rounded-full transition-colors ${
          checked ? "bg-accent" : "bg-border"
        }`}
        aria-pressed={checked}
      >
        <span
          className={`absolute top-0.5 h-5 w-5 rounded-full bg-white transition-transform ${
            checked ? "translate-x-[22px]" : "translate-x-0.5"
          }`}
        />
      </button>
    </div>
  );
}
