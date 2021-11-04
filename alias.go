package it

import (
	"github.com/reusee/dscope"
	"github.com/reusee/e4"
	"github.com/reusee/pp"
)

type (
	Scope = dscope.Scope
	Src   = pp.Src
	Sink  = pp.Sink
)

var (
	we = e4.Wrap
)
