import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { NotificationsApp } from "@/components/notifications/notifications-app";

export default function NotificationsPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr]">
      <DashboardSidebar active="/dashboard/notifications" />
      <NotificationsApp />
    </div>
  );
}
