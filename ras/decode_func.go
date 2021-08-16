package ras

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

var decoderFunc = map[string]TypeDecoderFunc{}

type TypeDecoderFunc func(r io.Reader, into interface{}) error

func init() {
	RegisterDecoderType("time", decodeTime)
	RegisterDecoderType("type", decodeType)
	RegisterDecoderType("bool", decodeBool)
	RegisterDecoderType("byte", decodeByte)
	RegisterDecoderType("char short int16 uint16", decodeUint16)
	RegisterDecoderType("int int32 uint32", decodeUint32)
	RegisterDecoderType("int64 uint64", decodeUint64)
	RegisterDecoderType("float32", decodeFloat32)
	RegisterDecoderType("float64 double", decodeFloat64)
	RegisterDecoderType("string", decodeString)
	RegisterDecoderType("null-size", decodeNullableSize)
	RegisterDecoderType("size", decodeSize)
}

func RegisterDecoderType(name string, dec TypeDecoderFunc) {

	names := strings.Fields(strings.ToLower(name))

	for _, s := range names {
		decoderFunc[s] = dec
	}
}

func decodeTime(r io.Reader, into interface{}) error {

	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"time",
			err.Error(),
		}
	}

	val := binary.BigEndian.Uint64(buf)
	timeValue := dateFromTicks(int64(val))

	switch typed := into.(type) {
	case *uint64:
		*typed = uint64(timeValue.Unix())
	case *int64:
		*typed = timeValue.Unix()
	case *time.Time:
		*typed = timeValue
	default:
		return &TypeDecodeError{"time",
			fmt.Sprintf("decode time to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeType(r io.Reader, into interface{}) error {

	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"type",
			err.Error(),
		}
	}

	b1 := buf[0]
	cur := b1 & 0xFF

	switch typed := into.(type) {
	case *byte:
		*typed = cur
	default:
		return &TypeDecodeError{"type",
			fmt.Sprintf("decode byte to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeByte(r io.Reader, into interface{}) error {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"byte",
			err.Error(),
		}
	}

	b1 := buf[0]

	switch typed := into.(type) {
	case *byte:
		*typed = b1
	default:
		return &TypeDecodeError{"byte",
			fmt.Sprintf("decode byte to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeBool(r io.Reader, into interface{}) error {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"bool",
			err.Error(),
		}
	}

	b1 := buf[0]

	var val bool

	switch b1 {
	case TRUE_BYTE:
		val = true
	case FALSE_BYTE:
		val = false
	}

	switch typed := into.(type) {
	case *bool:
		*typed = val
	case *int:
		if val {
			*typed = 1
		} else {
			*typed = 0
		}
	default:
		return &TypeDecodeError{"bool",
			fmt.Sprintf("decode byte to <%s> unsupporsed", typed)}
	}
	return nil

}

func decodeUint16(r io.Reader, into interface{}) error {

	buf := make([]byte, 2)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"uint16",
			err.Error(),
		}
	}

	val := binary.BigEndian.Uint16(buf)

	switch typed := into.(type) {
	case *int:
		*typed = int(val)
	case *uint16:
		*typed = val
	case *int16:
		*typed = int16(val)
	case *uint32:
		*typed = uint32(val)
	case *int32:
		*typed = int32(val)
	case *uint64:
		*typed = uint64(val)
	case *int64:
		*typed = int64(val)
	default:
		return &TypeDecodeError{"uint16",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil

}

func decodeUint32(r io.Reader, into interface{}) error {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"uint32",
			err.Error(),
		}
	}

	val := binary.BigEndian.Uint32(buf)

	switch typed := into.(type) {
	case *int:
		*typed = int(val)
	case *uint16:
		*typed = uint16(val)
	case *int16:
		*typed = int16(val)
	case *uint32:
		*typed = uint32(val)
	case *int32:
		*typed = int32(val)
	case *uint64:
		*typed = uint64(val)
	case *int64:
		*typed = int64(val)
	default:
		return &TypeDecodeError{"uint32",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil

}

func decodeUint64(r io.Reader, into interface{}) error {
	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"uint64",
			err.Error(),
		}
	}

	val := binary.BigEndian.Uint64(buf)

	switch typed := into.(type) {
	case *int:
		*typed = int(val)
	case *uint16:
		*typed = uint16(val)
	case *int16:
		*typed = int16(val)
	case *uint32:
		*typed = uint32(val)
	case *int32:
		*typed = int32(val)
	case *uint64:
		*typed = uint64(val)
	case *int64:
		*typed = int64(val)
	default:
		return &TypeDecodeError{"uint64",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil

}

func decodeFloat32(r io.Reader, into interface{}) error {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"float32",
			err.Error(),
		}
	}

	val := math.Float32frombits(binary.BigEndian.Uint32(buf))

	switch typed := into.(type) {
	case *float32:
		*typed = float32(val)
	case *float64:
		*typed = float64(val)
	default:
		return &TypeDecodeError{"float32",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeFloat64(r io.Reader, into interface{}) error {

	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{
			"float64",
			err.Error(),
		}
	}

	val := math.Float64frombits(binary.BigEndian.Uint64(buf))

	switch typed := into.(type) {
	case *float32:
		*typed = float32(val)
	case *float64:
		*typed = float64(val)
	default:
		return &TypeDecodeError{"float64",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeString(r io.Reader, into interface{}) error {

	var size int
	err := decodeNullableSize(r, &size)
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	n, err := r.Read(buf)
	if err != nil {
		return &TypeDecodeError{"string",
			fmt.Sprintf("read bytes<%d> err: <%s>", n, err.Error())}
	}

	switch typed := into.(type) {
	case *string:
		*typed = string(buf)
	case *[]byte:
		*typed = buf
	case []byte:
		copy(typed, buf)
	default:
		return &TypeDecodeError{"string",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}

	return nil
}

func decodeNullableSize(r io.Reader, into interface{}) error {

	readByte := func(fnName string) (byte, error) {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err != nil {
			return 0, &TypeDecodeError{
				fnName,
				err.Error(),
			}
		}
		b1 := buf[0]
		return b1, err
	}

	size := 0
	b1, err := readByte("nullableSize")
	if err != nil {
		return err
	}

	cur := int(b1 & 0xFF)
	if (cur & 0xFFFFFF80) == 0x0 {
		size = cur & 0x3F
		if cur&0x40 == 0x0 {
			return applyNullableSize(size, into)
		}

		shift := NULL_SHIFT
		b1, err := readByte("nullableSize")
		if err != nil {
			return err
		}
		cur := int(b1 & 0xFF)
		size += (cur & 0x7F) << NULL_SHIFT
		shift += MAX_SHIFT

		for (cur & 0xFFFFFF80) != 0x0 {

			b1, err := readByte("nullableSize")
			if err != nil {
				return err
			}

			cur = int(b1 & 0xFF)
			size += (cur & 0x7F) << shift
			shift += MAX_SHIFT

		}
		return applyNullableSize(size, into)
	}

	if (cur & 0x7F) != 0x0 {
		return &TypeDecodeError{
			"nullableSize",
			"null expected",
		}
	}

	return applyNullableSize(size, into)
}

func applyNullableSize(val int, into interface{}) error {
	switch typed := into.(type) {
	case *int:
		*typed = int(val)
	case *uint16:
		*typed = uint16(val)
	case *int16:
		*typed = int16(val)
	case *uint32:
		*typed = uint32(val)
	case *int32:
		*typed = int32(val)
	case *uint64:
		*typed = uint64(val)
	case *int64:
		*typed = int64(val)
	default:
		return &TypeDecodeError{"nullableSize",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func decodeSize(r io.Reader, into interface{}) error {

	readByte := func(fnName string) (byte, error) {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err != nil {
			return 0, &TypeDecodeError{
				fnName,
				err.Error(),
			}
		}
		b1 := buf[0]
		return b1, err
	}

	ff := 0xFFFFFF80
	b1, err := readByte("size")
	if err != nil {
		return err
	}

	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := MAX_SHIFT; (cur & ff) != 0x0; {

		b1, err = readByte("size")
		if err != nil {
			return err
		}

		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += MAX_SHIFT
	}

	switch typed := into.(type) {
	case *int:
		*typed = int(size)
	case *uint16:
		*typed = uint16(size)
	case *int16:
		*typed = int16(size)
	case *uint32:
		*typed = uint32(size)
	case *int32:
		*typed = int32(size)
	case *uint64:
		*typed = uint64(size)
	case *int64:
		*typed = int64(size)
	default:
		return &TypeDecodeError{"size",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

type TypeDecodeError struct {
	Mame string
	Msg  string
}

func (e *TypeDecodeError) Error() string {
	// if e.Type == nil {
	// 	return "ras: Decode(nil)"
	// }

	// if e.Type.Kind() != reflect.Ptr {
	// 	return "ras: Decode(non-pointer " + e.Type.String() + ")"
	// }
	return "ras: (decoderFunc " + e.Mame + ") " + e.Msg + ""
}
