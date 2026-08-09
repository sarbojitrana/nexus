import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { CommunitiesApp } from "@/components/communities/communities-app";

export default function CommunitiesPage() {
  return (
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <DashboardSidebar active="/dashboard/communities" />
      <CommunitiesApp />
    </div>
  );
}
