package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type Ant struct {
	Pos   *pkg.Pos
	Role  Role
	Order Order
}

func (a *Ant) RelativePos(x int, y int) *pkg.Pos {
	return &pkg.Pos{
		X: x - a.Pos.X + 2,
		Y: y - a.Pos.Y + 2,
	}
}

func (a *Ant) RelativeNearestPos(x int, y int) *pkg.Pos {
	targetX := x - a.Pos.X + 2
	targetY := y - a.Pos.Y + 2
	if targetX > 1 || targetX < -1 {
		targetX = int(float64(targetX) / math.Abs(float64(targetX)))
	}
	if targetY > 1 || targetY < -1 {
		targetY = int(float64(targetY) / math.Abs(float64(targetY)))
	}

	return &pkg.Pos{
		X: targetX,
		Y: targetY,
	}
}
