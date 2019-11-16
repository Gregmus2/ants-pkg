package pkg

type FieldType uint8

const (
	NoField = iota
	EmptyField
	FoodField
	AllyField
	EnemyField
	WallField
	AntField
	AnthillField
)

type Action uint8

const (
	NoAction = iota
	AttackAction
	EatAction
	MoveAction
	DieAction
)
