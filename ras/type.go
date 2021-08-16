package ras

import "reflect"

type Type struct {
	kind reflect.Kind

	fields []*field
}

type field struct {
	Number  int
	Ignore  bool
	Version int64
	Kind    CodecKind

	fieldIdx int
}
