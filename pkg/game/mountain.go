package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type mountain struct {
	slopes []slope
	config *Config
}

func (m *mountain) init() {
	m.slopes = make([]slope, 0)
}

func (m *mountain) draw() {
	for i := 0; i < len(m.slopes); i++ {
		m.slopes[i].draw()
	}
}

func (m *mountain) update(scrollSpeed float32) {
	p := rl.GetMousePosition()
	for i := 0; i < len(m.slopes); i++ {
		m.slopes[i].scroll(scrollSpeed)
	}

	if (p.X >= 0 && p.X <= float32(m.config.WindowWidth)) && (p.Y >= 0 && p.Y <= float32(m.config.WindowHeight)) {
		existingH, _ := m.heightAt(p.X)
		slopeExists := existingH != -1

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !slopeExists {
			m.slopes = append(m.slopes, newSlope(m.config))
			m.slopes[len(m.slopes)-1].add()
			m.slopes[len(m.slopes)-1].active = true
		}
		if rl.IsMouseButtonDown(rl.MouseLeftButton) && len(m.slopes) > 0 && m.slopes[len(m.slopes)-1].active {
			if p.X > m.slopes[len(m.slopes)-1].lastPoint().X && !slopeExists {
				m.slopes[len(m.slopes)-1].add()
			}
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) && len(m.slopes) > 0 {
			m.slopes[len(m.slopes)-1].active = false
		}
	}
	var activeSlopes = make([]slope, 0)
	for _, s := range m.slopes {
		if s.lastPoint().X >= 0 {
			activeSlopes = append(activeSlopes, s)
		}
	}
	m.slopes = activeSlopes
}

func (m *mountain) heightAt(x float32) (float32, *slope) {
	for i := len(m.slopes) - 1; i >= 0; i-- {
		height := m.slopes[i].heightAt(x)
		if height != -1 {
			return height, &m.slopes[i]
		}
	}
	return -1, nil
}

func (m *mountain) angleAt(x float32) float32 {
	for i := len(m.slopes) - 1; i >= 0; i-- {
		h := m.slopes[i].heightAt(x)
		if h != -1 {
			return m.slopes[i].angleAt(x)
		}
	}
	return 0
}

func newMountain(config *Config) (m mountain) {
	m.init()
	m.config = config
	return m
}
