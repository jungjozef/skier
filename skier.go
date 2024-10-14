package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	gravity = 9.81
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
		//rl.DrawPixel(int32(s.points[0].X), 10, rl.Red)
		rl.DrawCircle(int32(s.points[0].X), 1080, 6, rl.Red)
		for i := 1; i < len(s.points); i++ {
			for j := 0; j <= 2; j++ {
				v0 := rl.NewVector2(s.points[i-1].X, s.points[i-1].Y+float32(j)*30)
				v1 := rl.NewVector2(v0.X, float32(s.config.windowHeight))
				v2 := rl.NewVector2(s.points[i].X, s.points[i].Y+float32(j)*30)
				v4 := rl.NewVector2(v2.X, float32(s.config.windowHeight))
				color := rl.LightGray
				if j == 1 {
					color = rl.Gray
				} else if j == 2 {
					color = rl.DarkGray
				}
				rl.DrawTriangleStrip([]rl.Vector2{v0, v1, v2, v4}, color)
			}

			rl.DrawLineEx(s.points[i-1], s.points[i], 5, rl.RayWhite)
			//rl.DrawTriangleFan(s.points, rl.White)

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

func (s *Slope) heightAt(x float32) float32 {
	if s.lastPoint().X < x {
		return -1
	}

	for i := len(s.points) - 1; i >= 1; i-- {
		p1 := s.points[i-1]
		p2 := s.points[i]
		if p1.X == x {
			return p1.Y
		}
		if p2.X == x {
			return p2.Y
		}
		if p1.X <= x && p2.X >= x {
			return ((p2.Y-p1.Y)/(p2.X-p1.X))*(x-p1.X) + p1.Y
		}
	}

	return -1

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

func remove[T any](slice []T, i int) []T {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func (m *Mountain) update() {
	p := rl.GetMousePosition()
	for i := 0; i < len(m.slopes); i++ {
		m.slopes[i].scroll(12)
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
	var activeSlopes = make([]Slope, 0)
	for _, s := range m.slopes {
		if s.lastPoint().X >= 0 {
			activeSlopes = append(activeSlopes, s)
		}
	}
	m.slopes = activeSlopes

}

func (m *Mountain) heightAt(x float32) (float32, *Slope) {
	for i := len(m.slopes) - 1; i >= 0; i-- {
		height := m.slopes[i].heightAt(x)
		if height != -1 {
			return height, &m.slopes[i]
		}
	}
	return -1, nil
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
	texture  rl.Texture2D
	velocity float32
	mountain *Mountain
	onSlope  bool
}

func (s *Skier) draw() {
	rl.DrawTextureEx(s.texture, rl.NewVector2(s.position.X-32, s.position.Y-32), 0, 1, rl.White)
}

func (s *Skier) init() {
	s.position = rl.NewVector2(400, 20)
	img := rl.GenImageGradientRadial(64, 64, 0.5, rl.Yellow, rl.Blank)
	s.texture = rl.LoadTextureFromImage(img)
	s.velocity = 20
	s.onSlope = false
}

func (s *Skier) update() {
	if s.position.Y > 1080 {
		s.position.Y = 0
	}
	h, _ := s.mountain.heightAt(s.position.X)
	if h == -1 {
		s.position.Y += s.velocity
		s.onSlope = false
		return
	}
	if s.onSlope {
		s.position.Y = h
		return
	}
	nextPosY := s.position.Y + s.velocity
	if s.position.Y < h && nextPosY > h {
		s.position.Y = h
		s.onSlope = true
	} else {
		s.position.Y = nextPosY
		s.onSlope = false
	}

}

func NewSkier(mountain *Mountain) (s Skier) {
	s.init()
	s.mountain = mountain
	return s
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
	skier := NewSkier(&mountain)
	bkg := NewParallaxBackground(&cfg)
	bkg.add("assets/landscape_0004_5_clouds.png", 0.5, rl.NewVector2(0, -100))
	bkg.add("assets/landscape_0003_4_mountain.png", 2, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0002_3_trees.png", 3, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0001_2_trees.png", 4, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0000_1_trees.png", 5, rl.NewVector2(0, 0))

	bkg2 := NewParallaxBackground(&cfg)
	bkg2.add("assets/landscape_0001_2_trees_green.png", 7, rl.NewVector2(0, 350))
	bkgGrad := rl.LoadTextureFromImage(
		rl.GenImageGradientRadial(int(cfg.windowWidth), int(0.65*float32(cfg.windowHeight)), 12, rl.SkyBlue, rl.Beige))
	sun := rl.LoadTexture("assets/sun.png")
	for !rl.WindowShouldClose() {
		//rl.UpdateMusicStream(msx)
		if !rl.IsMusicStreamPlaying(msx) {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				rl.PlayMusicStream(msx)
			}
		}
		mountain.update()
		skier.update()
		bkg.update()
		bkg2.update()

		rl.BeginDrawing()
		//rl.ClearBackground(rl.NewColor(235, 239, 242, 255))
		rl.DrawTexture(bkgGrad, 0, 0, rl.White)
		rl.DrawTexture(sun, cfg.windowWidth/2-128, cfg.windowHeight/2, rl.White)

		bkg.draw()
		mountain.draw()
		skier.draw()
		rl.DrawCircleV(rl.GetMousePosition(), 10, rl.RayWhite)
		bkg2.draw()
		h, _ := mountain.heightAt(skier.position.X)
		rl.DrawText(fmt.Sprintf("Slope: %f", h), 10, 10, 30, rl.Black)
		rl.DrawText(fmt.Sprintf("Skier: %f", skier.position.Y), 10, 50, 30, rl.Black)

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
