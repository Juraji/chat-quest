# AGENTS.md

## Architecture

- **Frontend**: Angular 20 (root `src/`, runs on port 4200, served to port 8080)
- **Backend**: Go (in `backend/`, serves API on `/api`, serves UI from `./browser`)
- **Database**: SQLite with golang-migrate migrations in `backend/core/database/migrations/`

## Developer Commands

```bash
npm start         # Start Angular dev server
npm run build     # Build Angular UI
./build.sh        # Build full release (UI + Go binaries + zip packages)
cd backend && go run .   # Run Go backend (requires UI built or CHAT_QUEST_UI_ROOT set)
```

## Environment Variables (Go backend)

| Variable                        | Default          | Description            |
|---------------------------------|------------------|------------------------|
| `CHAT_QUEST_DATA_DIR`           | `./data`         | Data storage directory |
| `CHAT_QUEST_UI_ROOT`            | `./browser`      | Angular build output   |
| `CHAT_QUEST_DEBUG`              | `false`          | Enable debug logging   |
| `CHAT_QUEST_APPLICATION_PORT`   | `8080`           | HTTP server port       |
| `CHAT_QUEST_CORS_ALLOW_ORIGINS` | `localhost:8080` | Allowed CORS origins   |

## Key Entry Points

- Backend entry: `backend/main.go` (lines 1-85)
- Template system: `backend/core/util/templates.go`
- Default instructions: `backend/model/instructions/default_templates.go`
- Providers: `backend/core/providers/` (OpenAI-compatible API clients)

## Testing

- We don't care about tests.

## Important Patterns

- Go templates use `{{ .Variable }}` syntax for dynamic content
- Character data: `Character`, `SparseCharacter`, `WorldVars` types in templates
- Multi-char chats: Use `<characterid>[id]</characterid>` XML tags in responses
- Connection profiles support OpenAI-compatible APIs (OpenRouter, LM Studio, Ollama, Koboldcpp)
