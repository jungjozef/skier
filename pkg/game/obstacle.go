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
	gapY      float32
	gapHeight float32
	velocityX float32
	velocityY float32
	active    bool
	scored    bool // for near-miss tracking
}

func (o *obstacle) topRect() rl.Rectangle {
	return rl.NewRectangle(o.position.X, 0, o.size.X, o.gapY)
}

func (o *obstacle) bottomRect() rl.Rectangle {
	return rl.NewRectangle(o.position.X, o.gapY+o.gapHeight, o.size.X, 1080-(o.gapY+o.gapHeight))
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

	img := rl.LoadImage("assets/blizzardskull.png")
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
			// Spawn gate
			gapY := float32(rl.GetRandomValue(200, 800))
			gapH := float32(150)
			om.obstacles = append(om.obstacles, obstacle{
				position:  rl.NewVector2(1920+20, 0),
				size:      rl.NewVector2(40, 1080),
				otype:     obstacleGate,
				gapY:      gapY,
				gapHeight: gapH,
				active:    true,
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
		if o.otype == obstacleGate {
			if rl.CheckCollisionRecs(skierRect, o.topRect()) || rl.CheckCollisionRecs(skierRect, o.bottomRect()) {
				collision = true
			}
			// Near-miss check
			if !o.scored && o.position.X < skierPos.X {
				o.scored = true
				// Check if skier was within 80px vertically of the bars
				distToTop := skierPos.Y - o.gapY
				distToBottom := (o.gapY + o.gapHeight) - skierPos.Y
				if distToTop > 0 && distToBottom > 0 && (distToTop < 80 || distToBottom < 80) {
					nearMissBonus += 50
				}
			}
		} else {
			if rl.CheckCollisionRecs(skierRect, o.rect()) {
				collision = true
			}
			// Near-miss for airplane
			if !o.scored && o.position.X < skierPos.X-40 {
				o.scored = true
				dx := skierPos.X - (o.position.X + o.size.X/2)
				dy := skierPos.Y - (o.position.Y + o.size.Y/2)
				dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				if dist < 80 {
					nearMissBonus += 50
				}
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
			// Top bar
			rl.DrawRectangle(int32(o.position.X), 0, int32(o.size.X), int32(o.gapY), rl.NewColor(255, 50, 50, 200))
			// Bottom bar
			bottomY := o.gapY + o.gapHeight
			rl.DrawRectangle(int32(o.position.X), int32(bottomY), int32(o.size.X), int32(1080-bottomY), rl.NewColor(255, 50, 50, 200))
			// Gap indicators
			rl.DrawLineEx(
				rl.NewVector2(o.position.X, o.gapY),
				rl.NewVector2(o.position.X+o.size.X, o.gapY),
				3, rl.Yellow)
			rl.DrawLineEx(
				rl.NewVector2(o.position.X, o.gapY+o.gapHeight),
				rl.NewVector2(o.position.X+o.size.X, o.gapY+o.gapHeight),
				3, rl.Yellow)
		} else {
			rl.DrawTexture(om.skullTexture, int32(o.position.X), int32(o.position.Y), rl.White)
		}
	}
}
