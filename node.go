package it

import (
	"io"
)

type Node interface {
	KeyRange() (Key, Key)
	Mutate(
		ctx Scope,
		path KeyPath,
		fn func(Node) (Node, error),
	) (
		retNode Node,
		err error,
	)
	Dump(w io.Writer, level int)
	Walk(cont Src) Src
}
