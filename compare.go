package it

import (
	"strings"

	"github.com/reusee/sb"
)

type Comparable interface {
	Cmp(b any) int
}

func Compare(a, b any) int {
	switch a := a.(type) {
	case string:
		if b, ok := b.(string); ok {
			return strings.Compare(a, b)
		}
	case uint8:
		if b, ok := b.(uint8); ok {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}
	case Comparable:
		return a.Cmp(b)
	}
	return sb.MustCompare(
		sb.Marshal(a),
		sb.Marshal(b),
	)
}
