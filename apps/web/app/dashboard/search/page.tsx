import { Suspense } from "react";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { HudCorners } from "@/components/logo";
import { SearchResults } from "@/components/search/search-results";

export default async function SearchPage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string }>;
}) {
  const { q } = await searchParams;

  return (
    <div className="grid h-dvh grid-cols-1 overflow-hidden lg:grid-cols-[236px_1fr]">
      <HudCorners />
      <DashboardSidebar active="/dashboard/search" />
      <div className="mx-auto flex min-h-0 w-full max-w-[880px] flex-col">
        <Suspense>
          <SearchResults initialQuery={q ?? ""} />
        </Suspense>
      </div>
    </div>
  );
}
