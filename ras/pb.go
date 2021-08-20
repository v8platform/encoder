package ras

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sort"
)

func encode(message proto.Message) []byte {

	//message, ok := m.(proto.Message)

	//if !ok {
	//	panic("Non proto message")
	//}

	buf := &bytes.Buffer{}

	type encodeField struct {
		order   int32
		encoder string
		value   interface{}
	}

	var fields []encodeField

	message.ProtoReflect().Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {

		encoderOptions := proto.GetExtension(descriptor.Options(), extpb.E_Field).(*extpb.EncodingFieldOptions)

		fields = append(fields, encodeField{
			order:   encoderOptions.Order,
			encoder: *encoderOptions.Encoder,
			value:   value.Interface(),
		})
		return true
	})

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].order < fields[j].order
	})

	//pp.Println(fields)

	for _, field := range fields {

		_, err := EncodeValue(field.encoder, buf, field.value)
		if err != nil {
			panic(err)
		}
	}

	return buf.Bytes()
}
