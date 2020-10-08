package main

type caseName string

const (
	camel      caseName = "camel"
	upperCamel caseName = "upper-camel"
	kebab      caseName = "kebab"
	upperKebab caseName = "upper-kebab"
	snake      caseName = "snake"
	upperSnake caseName = "upper-snake"
	space      caseName = "space"
	upperSpace caseName = "upper-space"
)

var allPatternNames = map[caseName]struct{}{
	camel:      {},
	upperCamel: {},
	kebab:      {},
	upperKebab: {},
	snake:      {},
	upperSnake: {},
	space:      {},
	upperSpace: {},
}
