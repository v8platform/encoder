package ras

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

const TagNamespace = "rac"

type decodeFn func(dec *Decoder, rValue reflect.Value) error

type CodecField struct {
	Number   int
	Ignore   bool
	Version  int64
	Kind     CodecKind
	decoder  string
	fieldIdx int
}

func decodeValue(dec *Decoder, rValue reflect.Value) error {
	return dec.decodeValue(rValue)
}

// Unmarshal decodes the tag into a prototype.CodecField.
func unmarshalTag(tag string, fieldIdx int, goType reflect.Type) CodecField {

	f := CodecField{
		fieldIdx: fieldIdx,
	}

	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		f.Ignore = true
		return f
	}

	for idx, v := range tags {
		switch idx {
		case 0:
			switch v {
			case "-":
				f.Ignore = true
			default:
				f.decoder = v
			}
		case 1:
			n, _ := strconv.ParseInt(v, 10, 32)
			f.Number = int(n)
		case 2:
			n, _ := strconv.ParseInt(v, 10, 32)
			f.Version = n
		default:
			log.Fatalf("to many value in tag for field %s", goType.Name())
		}

	}
	return f
}
