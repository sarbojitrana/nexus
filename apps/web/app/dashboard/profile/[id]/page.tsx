import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { ProfileView } from "@/components/profile/profile-view";

export default async function UserProfilePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;

  return (
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/profile" />
      <ProfileView userId={id} />
    </div>
  );
}
