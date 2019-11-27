package example

import pkg "github.com/gregmus2/ants-pkg"

type Ant struct {
	Pos   *pkg.Pos
	Role  Role
	Order Order
}

func (a *Ant) RelativePos(x int, y int) *pkg.Pos {
	return &pkg.Pos{
		X: x - a.Pos.X + 3,
		Y: y - a.Pos.Y + 3,
	}
}
