package main

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type Ant struct {
	ID    int
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

func (a *Ant) CalcStep(x int, y int) (*pkg.Pos, bool) {
	isGoal := true
	targetX := x - a.Pos.X
	targetY := y - a.Pos.Y
	if targetX > 1 || targetX < -1 {
		isGoal = false
		targetX = int(float64(targetX) / math.Abs(float64(targetX)))
	}
	if targetY > 1 || targetY < -1 {
		isGoal = false
		targetY = int(float64(targetY) / math.Abs(float64(targetY)))
	}

	return &pkg.Pos{
		X: targetX,
		Y: targetY,
	}, isGoal
}
