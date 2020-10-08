package main

type patternName string

const (
	camel      patternName = "camel"
	upperCamel patternName = "upper-camel"
	kebab      patternName = "kebab"
	upperKebab patternName = "upper-kebab"
	snake      patternName = "snake"
	upperSnake patternName = "upper-snake"
	space      patternName = "space"
	upperSpace patternName = "upper-space"
)

var allPatternNames = map[patternName]struct{}{
	camel:      {},
	upperCamel: {},
	kebab:      {},
	upperKebab: {},
	snake:      {},
	upperSnake: {},
	space:      {},
	upperSpace: {},
}
