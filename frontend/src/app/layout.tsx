import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import type React from "react";
import Navigation from "@/components/Navigation";
import { ClerkProvider } from "@clerk/nextjs";
import Footer from "@/components/Footer";
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
          <Footer />
        </body>
      </html>
    </ClerkProvider>
  );
}
