// Package computer_use provides autonomous computer use agent.
package computer_use

import (
	"math"
)

// CoordinateSystem normalizes coordinates to 0-999 grid.
type CoordinateSystem struct {
	ScreenWidth  int
	ScreenHeight int
}

// NewCoordinateSystem creates a new coordinate system.
func NewCoordinateSystem(width, height int) *CoordinateSystem {
	return &CoordinateSystem{ScreenWidth: width, ScreenHeight: height}
}

// NormalizeToGrid converts pixel coordinates to 0-999 grid.
func (c *CoordinateSystem) NormalizeToGrid(x, y int) (int, int) {
	if c.ScreenWidth == 0 || c.ScreenHeight == 0 {
		return 0, 0
	}
	normX := clamp(int(math.Round(float64(x)*999/float64(c.ScreenWidth))), 0, 999)
	normY := clamp(int(math.Round(float64(y)*999/float64(c.ScreenHeight))), 0, 999)
	return normX, normY
}

// DenormalizeFromGrid converts 0-999 grid back to pixel coordinates.
func (c *CoordinateSystem) DenormalizeFromGrid(normX, normY int) (int, int) {
	if c.ScreenWidth == 0 || c.ScreenHeight == 0 {
		return 0, 0
	}
	x := clamp(int(math.Round(float64(normX)*float64(c.ScreenWidth)/999)), 0, c.ScreenWidth)
	y := clamp(int(math.Round(float64(normY)*float64(c.ScreenHeight)/999)), 0, c.ScreenHeight)
	return x, y
}

// GridRegion represents a region in normalized coordinates.
type GridRegion struct {
	X, Y   int
	Width  int
	Height int
}

// RegionToGrid converts a pixel region to normalized grid region.
func (c *CoordinateSystem) RegionToGrid(x1, y1, x2, y2 int) GridRegion {
	nx1, ny1 := c.NormalizeToGrid(x1, y1)
	nx2, ny2 := c.NormalizeToGrid(x2, y2)
	return GridRegion{X: nx1, Y: ny1, Width: nx2 - nx1, Height: ny2 - ny1}
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// NormalizedPoint represents a point in 0-999 grid.
type NormalizedPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ToPixels converts normalized point to pixel coordinates.
func (p NormalizedPoint) ToPixels(cs *CoordinateSystem) (int, int) {
	return cs.DenormalizeFromGrid(p.X, p.Y)
}

// FromPixels creates normalized point from pixel coordinates.
func FromPixels(x, y int, cs *CoordinateSystem) NormalizedPoint {
	nx, ny := cs.NormalizeToGrid(x, y)
	return NormalizedPoint{X: nx, Y: ny}
}
