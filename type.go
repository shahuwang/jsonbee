package main

import ()

type ItemType int

const (
	ItemArrayLeft ItemType = iota
	ItemArrayRight
	ItemDictLeft
	ItemDictRight
	ItemString
	ItemInteger
	ItemFloat
	ItemNull
	ItemError
	ItemComma
	ItemColon
)
