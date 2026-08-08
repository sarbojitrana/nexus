"use client";

// Small wrapper around the Notifications API: asks once, stays quiet if the
// user said no, and never throws in browsers that don't support it.
export async function notify(title: string, body?: string) {
  if (typeof window === "undefined" || !("Notification" in window)) return;

  let permission = Notification.permission;
  if (permission === "default") {
    permission = await Notification.requestPermission().catch(() => "denied" as NotificationPermission);
  }
  if (permission !== "granted") return;

  try {
    new Notification(title, { body, icon: "/icon.svg" });
  } catch {
    // Some browsers only allow notifications from a service worker; a failure
    // here isn't worth surfacing since it's a nicety, not the action itself.
  }
}
