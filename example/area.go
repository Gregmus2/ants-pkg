package example

import pkg "github.com/gregmus2/ants-pkg"

type Area [][]pkg.FieldType

func NewArea(w, h int) Area {
	area := make(Area, w)
	for x := range area {
		area[x] = make([]pkg.FieldType, h)
		for y := range area[x] {
			area[x][y] = unknownField
		}
	}

	return area
}

func (a Area) Closest(point *pkg.Pos, sought pkg.FieldType) *pkg.Pos {
	// limitations for every direction in order: > v < ^
	directions := []int{
		len(a) - 1 - point.X, len(a[0]) - 1 - point.Y,
		point.X, point.Y,
	}
	// map[X|Y][polarity(-1|1)] = limit for this direction
	limit := map[bool]map[int]int{
		true:  {1: directions[0], -1: directions[2]},
		false: {1: directions[1], -1: directions[3]},
	}

	// the last frame in which we need to search before the array ends
	lastFrame := directions[0]
	for i := 1; i < len(directions); i++ {
		if directions[i] > lastFrame {
			lastFrame = directions[i]
		}
	}
	xLine := true
	lastPolarity := -1
	for frame := 1; frame <= lastFrame; frame++ {
		baseFrom := -(frame - 1)
		x := baseFrom + point.X
		y := -frame + point.X
		/*
			---++++++
			----+++++
			----p++++
			-----++++
			------+++
		*/
		for polarity := 1; polarity >= -1; polarity -= 2 {
			/*
				y x x x x
				y y x x y
				y y p y y
				y x x y y
				x x x x y
			*/
			for axis := 0; axis <= 1; axis++ {
				if frame > limit[!xLine][lastPolarity] {
					continue
				}

				var to int
				if limit[xLine][polarity] > frame {
					to = frame
				} else {
					to = limit[xLine][polarity]
				}

				from := baseFrom
				if limit[xLine][polarity*-1] > from {
					from = limit[xLine][polarity*-1]
				}

				for j := from; j <= to; j++ {
					if xLine {
						x = j*polarity + point.X
					} else {
						y = j*polarity + point.Y
					}

					if a[x][y] == sought {
						return &pkg.Pos{X: x, Y: y}
					}
				}

				xLine = !xLine
				lastPolarity = polarity
			}
		}
	}
}
