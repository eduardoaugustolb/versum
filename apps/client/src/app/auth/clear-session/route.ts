import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getSessionCookieName } from "@/lib/auth";

export async function GET() {
  (await cookies()).delete(getSessionCookieName());
  redirect("/login");
}
