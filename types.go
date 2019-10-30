package pkg

import "math"

type Algorithm interface {
	Start(anthill Pos, birth Pos)
	Do(fields [5][5]FieldType, round int) (target Pos, action Action)
}

type Pos [2]uint

func (p *Pos) X() uint {
	return p[0]
}

func (p *Pos) Y() uint {
	return p[1]
}

/* add field with format
	0 1 2
	3 4 5
	6 7 8
to input position
*/
func (p Pos) RelativePosition(field uint8) Pos {
	return Pos{
		p.X() + uint(math.Mod(float64(field+3), 3)-1),
		p.Y() + uint(math.Floor(float64(field/3))-1), //nolint
	}
}
