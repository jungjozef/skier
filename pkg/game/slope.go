package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type slope struct {
	points     []rl.Vector2
	pointCount int32
	active     bool
	config     *Config
}

func (s *slope) init() {
	s.pointCount = 0
	s.points = make([]rl.Vector2, 0)
	s.active = false
}

func (s *slope) add() {
	s.pointCount++
	s.points = append(s.points, rl.GetMousePosition())
}

func (s *slope) draw() {
	if !s.active && s.lastPoint().X < -10 {
		return
	}
	if s.pointCount > 0 {
		rl.DrawCircle(int32(s.points[0].X), 1080, 6, rl.Red)
		for i := 1; i < len(s.points); i++ {
			for j := 0; j <= 2; j++ {
				v0 := rl.NewVector2(s.points[i-1].X, s.points[i-1].Y+float32(j)*30)
				v1 := rl.NewVector2(v0.X, float32(s.config.WindowHeight))
				v2 := rl.NewVector2(s.points[i].X, s.points[i].Y+float32(j)*30)
				v4 := rl.NewVector2(v2.X, float32(s.config.WindowHeight))
				color := rl.LightGray
				if j == 1 {
					color = rl.Gray
				} else if j == 2 {
					color = rl.DarkGray
				}
				rl.DrawTriangleStrip([]rl.Vector2{v0, v1, v2, v4}, color)
			}
			rl.DrawLineEx(s.points[i-1], s.points[i], 5, rl.RayWhite)
		}
	}
}

func (s *slope) scroll(speed float32) {
	for i := 0; i < int(s.pointCount); i++ {
		s.points[i].X -= speed
	}
}

func (s *slope) lastPoint() rl.Vector2 {
	return s.points[len(s.points)-1]
}

func (s *slope) heightAt(x float32) float32 {
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

func (s *slope) angleAt(x float32) float32 {
	for i := len(s.points) - 1; i >= 1; i-- {
		p1 := s.points[i-1]
		p2 := s.points[i]
		if p1.X <= x && p2.X >= x {
			dx := p2.X - p1.X
			dy := p2.Y - p1.Y
			if dx == 0 {
				return 0
			}
			return float32(math.Atan2(float64(dy), float64(dx)))
		}
	}
	return 0
}

func newSlope(config *Config) (s slope) {
	s.init()
	s.config = config
	return s
}
