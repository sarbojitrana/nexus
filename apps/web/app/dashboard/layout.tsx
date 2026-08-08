import { LightboxProvider } from "@/components/media/lightbox";

// One overlay for the whole dashboard, so avatars, banners and post media can
// all open full screen without each screen wiring its own viewer.
export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return <LightboxProvider>{children}</LightboxProvider>;
}
