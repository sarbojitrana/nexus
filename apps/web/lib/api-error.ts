// ts-rest hands back non-2xx responses as {status, body} rather than throwing,
// and the Go error handler always shapes the body as {code, message, status}.
// Surfacing that message beats a hardcoded guess at what went wrong.
export function apiErrorMessage(res: { status: number; body: unknown } | null, fallback: string) {
  if (!res) return `${fallback} (couldn't reach the server)`;

  const body = res.body as { message?: unknown; code?: unknown } | null;
  const message = typeof body?.message === "string" ? body.message : null;
  const code = typeof body?.code === "string" ? body.code : null;

  if (message && code && code !== "INTERNAL_SERVER_ERROR") return message;
  if (message && !code) return message;

  if (res.status === 401) return "You're signed out — sign in and try again.";
  if (res.status === 500) return `${fallback} (server error — check the backend logs)`;
  return message ?? fallback;
}
