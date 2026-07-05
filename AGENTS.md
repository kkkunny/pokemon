# AGENTS.md

## Run

```bash
go run .
```

No tests, no CI, no lint config, no Makefile.

## Architecture

Stack-based subsystem chain managed in `src/game.go`. Only the **last** subsystem receives `OnAction`/`OnUpdate`/`OnDraw`. Subsystems delegate to `system.Next()` when inactive.

- `src/system/` — `System` interface + `SystemManager` cursor
- `src/system/world/` — overworld map / sprite movement
- `src/system/battle/` — turn-based battle UI
- `src/system/dialogue/` — text box overlay

## Animation timing (critical gotcha)

**`Player.Update()` must be called in `OnUpdate` (TPS rate), never in `OnDraw` (FPS rate).**

- `animation.NewPlayer(fps)` precomputes frame counters based on `fps` (passed as `config.DefaultTPS` = 60). Each counter = `int(frameTime / (1s/fps))` ticks per frame.
- `Player.Update()` advances by exactly 1 tick per call.
- Ebitengine calls `OnUpdate` at TPS rate and `OnDraw` at display refresh rate. These can differ (e.g. TPS=60, FPS=100).
- Putting `.Update()` in `OnDraw` makes animation speed scale with FPS, desyncing from the precomputed counters.

GIF delay parsing: `time.Second / 100 * time.Duration(g.Delay[i])`. GIF delays are in hundredths of a second, converted to `time.Duration` nanoseconds.

## Key config values

- `config.DefaultTPS = 60` — controls animation `NewPlayer` tick rate
- `config.Scale = 2`, `config.TileSize = 16`

## Ebitengine v2.8.8 quirks

- No `SetMaxFPS`. `SetFPSMode` exists but is deprecated. Use `SetVsyncEnabled` instead.
- `ebitenutil.DebugPrint` shows FPS/TPS overlay (hardcoded in `main.go`).

## Data paths

Asset root: `data/`. Paths defined in `src/consts/path.go`:

| Var | Path |
|---|---|
| `PokemonDefinePath` | `data/pokemons/` |
| `GFXBattleSitesPath` | `data/gfx/battle_sites/` |
| `GFXMapPath` | `data/gfx/map/` |
| `MapsPath` | `data/world/maps/` |

Pokemon species sprites: `data/pokemons/<id>/front.gif` + `back.gif`. Loaded via `util.FindFileAndThenParse` (matches filename ignoring extension).

## Code conventions

- Enum pattern: `var FooEnum = enum.New[struct{ ... }]()` from `github.com/tnnmigga/enum`
- Drawing: fluent `option.DrawXxxOptions` via `draw.PrepareDrawXxx(drawer, ...).Move(...).Draw()`
- No comments in production code — only user-facing comments in Chinese
- Collection utils: `github.com/kkkunny/stl` (`stlslices`, `stlmaps`, `stlval`)

## Known build break

`src/script/script.go` references `config.ScriptsPath` which doesn't exist (moved to `consts.ScriptsPath`). Pre-existing — don't fix unless asked.
