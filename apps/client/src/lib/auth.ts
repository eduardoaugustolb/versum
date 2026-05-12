export function getSessionCookieName() {
  return process.env.NEXT_PUBLIC_API_URL?.startsWith("https")
    ? "__Host-session"
    : "session";
}
