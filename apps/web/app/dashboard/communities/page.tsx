import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { CommunitiesApp } from "@/components/communities/communities-app";

export default function CommunitiesPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr]">
      <DashboardSidebar active="/dashboard/communities" />
      <CommunitiesApp />
    </div>
  );
}
