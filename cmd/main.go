package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"skier/pkg/game"
)

func main() {
	cfg := game.Config{
		WindowWidth:  1920,
		WindowHeight: 1080,
		WindowTitle:  "Skier",
	}
	rl.InitWindow(cfg.WindowWidth, cfg.WindowHeight, cfg.WindowTitle)
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)
	rl.HideCursor()

	g := game.NewGame(&cfg)

	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
	rl.CloseWindow()
}
