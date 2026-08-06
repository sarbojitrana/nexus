import { Suspense } from "react";
import { DashboardSidebar } from "@/components/dashboard-sidebar";
import { ChatApp } from "@/components/chat/chat-app";

export default function MessagesPage() {
  return (
    <div className="grid min-h-dvh grid-cols-1 lg:grid-cols-[232px_1fr]">
      <DashboardSidebar active="/dashboard/messages" />
      <Suspense>
        <ChatApp />
      </Suspense>
    </div>
  );
}
