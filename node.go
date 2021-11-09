package it

import (
	"io"
)

type Node interface {
	NodeID() int64
	KeyRange() (Key, Key)
	Mutate(
		ctx Scope,
		path KeyPath,
		fn func(Node) (Node, error),
	) (
		retNode Node,
		err error,
	)
	Merge(
		ctx Scope,
		node2 Node,
	) (
		newNode Node,
		err error,
	)
	Dump(w io.Writer, level int)
	Walk(cont Src) Src
}
