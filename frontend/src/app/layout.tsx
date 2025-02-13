import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import type React from "react";
import Navigation from "@/components/Navigation";
import { ClerkProvider } from "@clerk/nextjs";
const inter = Inter({ subsets: ["latin"] });
export const metadata: Metadata = {
  title: "SwipeRank - Simplify Your Applicant Ranking Process",
  description:
    "SwipeRank helps clubs and organizations efficiently review and rank applicants with an intuitive, swipe-based interface.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body className={inter.className}>
          <header className="p-4 border-b">
            <Navigation />
          </header>
          {children}
          <footer className="p-4 border-t mt-auto">
            <div className="container mx-auto text-center">
              <a
                href="https://github.com/your-repo-here"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-600 hover:text-gray-800"
              >
                View on GitHub
              </a>
            </div>
          </footer>
        </body>
      </html>
    </ClerkProvider>
  );
}
