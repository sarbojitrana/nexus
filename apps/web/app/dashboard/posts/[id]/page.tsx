import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { PostDetail } from "@/components/posts/post-detail";

export default async function PostPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;

  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard" />
      <PostDetail postId={id} />
    </div>
  );
}
