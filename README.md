# Nexus

A community-first social platform — posts live inside communities you choose, your feed blends what's trending with the people and spaces you actually follow, and every community governs itself with its own moderators, roles and report queue.

Go + Echo API, Next.js frontend, contracts shared between them in a Turborepo monorepo.

---

## What it does

**Posts and communities** — create communities, join or follow them, post with images and video, vote, comment, and reply to replies. Each community has members, roles (`admin` / `moderator` / `member`), a report queue, bans, and its own settings.

**Chat** — 1:1 and group conversations over WebSockets, with live presence, read receipts, file attachments, and an invite → accept/reject flow for groups.

**Search** — full-text across posts, people and communities via OpenSearch, with a Postgres fallback so search keeps working when the index is unavailable or empty.

**Notifications** — in-app feed for follows, group invitations and messages, delivered live over the same socket the chat uses.

**Privacy** — per-user profile visibility, online-status visibility, group-invite permissions, read-receipt sharing, and user blocking.

**Emails** — welcome, new sign-in, and password-changed notices, sent through a background job queue.

---

## Stack

| Layer | Choice |
|---|---|
| API | Go 1.25, Echo, pgx/v5 (no ORM) |
| Database | PostgreSQL, migrations via `tern` |
| Auth | Clerk (passwords and OAuth never touch this codebase) |
| Cache + jobs | Redis — `asynq` queue plus a cache-aside layer |
| Search | OpenSearch, with a Postgres `ILIKE` fallback |
| Object storage | Any S3-compatible provider (built against Cloudflare R2) |
| Real-time | gorilla/websocket |
| Frontend | Next.js 16 (App Router), React 19, Tailwind v4 |
| Contracts | Zod schemas → ts-rest, shared by both sides |
| Tooling | Turborepo, Bun |
| Observability | zerolog, New Relic |

---

## Layout

```
apps/
  backend/          Go API — repository → service → handler → router
  web/              Next.js app
packages/
  zod/              Request/response schemas, single source of truth
  openapi/          ts-rest contracts + the typed API client
```

The frontend imports `@nexus/openapi` as a workspace package, so an API change that breaks the contract surfaces as a TypeScript error at build time rather than a runtime failure.

---

## Running it locally

**Prerequisites:** Docker, Go 1.25+, Bun.

```bash
# 1. Infrastructure (Postgres, Redis, OpenSearch)
docker compose up -d postgres redis opensearch

# 2. Backend config
cp apps/backend/.env.example apps/backend/.env
# fill in NEXUS_AUTH.SECRET_KEY, NEXUS_AUTH.WEBHOOK_SECRET and the SMTP values

# 3. Frontend config
cp apps/web/.env.example apps/web/.env
# fill in the Clerk keys; NEXT_PUBLIC_API_URL defaults to http://localhost:8080

# 4. Install and run
bun install
cd apps/backend && go run ./cmd/nexus     # API
bun run dev --filter=@nexus/web            # frontend
```

Migrations run automatically at startup whenever `NEXUS_PRIMARY.ENV` is anything other than `local`; apply them by hand for local development.

**Environment variables** use `NEXUS_`-prefixed keys with **dots** for nesting — `NEXUS_DATABASE.HOST`, not `NEXUS_DATABASE_HOST`. `NEXUS_SERVER.CORS_ALLOWED_ORIGINS` accepts a comma-separated list.

Search and object storage are both optional: leave `NEXUS_SEARCH.URL` or the `NEXUS_AWS.*` values empty and those features degrade to disabled rather than blocking startup.

---

## Common commands

```bash
bun run build          # build every workspace
bun run typecheck      # tsc across the monorepo
go build ./...         # from apps/backend
go vet ./...
```

---

## Deployment

Backend on Render (Docker, `render.yaml`), frontend on Vercel (root directory `apps/web`), Postgres on Neon, Redis on Render Key Value, search on Bonsai, object storage on Cloudflare R2. A GitHub Actions workflow pings `/status` on a schedule to keep the free-tier instance warm.

Two things that bite on a fresh deploy:

- **Clerk webhook** — point it at `https://<backend>/webhooks/clerk` for `user.created`, `user.updated`, `user.deleted` and `session.created`. Without it, profile changes and deletions won't sync. Accounts still self-provision on first authenticated request, so a missing webhook won't lock anyone out.
- **R2 CORS** — the bucket must allow `PUT` from the frontend's origin, or uploads fail in the browser even with valid credentials.

---

## Status

Working and deployed. Some known gaps:

- The WebSocket hub is in-memory, so real-time only works correctly on a single backend instance. Horizontal scaling needs Redis pub/sub fan-out first.
- Search relevance falls back to substring matching when OpenSearch is unavailable — no ranking or typo tolerance in that mode.
- Email goes out over Gmail SMTP, which is fine at low volume but has send caps and no deliverability guarantees. Swapping in a transactional provider touches one file.
- There are no automated tests yet.
