package main

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type Area struct {
	w, h   int
	matrix [][]pkg.FieldType
}

func NewArea(w, h int) *Area {
	matrix := make([][]pkg.FieldType, w)
	for x := range matrix {
		matrix[x] = make([]pkg.FieldType, h)
		for y := range matrix[x] {
			matrix[x][y] = unknownField
		}
	}

	return &Area{
		w:      w,
		h:      h,
		matrix: matrix,
	}
}

func (a *Area) SetByPos(pos *pkg.Pos, field pkg.FieldType) {
	a.matrix[pos.X][pos.Y] = field
}

func (a *Area) Closest(point *pkg.Pos, sought pkg.FieldType) *pkg.Pos {
	// limitations for every direction in order: > v < ^
	directions := []int{
		a.w - 1 - point.X, a.h - 1 - point.Y,
		point.X - 1, point.Y - 1,
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

	// debug
	//mm := make([][]string, a.h)
	//for i := range mm {
	//	mm[i] = make([]string, a.w)
	//	for j := range mm[i] {
	//		if a.matrix[j][i] == sought {
	//			mm[i][j] = "x"
	//		} else {
	//			mm[i][j] = "-"
	//		}
	//	}
	//}
	//mm[point.Y][point.X] = "p"
	for frame := 1; frame <= lastFrame; frame++ {
		/* s - start (x, y)
		^s-->
		|-p-|
		<---v
		*/
		baseFrom := -(frame - 1)
		x := baseFrom + point.X
		y := -frame + point.Y
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
				// if we trying to check x line out of map (negative)
				if xLine == true && y < 0 {
					x = frame*polarity + point.X
					xLine = !xLine
					lastPolarity = polarity
					continue
				}
				// if we trying to check y line out of map (negative)
				if xLine == false && x < 0 {
					y = frame*polarity + point.Y
					xLine = !xLine
					lastPolarity = polarity
					continue
				}

				// if we trying to check line out of map (positive)
				if frame > limit[!xLine][lastPolarity] {
					if xLine {
						x = frame*polarity + point.X
					} else {
						y = frame*polarity + point.Y
					}
					xLine = !xLine
					lastPolarity = polarity
					continue
				}

				var to int
				if limit[xLine][polarity] > frame {
					to = frame
				} else {
					to = limit[xLine][polarity]
				}

				from := baseFrom
				if limit[xLine][polarity*-1] < int(math.Abs(float64(from))) {
					from = -limit[xLine][polarity*-1]
				}

				for j := from; j <= to; j++ {
					if xLine {
						x = j*polarity + point.X
					} else {
						y = j*polarity + point.Y
					}

					//mm[y][x] = "+"
					if a.matrix[x][y] == sought {
						return &pkg.Pos{X: x, Y: y}
					}

					// debug
					//out := ""
					//for x := range mm{
					//	out = out + strings.Join(mm[x], "") + "\n"
					//}
					//fmt.Printf("\x0c%s", out)
					//out = ""
					//time.Sleep(800 * time.Millisecond)
				}

				xLine = !xLine
				lastPolarity = polarity
			}
		}
	}

	return nil
}

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func (a *Area) RewriteMap(nX, nY int, t pkg.FieldType, ai *AI) bool {
	if t == pkg.NoField || (nX >= 0 && nY >= 0 && nX < a.w && nY < a.h) {
		return false
	}

	// todo handle wall and noField case, when we need to make area shorter

	if nX >= a.w {
		for x := a.w; x < a.w*2; x++ {
			a.matrix = append(a.matrix, make([]pkg.FieldType, a.h))
			for y := range a.matrix[x] {
				a.matrix[x][y] = unknownField
			}
		}

		a.w = a.w * 2
	}

	if nY >= a.h {
		for y := a.h; y < a.h*2; y++ {
			for x := range a.matrix {
				a.matrix[x] = append(a.matrix[x], unknownField)
			}
		}

		a.h = a.h * 2
	}

	if nX < 0 {
		expandedMatrix := make([][]pkg.FieldType, a.w, a.w*2)
		// fill new part
		for x := 0; x < a.w; x++ {
			expandedMatrix[x] = make([]pkg.FieldType, a.h)
			for y := range expandedMatrix[x] {
				expandedMatrix[x][y] = unknownField
			}
		}

		// fill old part
		expandedMatrix = append(expandedMatrix, a.matrix...)

		posDiff := &pkg.Pos{X: a.w}
		for _, ant := range ai.ants {
			ant.Pos.Add(posDiff)
		}
		for _, anthill := range ai.anthills {
			anthill.Pos.Add(posDiff)
		}
		for _, pos := range ai.enemyAnthills {
			pos.Add(posDiff)
		}

		a.w = a.w * 2
		a.matrix = expandedMatrix
	}

	if nY < 0 {
		expandedMatrix := make([][]pkg.FieldType, a.w)
		// fill new part
		for x := range expandedMatrix {
			expandedMatrix[x] = make([]pkg.FieldType, a.h*2)
			for y := 0; y < a.h; y++ {
				expandedMatrix[x][y] = unknownField
			}
		}

		// fill old part
		for x := range expandedMatrix {
			for y := a.h; y < a.h*2; y++ {
				expandedMatrix[x][y] = a.matrix[x][y-a.h]
			}
		}

		posDiff := &pkg.Pos{Y: a.h}
		for _, ant := range ai.ants {
			ant.Pos.Add(posDiff)
		}
		for _, anthill := range ai.anthills {
			anthill.Pos.Add(posDiff)
		}
		for _, pos := range ai.enemyAnthills {
			pos.Add(posDiff)
		}

		a.h = a.h * 2
		a.matrix = expandedMatrix
	}

	return true
}
