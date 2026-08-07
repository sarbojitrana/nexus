import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { SettingsApp } from "@/components/settings/settings-app";

export default function SettingsPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/settings" />
      <SettingsApp />
    </div>
  );
}
