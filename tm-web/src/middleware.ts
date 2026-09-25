import { NextRequest, NextResponse } from "next/server";
import { TOKEN, TOKEN_ORG } from "./redux";

const publicRoute = [
  "/login",
  "/qr-scan",
  "/sign-in",
  "/organizations/login",
  "/invite-expired",
  "/organizations/invite",
];

export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  // Allow any path that starts with the excluded ones
  const isPublicRoute = publicRoute.some((url) => pathname.startsWith(url));

  if (!isPublicRoute) {
    if (pathname.startsWith("/organisation")) {
      // check org token cookie
      const orgToken = request.cookies.get(TOKEN_ORG);
      if (!orgToken) {
        return NextResponse.redirect(
          new URL("/organizations/login", request.url)
        );
      }
    } else {
      // check user token cookie
      const token = request.cookies.get(TOKEN);
      if (!token) {
        return NextResponse.redirect(new URL("/sign-in", request.url));
      }
    }
  }
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};
