# CLAUDE.md

`go-cli-pokedex`: Go CLI (stdlib only) that looks up Pokémon on PokeAPI.
Usage: `go-cli-pokedex get <pokemon> [pokemon...]`

## Commands

- Build: `go build -o go-cli-pokedex .`
- Run: `go run . get pikachu arceus`
- Vet: `go vet ./...`
- Test: `go test ./...` (offline; uses `httptest`)

`go-cli-pokedex` in the repo root is a build artifact; keep it out of git.

## Conventions

- Everything is `package main` in the repo root. `main.go` handles args and output, `pokeapi.go` handles the API.
- PokeAPI weight is in hectograms; `FetchPokemon` converts to kg (÷10).
- Wrap errors with `%w`. Only `main` prints errors and exits. A 404 returns the `ErrNotFound` sentinel.
- With multiple names, a failure on one is reported to stderr and the rest still run. Exit code is 1 if any failed.
- Use the shared `client` (10s timeout), never `http.Get`. Escape user input with `url.PathEscape`.
- `main` just calls `run(args, stdout, stderr) int`; test logic through `run`.
- `baseURL` and `client` are package vars so tests can swap them; use the `fakeAPI` helper in tests, never the real API.
- Keep the raw `apiResponse` separate from the domain `Pokemon` type.
