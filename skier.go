package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Config related data structure
type Config struct {
	windowWidth  int32
	windowHeight int32
	windowTitle  string
}

// Slope related data structure and functions
type Slope struct {
	points     []rl.Vector2
	pointCount int32
	active     bool
	config     *Config
}

func (s *Slope) init() {
	s.pointCount = 0
	s.points = make([]rl.Vector2, 0)
	s.active = false
}

func (s *Slope) add() {
	s.pointCount++
	s.points = append(s.points, rl.GetMousePosition())
}

func (s *Slope) draw() {
	if !s.active && s.lastPoint().X < -10 {
		return
	}
	if s.pointCount > 0 {
		for i := 1; i < len(s.points); i++ {

			pointA := s.points[i-1]
			pointB := s.points[i]
			pointNum := pointB.X - pointA.X
			diffX := pointB.X - pointA.X
			diffY := pointB.Y - pointA.Y
			intervalX := diffX / (pointNum + 1)
			intervalY := diffY / (pointNum + 1)
			rl.DrawLineEx(s.points[i-1], s.points[i], 5, rl.RayWhite)
			for j := 0; j < int(pointNum); j++ {
				x1, y1 := pointA.X+intervalX*float32(j), pointA.Y+intervalY*float32(j)
				x2, y2 := pointA.X+intervalX*float32(j), float32(s.config.windowHeight)+pointA.Y+intervalY*float32(j)
				rl.DrawLineEx(rl.NewVector2(x1, y1+2.5), rl.NewVector2(x2, y2), 4, rl.LightGray)
				rl.DrawLineEx(rl.NewVector2(x1, y1+30), rl.NewVector2(x2, y2), 4, rl.Gray)
				rl.DrawLineEx(rl.NewVector2(x1, y1+60), rl.NewVector2(x2, y2), 4, rl.DarkGray)
			}

		}

	}
}

func (s *Slope) scroll(speed float32) {
	for i := 0; i < int(s.pointCount); i++ {
		s.points[i].X -= speed
	}
}

func (s *Slope) lastPoint() rl.Vector2 {
	return s.points[len(s.points)-1]
}

func NewSlope(config *Config) (slope Slope) {
	slope.init()
	slope.config = config
	return slope
}

//Slope end

// Mountain related data structure and functions
type Mountain struct {
	slopes []Slope
	config *Config
}

func (m *Mountain) init() {
	m.slopes = make([]Slope, 0)
}

func (m *Mountain) draw() {
	for i := 0; i < len(m.slopes); i++ {
		m.slopes[i].draw()
	}
}

func (m *Mountain) update() {
	p := rl.GetMousePosition()
	for i := 0; i < len(m.slopes); i++ {
		m.slopes[i].scroll(6)
	}

	if (p.X >= 0 && p.X <= float32(m.config.windowWidth)) && (p.Y >= 0 && p.Y <= float32(m.config.windowHeight)) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			m.slopes = append(m.slopes, NewSlope(m.config))
			m.slopes[len(m.slopes)-1].add()
			m.slopes[len(m.slopes)-1].active = true
		}
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if p.X > m.slopes[len(m.slopes)-1].lastPoint().X {
				m.slopes[len(m.slopes)-1].add()
			}
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			m.slopes[len(m.slopes)-1].active = false
		}
	}
}

func NewMountain(config *Config) (mountain Mountain) {
	mountain.init()
	mountain.config = config
	return mountain
}

// Mountain end

// ParallaxBackground data structure and functions
type ParallaxBackground struct {
	textures         []rl.Texture2D
	speeds           []float32
	positions        []rl.Vector2
	initialPositions []rl.Vector2
	speedModifier    float32
	config           *Config
}

func (b *ParallaxBackground) init() {
	b.textures = make([]rl.Texture2D, 0)
	b.speeds = make([]float32, 0)
	b.positions = make([]rl.Vector2, 0)
	b.initialPositions = make([]rl.Vector2, 0)
	b.speedModifier = 1.0
}

func (b *ParallaxBackground) add(fileName string, scrollSpeed float32, initialPosition rl.Vector2) {
	img := rl.LoadImage(fileName)
	if img.Width != b.config.windowWidth || img.Height != b.config.windowHeight {
		rl.ImageResize(img, b.config.windowWidth, b.config.windowHeight)
	}

	b.textures = append(b.textures, rl.LoadTextureFromImage(img))
	b.speeds = append(b.speeds, scrollSpeed)
	b.positions = append(b.positions, initialPosition)
	b.initialPositions = append(b.initialPositions, initialPosition)
}

func (b *ParallaxBackground) update() {
	for i := 0; i < len(b.textures); i++ {
		b.positions[i].X -= b.speeds[i] * b.speedModifier
		if b.positions[i].X <= -float32(b.textures[i].Width) {
			b.positions[i].X = b.initialPositions[i].X
		}
	}
}

func (b *ParallaxBackground) draw() {
	for i, texture := range b.textures {
		rl.DrawTextureEx(texture, b.positions[i], 0, 1, rl.White)
		rl.DrawTexture(texture, int32(b.positions[i].X)+texture.Width, int32(b.positions[i].Y), rl.White)
	}
}

func NewParallaxBackground(config *Config) (b ParallaxBackground) {
	b.init()
	b.config = config
	return b
}

// ParallaxBackground end

// Skier related data structures and functions
type Skier struct {
	position rl.Vector2
}

// Skier end

func main() {
	cfg := Config{
		1920, 1080, "Skier",
	}
	rl.InitWindow(cfg.windowWidth, cfg.windowHeight, cfg.windowTitle)
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)
	rl.HideCursor()

	msx := rl.LoadMusicStream("assets/music1.mp3")
	mountain := NewMountain(&cfg)
	bkg := NewParallaxBackground(&cfg)
	bkg.add("assets/landscape_0004_5_clouds.png", 0.5, rl.NewVector2(0, -100))
	bkg.add("assets/landscape_0003_4_mountain.png", 2, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0002_3_trees.png", 3, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0001_2_trees.png", 4, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0000_1_trees.png", 5, rl.NewVector2(0, 0))

	bkg2 := NewParallaxBackground(&cfg)
	bkg2.add("assets/landscape_0001_2_trees_green.png", 7, rl.NewVector2(0, 350))
	bkgGrad := rl.LoadTextureFromImage(
		rl.GenImageGradientV(int(cfg.windowWidth), int(0.65*float32(cfg.windowHeight)), rl.SkyBlue, rl.Beige))
	for !rl.WindowShouldClose() {
		rl.UpdateMusicStream(msx)
		if !rl.IsMusicStreamPlaying(msx) {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				rl.PlayMusicStream(msx)
			}
		}
		mountain.update()
		bkg.update()
		bkg2.update()

		rl.BeginDrawing()
		//rl.ClearBackground(rl.NewColor(235, 239, 242, 255))
		rl.DrawTexture(bkgGrad, 0, 0, rl.White)

		bkg.draw()
		mountain.draw()
		rl.DrawCircleV(rl.GetMousePosition(), 10, rl.RayWhite)
		bkg2.draw()

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
