package main

type Role uint8

const (
	explorer Role = iota
	defender
	attacker
)
