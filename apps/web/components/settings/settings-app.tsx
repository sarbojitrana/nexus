"use client";

import { useEffect, useState } from "react";
import { useApi } from "@/lib/use-api";
import { ProfileTab } from "@/components/settings/profile-tab";
import { PrivacyTab } from "@/components/settings/privacy-tab";
import { AccountTab } from "@/components/settings/account-tab";
import type { User } from "@nexus/zod";

const TABS = ["Profile", "Privacy", "Account"] as const;
type Tab = (typeof TABS)[number];

export function SettingsApp() {
  const api = useApi();
  const [tab, setTab] = useState<Tab>("Profile");
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    api.User.getMe().then((res) => {
      if (res.status === 200) setUser(res.body);
      setIsLoading(false);
    });
  }, [api]);

  return (
    <div className="mx-auto flex w-full max-w-[620px] flex-col gap-6 px-6 py-8">
      <h1 className="font-display text-[1.4rem] font-extrabold">Settings</h1>

      <div className="flex gap-1 border-b border-border-soft">
        {TABS.map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`rounded-t-[9px] px-4 py-2.5 text-[0.84rem] font-bold ${
              tab === t
                ? "border-b-2 border-accent text-accent-strong"
                : "text-text-faint hover:text-text-muted"
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      {isLoading && <p className="text-[0.84rem] text-text-faint">Loading...</p>}

      {!isLoading && user && tab === "Profile" && (
        <ProfileTab user={user} onSaved={setUser} />
      )}
      {!isLoading && user && tab === "Privacy" && (
        <PrivacyTab user={user} onSaved={setUser} />
      )}
      {!isLoading && tab === "Account" && <AccountTab />}
    </div>
  );
}
