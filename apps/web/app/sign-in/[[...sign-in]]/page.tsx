import { SignIn } from "@clerk/nextjs";
import { Logo } from "@/components/logo";

export default function SignInPage() {
  return (
    <main className="flex min-h-dvh flex-col">
      <nav className="flex items-center justify-between px-[6vw] py-5">
        <Logo />
        <span className="font-mono text-[0.72rem] text-text-faint">welcome back</span>
      </nav>
      <div className="grid flex-1 place-items-center px-5 py-[6vh]">
        <SignIn
          fallbackRedirectUrl="/dashboard"
          appearance={{
            elements: {
              card: "shadow-[0_24px_60px_-20px_rgba(0,0,0,0.6)] border border-border",
            },
          }}
        />
      </div>
    </main>
  );
}
