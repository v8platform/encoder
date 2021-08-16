package ras

import (
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const TagNamespace = "rac"

type CodecField struct {
	Number   int
	Ignore   bool
	Version  int
	codec    string
	fieldIdx int
}

func getCodecFields(rType reflect.Type) []CodecField {
	if _, ok := rType.(reflect.Type); !ok {
		rType = rType.Elem()
	}

	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	fieldsCount := rType.NumField()

	var fields []CodecField

	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		tag := field.Tag.Get(TagNamespace)

		codecField := unmarshalTag(tag, i, rType)
		fields = append(fields, codecField)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Number < fields[j].Number
	})

	return fields
}

// Unmarshal decodes the tag into a prototype.CodecField.
func unmarshalTag(tag string, fieldIdx int, rType reflect.Type) CodecField {

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
				f.codec = v
			}
		case 1:
			n, _ := strconv.ParseInt(v, 10, 32)
			f.Number = int(n)
		case 2:
			n, _ := strconv.ParseInt(v, 10, 32)
			f.Version = int(n)
		default:
			log.Fatalf("to many value in tag for field %s", rType.Name())
		}

	}
	return f
}
