# minj.wrkbnch
 
A workbench for infrastructure concepts, each one with a working, interactive demo.
 
## Live
 
- Hub: [link once deployed]
- LRU Cache demo: [link once deployed]
 
## Structure
 
One shared React frontend, routing to each project's own page. Each project has
its own backend, deployed as its own service — not one monolith.
 
minj.wrkbnch/
├── frontend/          React + Vite — hub page + all project demo pages
└── CacheLRU/          Go backend for the LRU cache project
 
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
 
## Running locally
 
Backend
    cd CacheLRU
    go run main.go
 
Frontend
    cd frontend
    npm install
    npm run dev
 
Create frontend/.env with:
    VITE_CACHE_API_URL=http://localhost:8080
 
## Stack
 
- Frontend: React, Vite
- Backend: Go (net/http, no framework)
- Deploy: Render (one service per backend, one static site for the frontend)
