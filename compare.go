package it

import "github.com/reusee/sb"

func Compare(a, b any) int {
	return sb.MustCompare(
		sb.Marshal(a),
		sb.Marshal(b),
	)
}
