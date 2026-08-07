import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { CommunityDetail } from "@/components/communities/community-detail";

export default async function CommunityPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;

  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/communities" />
      <CommunityDetail slug={slug} />
    </div>
  );
}
