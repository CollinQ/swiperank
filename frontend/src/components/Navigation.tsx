"use client";

import { Button } from "@/components/ui/button";
import { SignInButton, SignUpButton, UserButton, useUser } from "@clerk/nextjs";
import Link from "next/link";

export default function Navigation() {
  const { isSignedIn } = useUser(); // Get user authentication status

  return (
    <header className="py-4 px-6 md:px-10 flex justify-between items-center bg-white shadow-sm">
      <div className="text-2xl font-bold text-blue-600">SwipeRank</div>
      <nav>
        <ul className="flex space-x-4">
          <li>
            <Button variant="ghost" asChild>
              <Link href="#features">Features</Link>
            </Button>
          </li>
          <li>
            <Button variant="ghost" asChild>
              <Link href="#how-it-works">How It Works</Link>
            </Button>
          </li>

          {isSignedIn ? (
            // If logged in, show User Avatar
            <li>
              <UserButton afterSignOutUrl="/" />
            </li>
          ) : (
            // If not logged in, show Sign In & Sign Up buttons
            <>
              <li>
                <Button variant="outline" asChild>
                  <SignInButton />
                </Button>
              </li>
              <li>
                <Button variant="default" asChild>
                  <SignInButton />
                </Button>
              </li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
}
