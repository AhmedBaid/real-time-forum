# Real-Time Forum

A simple real-time forum application (Go backend + static frontend) that includes user registration, login, posting, comments, reactions, and a WebSocket-based real-time feed.

This repository contains:

- `backend/` - Go server, handlers, middleware, SQLite database file (`database/db.db`).
- `frontend/` - Static frontend (HTML/CSS/JS) located in `frontend/main.html` and `frontend/src/`.

## Project structure

A quick tree of the important files and folders:

```
.
├─ backend/
│  ├─ cmd/
│  │  └─ main.go
│  ├─ config/
│  │  ├─ config.go
│  │  └─ ResponseJSON.go
│  ├─ database/
│  │  ├─ db.db
│  │  └─ query.sql
│  ├─ handler/
│  │  ├─ createCommentHandler.go
│  │  ├─ CreatePostHandler.go
│  │  ├─ CurrentUserHandler.go
│  │  ├─ getCommentsHandler.go
│  │  ├─ getPosts.go
│  │  ├─ getUsers.go
│  │  ├─ homeHandler.go
│  │  ├─ LoginHandler.go
│  │  ├─ LogoutHandler.go
│  │  ├─ ReactionHandler.go
│  │  ├─ RegisterHandler.go
│  │  ├─ StaticController.go
│  │  └─ websocket.go
│  ├─ helpers/
│  │  ├─ FetchCategories.go
│  │  └─ SessionChecked.go
│  ├─ middleware/
│  │  ├─ auth.go
│  │  └─ rateLimit.go
│  └─ router/
│     └─ router.go
├─ frontend/
│  ├─ main.html
│  └─ src/
│     ├─ assets/
│     │  └─ icon.png
│     ├─ css/
│     │  ├─ global.css
│     │  ├─ home.css
│     │  ├─ login.css
│     │  ├─ notfound.css
│     │  ├─ register.css
│     │  └─ responsive.css
│     └─ js/
│        ├─ config.js
│        ├─ loadPage.js
│        ├─ main.js
│        ├─ helpers/
│        │  ├─ api.js
│        │  ├─ randerComments.js
│        │  ├─ renderPosts.js
│        │  ├─ showToast.js
│        │  ├─ sortUsers.js
│        │  └─ timeFormat.js
│        └─ views/
│           ├─ createComments.js
│           ├─ createPost.js
│           ├─ HandleLikes.js
│           ├─ HandleMessages.js
│           ├─ home.js
│           ├─ login.js
│           ├─ notfound.js
│           └─ register.js
├─ go.mod
├─ go.sum
└─ README.md
```

(If your workspace has additional files, they will appear in the tree above when you list files locally.)

## Quick start

### Prerequisites

- Go 1.18+ installed
- SQLite (optional, the repo already contains a `backend/database/db.db` file)

### Run the backend (development)

From the repository root run:

```zsh
go run ./backend/cmd/main.go
```

This starts the Go server. The server entrypoint is `backend/cmd/main.go`.

If you prefer to build a binary:

```zsh
go build -o bin/forum ./backend/cmd
./bin/forum
```


## Routes and handlers

Handlers are implemented in `backend/handler/` and routes are wired in `backend/router/router.go`.

Notable handlers (for reference):

- `RegisterHandler` — user registration
- `LoginHandler` / `LogoutHandler` — authentication
- `CreatePostHandler`, `getPosts.go` — posts
- `createCommentHandler`, `getCommentsHandler` — comments
- `ReactionHandler` — likes/reactions
- `websocket.go` — WebSocket for real-time updates
- `StaticController` — serves static files (if used)

For exact endpoints and middleware (auth, rate limiting), inspect `backend/router/router.go` and `backend/middleware/`.

## Development notes

- Middleware: authentication and rate limiting are in `backend/middleware/`.
- Helpers: session checks and category fetching are in `backend/helpers/`.
- Config: application configuration is in `backend/config/`.

## Testing

There are no automated tests included by default. To manually test:

1. Start the backend.
2. Open the frontend and register/login.
3. Create posts and comments, check real-time behavior.

## Contributing

Contributions are welcome. Open issues or pull requests with a clear description of the changes.

## License

This project doesn't include a license file. Add a `LICENSE` file if you intend to open-source it.

---

If you want, I can:

- Add a more detailed list of exact HTTP routes and sample requests by reading `backend/router/router.go`.
- Add a small Makefile or Dockerfile for easy runs.
- Add a LICENSE file (MIT/Apache/etc.).

Tell me which follow-up you'd like and I'll implement it.
