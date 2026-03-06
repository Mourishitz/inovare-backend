# Copilot Instructions

## Build & Run

```sh
go build ./...          # compile
go run main.go          # run server
air                     # hot-reload dev server (uses .air.toml)
make migrate            # apply DB migrations via Atlas
```

Migrations are managed with **Atlas + GORM provider** (not `AutoMigrate`). To generate a new migration after changing models:
```sh
atlas migrate diff --env gorm
```

## Architecture

Layered: `routes → controllers → services → repositories → models`

- **routes/** — instantiate services/controllers and register Gin route groups
- **controllers/** — bind/validate requests, check auth, delegate to services
- **services/** — business logic; orchestrate multiple repositories
- **repositories/** — all GORM DB access; return sentinel errors from `utils/errors.go`
- **models/** — GORM model structs; enums as `int16` constants in `models/enums/`
- **requests/** — input structs with `binding:` tags validated via `utils.BindAndValidate()`
- **utils/** — shared: JWT, bcrypt, validator helper, sentinel errors
- **config/** — reads `.env` via `godotenv`; accessed globally via `config.GetConfig()`
- **database/** — exposes `database.DB` (`*gorm.DB`) after `database.Connect()`

## Key Conventions

### Interfaces & Constructors
Every service and repository is defined as an interface in its own file. The private struct implements it. Always use `NewXxxService()` / `NewXxxRepository()` which return the interface:
```go
type ShowerRepository interface { ... }
type showerRepository struct { db *gorm.DB }
func NewShowerRepository() ShowerRepository { return &showerRepository{db: database.DB} }
```

### Request Validation
Use `utils.BindAndValidate(ctx, &req)` in all controller handlers — it binds JSON, validates, and writes a 400 response on failure. Return immediately if it returns `false`.

### Error Handling
All domain errors are sentinel `var Err... = errors.New(...)` values in `utils/errors.go`. Repositories wrap `gorm.ErrRecordNotFound` → domain error. Controllers switch on `errors.Is(err, utils.ErrXxx)` to pick the right HTTP status.

### Partial Updates
Update request structs use pointer fields (`*string`, `*uint`, etc.) so only non-nil fields are written via `db.Model(&x).Updates(map[string]interface{}{...})`.

### Role-Based Access
`User.Role` is `int16`. Roles: `1 = RoleUser`, `2 = RoleAdmin` (see `models/enums/roles.go`). Role checks (`user.Role >= 2`) are done inline in controller methods, not in middleware.

### Auth Middleware
- `middlewares.Authenticate()` — required auth; sets `"userID"` (int) and `"username"` (string) in Gin context
- `middlewares.OptionalAuthMiddleware()` — same but allows unauthenticated requests through

JWT uses `JWT_SECRET_KEY` env var (read directly via `os.Getenv`, not from `config.GetConfig()`).

### Custom JSON Marshaling
`models.Catalog` uses a custom `MarshalJSON` to expose `url` as a full frontend URL and `package` as a human-readable string. The raw DB fields are tagged `json:"-"`.

### Enums
Preferences/catalog fields (style, colors, sizes, etc.) are stored as `int16`. Enum constants and name maps live in `models/enums/`. The `Preferences` model uses `models.Int16Array` for slice fields.

## Environment Variables

Copy `.env.example` to `.env`. Key vars:
| Variable | Purpose |
|---|---|
| `SERVER_PORT` | Gin listen port |
| `GIN_MODE` | `release` or `debug` |
| `DB_HOST/USER/PASSWORD/NAME` | PostgreSQL connection |
| `JWT_SECRET_KEY` | JWT signing secret |
| `FRONTEND_URL` | Used to build catalog public URLs |
