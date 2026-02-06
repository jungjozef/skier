package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Physics constants
const (
	gravityAccel     = 0.5
	baseScrollSpeed  = 18.0
	maxScrollSpeed   = 55.0
	minScrollSpeed   = 6.0
	slopeAccelFactor = 0.35
	friction         = 0.003
	launchFactor     = 1.2
	jumpFactor       = 0.8
	terminalVelocity = 30.0
)

type skier struct {
	position  rl.Vector2
	texture   rl.Texture2D
	velocityY float32
	rotation  float32
	mountain  *mountain
	onSlope   bool
}

func (s *skier) draw() {
	source := rl.NewRectangle(0, 0, float32(s.texture.Width), float32(s.texture.Height))
	dest := rl.NewRectangle(s.position.X, s.position.Y, float32(s.texture.Width), float32(s.texture.Height))
	origin := rl.NewVector2(float32(s.texture.Width)/2, float32(s.texture.Height)/2)
	rl.DrawTexturePro(s.texture, source, dest, origin, s.rotation*180/math.Pi, rl.White)
}

func (s *skier) init() {
	s.position = rl.NewVector2(400, 20)
	img := rl.GenImageGradientRadial(64, 64, 0.5, rl.Yellow, rl.Blank)
	s.texture = rl.LoadTextureFromImage(img)
	s.velocityY = 0
	s.rotation = 0
	s.onSlope = false
}

func (s *skier) reset() {
	s.position = rl.NewVector2(400, 20)
	s.velocityY = 0
	s.rotation = 0
	s.onSlope = false
}

func (s *skier) update(scrollSpeed float32) float32 {
	h, sl := s.mountain.heightAt(s.position.X)

	if s.onSlope {
		if sl == nil || h == -1 {
			// Fell off end of slope — launch using last known slope angle
			angle := s.rotation // preserved from last frame on slope
			s.velocityY = float32(math.Sin(float64(angle))) * scrollSpeed * launchFactor
			s.onSlope = false
		} else {
			s.position.Y = h
			angle := s.mountain.angleAt(s.position.X)
			s.rotation = angle
			// Slope acceleration
			scrollSpeed += float32(math.Sin(float64(angle))) * slopeAccelFactor * scrollSpeed
			scrollSpeed *= (1 - friction)
			if scrollSpeed > maxScrollSpeed {
				scrollSpeed = maxScrollSpeed
			}
			if scrollSpeed < minScrollSpeed {
				scrollSpeed = minScrollSpeed
			}
			// Right-click jump — strength scales with momentum
			if rl.IsMouseButtonPressed(rl.MouseRightButton) {
				jumpStrength := scrollSpeed * jumpFactor
				s.velocityY = -jumpStrength
				// Uphill ramp bonus: steeper uphill = bigger launch
				if angle < 0 {
					s.velocityY += float32(math.Sin(float64(angle))) * scrollSpeed * 0.5
				}
				s.onSlope = false
			}
		}
	} else {
		// Airborne
		s.velocityY += gravityAccel
		if s.velocityY > terminalVelocity {
			s.velocityY = terminalVelocity
		}
		s.position.Y += s.velocityY

		// Visual rotation based on velocity ratio
		targetRot := float32(math.Atan2(float64(s.velocityY), float64(scrollSpeed)))
		s.rotation += (targetRot - s.rotation) * 0.1

		// Check for landing
		if h != -1 && s.position.Y >= h && s.velocityY > 0 {
			// Capture approach angle before zeroing velocity
			approachAngle := float32(math.Atan2(float64(s.velocityY), float64(scrollSpeed)))
			s.position.Y = h
			s.onSlope = true
			s.velocityY = 0
			// Landing speed adjustment based on angle difference
			slopeAngle := s.mountain.angleAt(s.position.X)
			angleDiff := float32(math.Abs(float64(slopeAngle - approachAngle)))
			if angleDiff < 0.5 {
				scrollSpeed *= 1.1 // Smooth landing bonus
			} else {
				scrollSpeed *= 0.95 // Rough landing penalty
			}
			if scrollSpeed > maxScrollSpeed {
				scrollSpeed = maxScrollSpeed
			}
			if scrollSpeed < minScrollSpeed {
				scrollSpeed = minScrollSpeed
			}
			s.rotation = slopeAngle
		}
	}

	return scrollSpeed
}

func newSkier(m *mountain) (s skier) {
	s.init()
	s.mountain = m
	return s
}
