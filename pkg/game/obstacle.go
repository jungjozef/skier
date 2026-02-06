package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	obstacleGate     = 0
	obstacleAirplane = 1
)

type obstacle struct {
	position  rl.Vector2
	size      rl.Vector2
	otype     int
	velocityX float32
	velocityY float32
	active    bool
	scored    bool // for near-miss tracking
}

func (o *obstacle) rect() rl.Rectangle {
	return rl.NewRectangle(o.position.X, o.position.Y, o.size.X, o.size.Y)
}

type obstacleManager struct {
	obstacles    []obstacle
	spawnTimer   float32
	spawnRate    float32
	elapsed      float32
	skullTexture rl.Texture2D
}

func (om *obstacleManager) init() {
	om.obstacles = make([]obstacle, 0)
	om.spawnTimer = 0
	om.spawnRate = 300 // 5 seconds at 60fps
	om.elapsed = 0

	img := rl.LoadImage("../assets/blizzardskull.png")
	rl.ImageResize(img, 80, 80)
	om.skullTexture = rl.LoadTextureFromImage(img)
}

func (om *obstacleManager) reset() {
	om.obstacles = make([]obstacle, 0)
	om.spawnTimer = 0
	om.spawnRate = 300
	om.elapsed = 0
}

func (om *obstacleManager) update(scrollSpeed float32, skierPos rl.Vector2) (bool, float32) {
	om.elapsed++
	om.spawnTimer++

	// Difficulty scaling
	switch {
	case om.elapsed > 7200: // 120s+
		om.spawnRate = 60
	case om.elapsed > 3600: // 60-120s
		om.spawnRate = 120
	case om.elapsed > 1800: // 30-60s
		om.spawnRate = 200
	default: // 0-30s
		om.spawnRate = 300
	}

	// Spawn
	if om.spawnTimer >= om.spawnRate {
		om.spawnTimer = 0
		if om.elapsed > 1800 && rl.GetRandomValue(0, 2) == 0 {
			// Spawn airplane/skull
			yPos := float32(rl.GetRandomValue(100, 800))
			om.obstacles = append(om.obstacles, obstacle{
				position:  rl.NewVector2(1920+40, yPos),
				size:      rl.NewVector2(80, 80),
				otype:     obstacleAirplane,
				velocityX: -float32(rl.GetRandomValue(3, 8)),
				velocityY: float32(rl.GetRandomValue(-2, 2)),
				active:    true,
			})
		} else {
			// Spawn gate — a single bar of limited height
			barH := float32(rl.GetRandomValue(250, 400))
			barY := float32(rl.GetRandomValue(100, int32(1080-barH-100)))
			om.obstacles = append(om.obstacles, obstacle{
				position: rl.NewVector2(1920+20, barY),
				size:     rl.NewVector2(40, barH),
				otype:    obstacleGate,
				active:   true,
			})
		}
	}

	// Update and check collisions
	collision := false
	nearMissBonus := float32(0)
	skierRect := rl.NewRectangle(skierPos.X-24, skierPos.Y-24, 48, 48)

	active := make([]obstacle, 0, len(om.obstacles))
	for i := range om.obstacles {
		o := &om.obstacles[i]
		if o.otype == obstacleGate {
			o.position.X -= scrollSpeed // gates scroll with the world
		} else {
			// Airplanes have their own independent velocity
			o.position.X += o.velocityX
			o.position.Y += o.velocityY
		}

		// Check collision
		if rl.CheckCollisionRecs(skierRect, o.rect()) {
			collision = true
		}
		// Near-miss check
		if !o.scored && o.position.X < skierPos.X-o.size.X/2 {
			o.scored = true
			dx := skierPos.X - (o.position.X + o.size.X/2)
			dy := skierPos.Y - (o.position.Y + o.size.Y/2)
			dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			if dist < 100 {
				nearMissBonus += 50
			}
		}

		// Keep active if on screen
		if o.position.X > -100 && o.position.X < 2100 && o.position.Y > -200 && o.position.Y < 1300 {
			active = append(active, *o)
		}
	}
	om.obstacles = active

	return collision, nearMissBonus
}

func (om *obstacleManager) draw() {
	for _, o := range om.obstacles {
		if o.otype == obstacleGate {
			rl.DrawRectangle(int32(o.position.X), int32(o.position.Y), int32(o.size.X), int32(o.size.Y), rl.NewColor(255, 50, 50, 200))
			// Top and bottom edge highlights
			rl.DrawLineEx(
				rl.NewVector2(o.position.X, o.position.Y),
				rl.NewVector2(o.position.X+o.size.X, o.position.Y),
				3, rl.Yellow)
			rl.DrawLineEx(
				rl.NewVector2(o.position.X, o.position.Y+o.size.Y),
				rl.NewVector2(o.position.X+o.size.X, o.position.Y+o.size.Y),
				3, rl.Yellow)
		} else {
			rl.DrawTexture(om.skullTexture, int32(o.position.X), int32(o.position.Y), rl.White)
		}
	}
}
