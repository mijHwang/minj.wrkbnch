# minj.wrkbnch

A workbench for infrastructure concepts, each one with a working, interactive demo.

## Live

- Hub: [link once deployed]
- LRU Cache demo: [link once deployed]
- URL Shortener demo: [link once deployed]

## Structure

One shared React frontend, routing to each project's own page. Each project has
its own backend, deployed as its own service — not one monolith.

minj.wrkbnch/
├── frontend/          React + Vite — hub page + all project demo pages
├── CacheLRU/          Go backend for the LRU cache project
└── url-shortener/     Go backend for the URL shortener project

## Projects

### LRU Cache Server (CacheLRU/)

An LRU cache built from scratch in Go — no external cache library.

- Doubly-linked list + hash map, giving O(1) Get/Set
- TTL-based expiry on each entry
- sync.Mutex guarding concurrent access
- Eviction logic covered by tests

API
- POST /set — { key, value, ttl_seconds }
- GET /get?key=...
- GET /all — current cache state, used by the live dashboard

### URL Shortener (url-shortener/)

A URL shortener backed by Postgres (Neon) — long URLs in, short codes out.

- Auto-incrementing row id encoded to base62 for the short code — no randomness,
  no collisions possible, nothing stored beyond the original URL
- Short codes decoded back to an id on visit, looked up, and redirected with a
  real HTTP redirect
- No expiry, no deduplication — each submission gets its own short code

API
- POST /shorten — { url } → { short_url, long_url }
- GET /{code} — redirects to the original URL

## Running locally

LRU Cache backend
    cd CacheLRU
    go run main.go

URL Shortener backend
    cd url-shortener
    go run main.go

Create url-shortener/.env with:
    DATABASE_URL=your-neon-connection-string

Frontend
    cd frontend
    npm install
    npm run dev

Create frontend/.env with:
    VITE_CACHE_API_URL=http://localhost:8080
    VITE_SHORTENER_API_URL=http://localhost:8081

## Stack

- Frontend: React, Vite
- Backends: Go (net/http, no framework)
- Database: Postgres (Neon, free tier, no expiry)
- Deploy: Render (one service per backend, one static site for the frontend)