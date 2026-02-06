package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game state enum
const (
	statePlaying  = 0
	stateGameOver = 1
)

// Game struct centralizing all state
type Game struct {
	config    *Config
	state     int
	mountain  mountain
	skier     skier
	bkg       parallaxBackground
	bkg2      parallaxBackground
	bkgGrad   rl.Texture2D
	sun       rl.Texture2D
	msx       rl.Music
	obstacles obstacleManager

	scrollSpeed float32
	score       float32
	highScore   int
}

func NewGame(config *Config) *Game {
	g := &Game{
		config:      config,
		state:       statePlaying,
		scrollSpeed: baseScrollSpeed,
		score:       0,
	}

	g.highScore = loadHighScore()

	g.msx = rl.LoadMusicStream("assets/music1.mp3")
	g.mountain = newMountain(config)
	g.skier = newSkier(&g.mountain)

	g.bkg = newParallaxBackground(config)
	g.bkg.add("../assets/landscape_0004_5_clouds.png", 0.5, rl.NewVector2(0, -100))
	g.bkg.add("../assets/landscape_0003_4_mountain.png", 2, rl.NewVector2(0, 0))
	g.bkg.add("../assets/landscape_0002_3_trees.png", 3, rl.NewVector2(0, 0))
	g.bkg.add("../assets/landscape_0001_2_trees.png", 4, rl.NewVector2(0, 0))
	g.bkg.add("../assets/landscape_0000_1_trees.png", 5, rl.NewVector2(0, 0))

	g.bkg2 = newParallaxBackground(config)
	g.bkg2.add("../assets/landscape_0001_2_trees_green.png", 7, rl.NewVector2(0, 350))

	g.bkgGrad = rl.LoadTextureFromImage(
		rl.GenImageGradientRadial(int(config.WindowWidth), int(0.65*float32(config.WindowHeight)), 12, rl.SkyBlue, rl.Beige))
	g.sun = rl.LoadTexture("../assets/sun.png")

	g.obstacles.init()

	return g
}

func (g *Game) reset() {
	g.state = statePlaying
	g.scrollSpeed = baseScrollSpeed
	g.score = 0
	g.mountain.init()
	g.skier.reset()
	g.skier.mountain = &g.mountain
	g.bkg.reset()
	g.bkg2.reset()
	g.obstacles.reset()
}

func (g *Game) Update() {
	if g.state == stateGameOver {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			g.reset()
		}
		return
	}

	// Music
	if !rl.IsMusicStreamPlaying(g.msx) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			rl.PlayMusicStream(g.msx)
		}
	}
	rl.UpdateMusicStream(g.msx)

	g.mountain.update(g.scrollSpeed)
	g.scrollSpeed = g.skier.update(g.scrollSpeed)

	// Parallax speed modifier
	speedRatio := g.scrollSpeed / baseScrollSpeed
	g.bkg.speedModifier = speedRatio
	g.bkg2.speedModifier = speedRatio

	g.bkg.update()
	g.bkg2.update()

	// Obstacles
	collision, nearMissBonus := g.obstacles.update(g.scrollSpeed, g.skier.position)
	g.score += nearMissBonus

	// Scoring
	if g.skier.onSlope {
		g.score += speedRatio
	} else {
		g.score += 2 * speedRatio // Jump bonus
	}

	// Game over conditions
	if g.skier.position.Y > float32(g.config.WindowHeight)+64 {
		g.triggerGameOver()
	}
	if collision {
		g.triggerGameOver()
	}
}

func (g *Game) triggerGameOver() {
	g.state = stateGameOver
	finalScore := int(g.score)
	if finalScore > g.highScore {
		g.highScore = finalScore
		saveHighScore(g.highScore)
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.DrawTexture(g.bkgGrad, 0, 0, rl.White)
	rl.DrawTexture(g.sun, g.config.WindowWidth/2-128, g.config.WindowHeight/2, rl.White)

	g.bkg.draw()
	g.mountain.draw()
	g.obstacles.draw()
	g.skier.draw()
	rl.DrawCircleV(rl.GetMousePosition(), 10, rl.RayWhite)
	g.bkg2.draw()

	// HUD
	scoreText := fmt.Sprintf("Score: %d", int(g.score))
	scoreWidth := rl.MeasureText(scoreText, 36)
	rl.DrawText(scoreText, g.config.WindowWidth-scoreWidth-20, 20, 36, rl.Black)

	highText := fmt.Sprintf("Best: %d", g.highScore)
	highWidth := rl.MeasureText(highText, 28)
	rl.DrawText(highText, g.config.WindowWidth-highWidth-20, 62, 28, rl.DarkGray)

	speedText := fmt.Sprintf("Speed: %.1f", g.scrollSpeed)
	rl.DrawText(speedText, 20, 20, 24, rl.DarkGray)

	// Game over overlay
	if g.state == stateGameOver {
		rl.DrawRectangle(0, 0, g.config.WindowWidth, g.config.WindowHeight, rl.NewColor(0, 0, 0, 150))

		title := "GAME OVER"
		titleWidth := rl.MeasureText(title, 80)
		rl.DrawText(title, (g.config.WindowWidth-titleWidth)/2, g.config.WindowHeight/2-100, 80, rl.White)

		finalText := fmt.Sprintf("Score: %d", int(g.score))
		finalWidth := rl.MeasureText(finalText, 48)
		rl.DrawText(finalText, (g.config.WindowWidth-finalWidth)/2, g.config.WindowHeight/2, 48, rl.Yellow)

		bestText := fmt.Sprintf("Best: %d", g.highScore)
		bestWidth := rl.MeasureText(bestText, 36)
		rl.DrawText(bestText, (g.config.WindowWidth-bestWidth)/2, g.config.WindowHeight/2+60, 36, rl.LightGray)

		restartText := "Click to restart"
		restartWidth := rl.MeasureText(restartText, 30)
		rl.DrawText(restartText, (g.config.WindowWidth-restartWidth)/2, g.config.WindowHeight/2+120, 30, rl.RayWhite)
	}

	rl.EndDrawing()
}
