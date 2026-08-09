"use client";

import { useState } from "react";

// Uses the native share sheet where it exists (mobile, some desktop browsers)
// and falls back to copying the link, which is what most desktop users get.
export function ShareButton({
  path,
  title,
  label = "share",
  className = "",
}: {
  /** Path within the app, e.g. /dashboard/posts/123 */
  path: string;
  title?: string;
  label?: string;
  className?: string;
}) {
  const [copied, setCopied] = useState(false);

  async function share() {
    const url = typeof window === "undefined" ? path : new URL(path, window.location.origin).toString();

    if (navigator.share) {
      try {
        await navigator.share({ title, url });
        return;
      } catch {
        // A cancelled share sheet lands here too, so fall through to copying
        // rather than showing an error for what may have been a deliberate
        // dismissal.
      }
    }

    try {
      await navigator.clipboard.writeText(url);
      setCopied(true);
      setTimeout(() => setCopied(false), 1600);
    } catch {
      // Clipboard access can be denied; nothing useful to say beyond leaving
      // the label unchanged.
    }
  }

  return (
    <button
      onClick={share}
      aria-label="Share"
      className={`hover:text-accent-strong ${className}`}
    >
      {copied ? "link copied" : label}
    </button>
  );
}
