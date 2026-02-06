# Skier

A Solipskier-inspired slope-drawing ski game built with Go and raylib.

## Tech Stack

- **Language:** Go 1.23
- **Graphics/Audio:** raylib via [github.com/gen2brain/raylib-go/raylib](https://github.com/gen2brain/raylib-go)
- **Structure:** `cmd/` entry point + `pkg/game/` package, one type per file

## Build & Run

```sh
go build -o skier ./cmd/main.go
./skier
```

## Project Structure

```
cmd/
  main.go             # Entry point (package main), imports pkg/game
pkg/game/             # All game logic (package game)
  config.go           # Config struct (exported)
  game.go             # Game struct, NewGame, Update, Draw, reset
  skier.go            # skier struct + methods, physics constants
  slope.go            # slope struct + methods
  mountain.go         # mountain struct + methods
  background.go       # parallaxBackground struct + methods
  obstacle.go         # obstacle, obstacleManager structs + methods
  score.go            # loadHighScore, saveHighScore
assets/               # Images and audio
  music1.mp3
  sun.png
  blizzardskull.png   # Dynamic obstacle sprite
  landscape_*.png     # Parallax background layers
highscore.dat         # Persistent high score (auto-generated)
```

## Architecture

All state lives in the `Game` struct which owns:

- **Mountain** — collection of `Slope`s drawn by mouse input. Each slope has points, supports `heightAt(x)` and `angleAt(x)` interpolation.
- **Skier** — fixed at X=400, physics-driven Y position. On-slope: snaps to slope height, accelerates/decelerates based on slope angle. Airborne: gravity, terminal velocity, rotation.
- **ParallaxBackground** — layered scrolling backgrounds with speed modifier tied to game scroll speed.
- **ObstacleManager** — spawns gate obstacles (vertical bars with gap) and dynamic skull obstacles (`blizzardskull.png`). Difficulty scales over time.

### Game Loop

`main()` -> `game.Update()` + `game.Draw()` each frame at 60 FPS.

### Game States

- `StatePlaying` — active gameplay, mouse draws slopes
- `StateGameOver` — overlay with score, click to restart via `game.reset()`

## Physics Constants

| Constant | Value | Purpose |
|---|---|---|
| `baseScrollSpeed` | 12.0 | Default world scroll speed |
| `maxScrollSpeed` | 30.0 | Speed cap |
| `minScrollSpeed` | 4.0 | Speed floor |
| `gravityAccel` | 0.6 | Freefall acceleration per frame |
| `slopeAccelFactor` | 0.15 | Slope angle to speed change multiplier |
| `friction` | 0.02 | Per-frame speed reduction on slope |
| `launchFactor` | 0.8 | Slope speed to vertical launch velocity |
| `terminalVelocity` | 25.0 | Max falling speed |

## Key Patterns & Gotchas

- **Slice bounds on `Mountain.update()`:** Always guard `m.slopes[len(m.slopes)-1]` with `len(m.slopes) > 0`. `IsMouseButtonDown` / `IsMouseButtonReleased` can fire when the slopes slice is empty (after reset or all slopes scrolled off-screen).
- **Slope launch angle:** When the skier leaves a slope end, `angleAt()` may return 0 because the slope no longer covers skierX. The code uses `s.rotation` (last frame's slope angle) instead.
- **Landing approach angle:** Must capture `approachAngle` from `velocityY` *before* zeroing it on landing.
- **Obstacle movement:** Gates scroll with the world (`position.X -= scrollSpeed`). Airplanes use their own independent velocity, not world scroll.
- **High score:** Plain text integer in `highscore.dat`, read/written with `os.ReadFile`/`os.WriteFile`.

## Scoring

- On slope: `+1 * (scrollSpeed / baseScrollSpeed)` per frame
- Airborne: `+2 * (scrollSpeed / baseScrollSpeed)` per frame
- Near-miss (obstacle within 80px without collision): `+50` bonus

## Obstacle Difficulty Scaling

| Time | Spawn Interval | Types |
|---|---|---|
| 0-30s | 5s (300 frames) | Gates only |
| 30-60s | 3.3s (200 frames) | Gates + skulls |
| 60-120s | 2s (120 frames) | Mixed, faster |
| 120s+ | 1s (60 frames) | Dense |
