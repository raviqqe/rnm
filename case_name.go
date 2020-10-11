package main

type caseName string

const (
	bare       caseName = "bare"
	camel      caseName = "camel"
	upperCamel caseName = "upper-camel"
	kebab      caseName = "kebab"
	upperKebab caseName = "upper-kebab"
	snake      caseName = "snake"
	upperSnake caseName = "upper-snake"
	space      caseName = "space"
	upperSpace caseName = "upper-space"
)

var allCaseNames = map[caseName]struct{}{
	bare:       {},
	camel:      {},
	upperCamel: {},
	kebab:      {},
	upperKebab: {},
	snake:      {},
	upperSnake: {},
	space:      {},
	upperSpace: {},
}
