"use client";

import { useState } from "react";
import { useClerk, useUser } from "@clerk/nextjs";
import { useApi } from "@/lib/use-api";

export function AccountTab() {
  const api = useApi();
  const { signOut } = useClerk();
  const { user } = useUser();
  const [confirming, setConfirming] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [passwordMessage, setPasswordMessage] = useState<string | null>(null);

  async function deleteAccount() {
    setIsDeleting(true);
    const res = await api.User.deleteMe().catch(() => null);
    if (res && res.status === 204) {
      await signOut({ redirectUrl: "/" });
    } else {
      setIsDeleting(false);
    }
  }

  async function changePassword() {
    setPasswordMessage(null);
    if (newPassword.length < 8) {
      setPasswordMessage("New password must be at least 8 characters.");
      return;
    }
    if (newPassword !== confirmPassword) {
      setPasswordMessage("New passwords don't match.");
      return;
    }
    if (!user) return;

    setIsChangingPassword(true);
    try {
      await user.updatePassword({ currentPassword, newPassword, signOutOfOtherSessions: false });
      await api.User.notifyPasswordChanged().catch(() => {});
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      setPasswordMessage("Password updated.");
    } catch {
      setPasswordMessage("Couldn't update your password. Check your current password and try again.");
    } finally {
      setIsChangingPassword(false);
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3 rounded-2xl border border-border bg-surface p-4">
        <span className="text-[0.84rem] font-bold">Change password</span>
        <input
          type="password"
          value={currentPassword}
          onChange={(e) => setCurrentPassword(e.target.value)}
          placeholder="Current password"
          className="rounded-[9px] border border-border bg-bg px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />
        <input
          type="password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          placeholder="New password"
          className="rounded-[9px] border border-border bg-bg px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />
        <input
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          placeholder="Confirm new password"
          className="rounded-[9px] border border-border bg-bg px-3.5 py-2.5 text-[0.86rem] text-text placeholder:text-text-faint focus:border-accent focus:outline-none"
        />
        <div className="flex items-center gap-3">
          <button
            onClick={changePassword}
            disabled={isChangingPassword || !currentPassword || !newPassword}
            className="self-start rounded-[9px] bg-accent px-4 py-2 text-[0.8rem] font-bold text-accent-text hover:bg-accent-strong disabled:opacity-50"
          >
            {isChangingPassword ? "Updating..." : "Update password"}
          </button>
          {passwordMessage && (
            <span className="text-[0.76rem] text-text-muted">{passwordMessage}</span>
          )}
        </div>
      </div>

      <div className="flex items-center justify-between rounded-2xl border border-border bg-surface p-4">
        <div className="flex flex-col gap-0.5">
          <span className="text-[0.84rem] font-bold">Log out</span>
          <span className="text-[0.74rem] text-text-faint">
            Sign out of Nexus on this device.
          </span>
        </div>
        <button
          onClick={() => signOut({ redirectUrl: "/" })}
          className="rounded-[9px] border border-border px-4 py-2 text-[0.8rem] font-bold text-text-muted hover:bg-surface"
        >
          Log out
        </button>
      </div>

      <div className="flex items-center justify-between gap-4 rounded-2xl border border-accent/30 bg-accent/5 p-4">
        <div className="flex flex-col gap-0.5">
          <span className="text-[0.84rem] font-bold text-accent-strong">Delete account</span>
          <span className="text-[0.74rem] text-text-faint">
            Permanently deletes your account and all of its data. This can't be undone.
          </span>
        </div>
        {!confirming ? (
          <button
            onClick={() => setConfirming(true)}
            className="shrink-0 rounded-[9px] border border-accent px-4 py-2 text-[0.8rem] font-bold text-accent-strong hover:bg-accent/10"
          >
            Delete
          </button>
        ) : (
          <div className="flex shrink-0 gap-1.5">
            <button
              onClick={deleteAccount}
              disabled={isDeleting}
              className="rounded-[9px] bg-accent px-3 py-2 text-[0.78rem] font-bold text-accent-text disabled:opacity-50"
            >
              {isDeleting ? "Deleting..." : "Confirm"}
            </button>
            <button
              onClick={() => setConfirming(false)}
              className="rounded-[9px] px-3 py-2 text-[0.78rem] font-bold text-text-muted hover:bg-surface"
            >
              Cancel
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
