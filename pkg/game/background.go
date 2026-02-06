package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type parallaxBackground struct {
	textures         []rl.Texture2D
	speeds           []float32
	positions        []rl.Vector2
	initialPositions []rl.Vector2
	speedModifier    float32
	config           *Config
}

func (b *parallaxBackground) init() {
	b.textures = make([]rl.Texture2D, 0)
	b.speeds = make([]float32, 0)
	b.positions = make([]rl.Vector2, 0)
	b.initialPositions = make([]rl.Vector2, 0)
	b.speedModifier = 1.0
}

func (b *parallaxBackground) add(fileName string, scrollSpeed float32, initialPosition rl.Vector2) {
	img := rl.LoadImage(fileName)
	if img.Width != b.config.WindowWidth || img.Height != b.config.WindowHeight {
		rl.ImageResize(img, b.config.WindowWidth, b.config.WindowHeight)
	}
	b.textures = append(b.textures, rl.LoadTextureFromImage(img))
	b.speeds = append(b.speeds, scrollSpeed)
	b.positions = append(b.positions, initialPosition)
	b.initialPositions = append(b.initialPositions, initialPosition)
}

func (b *parallaxBackground) update() {
	for i := 0; i < len(b.textures); i++ {
		b.positions[i].X -= b.speeds[i] * b.speedModifier
		if b.positions[i].X <= -float32(b.textures[i].Width) {
			b.positions[i].X = b.initialPositions[i].X
		}
	}
}

func (b *parallaxBackground) draw() {
	for i, texture := range b.textures {
		rl.DrawTextureEx(texture, b.positions[i], 0, 1, rl.White)
		rl.DrawTexture(texture, int32(b.positions[i].X)+texture.Width, int32(b.positions[i].Y), rl.White)
	}
}

func (b *parallaxBackground) reset() {
	for i := range b.positions {
		b.positions[i] = b.initialPositions[i]
	}
	b.speedModifier = 1.0
}

func newParallaxBackground(config *Config) (b parallaxBackground) {
	b.init()
	b.config = config
	return b
}
