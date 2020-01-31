package main

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
	"math/rand"
	"time"
)

type Area struct {
	w, h   int
	matrix [][]pkg.FieldType
	r      *rand.Rand
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
		r:      rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (a *Area) SetByPos(pos *pkg.Pos, field pkg.FieldType) {
	a.matrix[pos.X][pos.Y] = field
}

func (a *Area) GetRelative(pos *pkg.Pos, diff pkg.Pos) pkg.FieldType {
	return a.matrix[pos.X+diff.X][pos.Y+diff.Y]
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
	isFound := false
	results := make([]*pkg.Pos, 0, 1)

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
						results = append(results, &pkg.Pos{X: x, Y: y})
						isFound = true
					}

					// debug
					//out := ""
					//for x := range mm{
					//	out += strings.Join(mm[x], "") + "\n"
					//}
					//fmt.Printf("\x0c%s", out)
					//out = ""
					//time.Sleep(800 * time.Millisecond)
				}

				xLine = !xLine
				lastPolarity = polarity
			}
		}

		if isFound {
			return results[a.r.Intn(len(results))]
		}
	}

	return nil
}

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func (a *Area) RewriteMap(nX, nY int, ai *AI) bool {
	if nX >= 0 && nY >= 0 && nX < a.w && nY < a.h {
		return false
	}

	if nX >= a.w {
		for x := a.w; x < a.w*2; x++ {
			a.matrix = append(a.matrix, make([]pkg.FieldType, a.h))
			for y := range a.matrix[x] {
				a.matrix[x][y] = unknownField
			}
		}

		a.w *= 2
	}

	if nY >= a.h {
		for y := a.h; y < a.h*2; y++ {
			for x := range a.matrix {
				a.matrix[x] = append(a.matrix[x], unknownField)
			}
		}

		a.h *= 2
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
		ai.moveObjects(&pkg.Pos{X: a.w})
		a.w *= 2
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

		ai.moveObjects(&pkg.Pos{Y: a.h})
		a.h *= 2
		a.matrix = expandedMatrix
	}

	return true
}

func (a *Area) CutArea(fields [5][5]pkg.FieldType, antPos *pkg.Pos, ai *AI) {
	if fields[0][2] == pkg.WallField && antPos.X-2 > 1 {
		diff := antPos.X - 3
		a.matrix = a.matrix[diff:]
		ai.moveObjects(&pkg.Pos{X: -diff})
		a.w -= diff
	}

	if fields[2][0] == pkg.WallField && antPos.Y-2 > 1 {
		diff := antPos.Y - 3
		for x := range a.matrix {
			a.matrix[x] = a.matrix[x][diff:]
		}

		ai.moveObjects(&pkg.Pos{Y: -diff})
		a.h -= diff
	}

	if fields[4][2] == pkg.WallField && antPos.X+2 < a.w-2 {
		a.w = antPos.X + 4
		a.matrix = a.matrix[:a.w]
	}

	if fields[2][4] == pkg.WallField && antPos.Y+2 < a.h-2 {
		a.h = antPos.Y + 4
		for x := range a.matrix {
			a.matrix[x] = a.matrix[x][:a.h]
		}
	}
}
