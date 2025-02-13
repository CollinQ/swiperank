import { clerkMiddleware } from "@clerk/nextjs/server";

// Apply Clerk authentication middleware
export default clerkMiddleware();

// Ensure middleware runs on all routes
export const config = {
  matcher: "/((?!api|_next|.*\\..*).*)",
};
