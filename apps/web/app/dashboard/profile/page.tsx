import { redirect } from "next/navigation";
import { getServerUserId } from "@/lib/api-server";

export default async function MyProfilePage() {
  const userId = await getServerUserId();
  if (!userId) redirect("/sign-in");
  redirect(`/dashboard/profile/${userId}`);
}
