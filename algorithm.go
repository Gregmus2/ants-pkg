package pkg

type Algorithm interface {
	Do(antID int, fields [5][5]FieldType, round int, posDiff *Pos) (target *Pos, action Action)
	OnAntDie(antID int)
	OnAnthillDie(anthillID int)
	OnAntBirth(antID int, anthillID int)
	OnNewAnthill(invaderID int, birthPos *Pos) // antID; position relative anthill
}
