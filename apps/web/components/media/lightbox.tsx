"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";

type LightboxItem = { url: string; kind: "image" | "video" };

const LightboxContext = createContext<(item: LightboxItem) => void>(() => {});

/** Opens any image or video full screen. Provided once at the dashboard layout
 *  so avatars, banners and post media all share one overlay. */
export function useLightbox() {
  return useContext(LightboxContext);
}

export function LightboxProvider({ children }: { children: React.ReactNode }) {
  const [item, setItem] = useState<LightboxItem | null>(null);

  const open = useCallback((next: LightboxItem) => setItem(next), []);

  useEffect(() => {
    if (!item) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setItem(null);
    };
    document.addEventListener("keydown", onKey);
    // Stop the page behind the overlay from scrolling.
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = previous;
    };
  }, [item]);

  return (
    <LightboxContext.Provider value={open}>
      {children}

      {item && (
        <div
          role="dialog"
          aria-modal="true"
          aria-label="Media viewer"
          onClick={() => setItem(null)}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/90 p-4"
        >
          <button
            onClick={() => setItem(null)}
            aria-label="Close"
            className="absolute top-4 right-5 font-mono text-[1.1rem] text-white/70 hover:text-white"
          >
            ✕
          </button>

          {item.kind === "video" ? (
            <video
              src={item.url}
              controls
              autoPlay
              onClick={(e) => e.stopPropagation()}
              className="max-h-full max-w-full"
            />
          ) : (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={item.url}
              alt=""
              onClick={(e) => e.stopPropagation()}
              className="max-h-full max-w-full object-contain"
            />
          )}
        </div>
      )}
    </LightboxContext.Provider>
  );
}
