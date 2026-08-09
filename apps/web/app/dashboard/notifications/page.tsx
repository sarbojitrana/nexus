import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { NotificationsApp } from "@/components/notifications/notifications-app";

export default function NotificationsPage() {
  return (
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <DashboardSidebar active="/dashboard/notifications" />
      <NotificationsApp />
    </div>
  );
}
