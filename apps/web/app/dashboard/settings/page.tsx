import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { SettingsApp } from "@/components/settings/settings-app";

export default function SettingsPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr]">
      <DashboardSidebar active="/dashboard/settings" />
      <SettingsApp />
    </div>
  );
}
