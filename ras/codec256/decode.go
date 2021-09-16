package ras

import (
	"encoding/binary"
	"fmt"
	uuid "github.com/satori/go.uuid"
	pb "google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"math"
	"reflect"
	"time"
)

const (
	UTF8_CHARSET   = "UTF-8"
	SIZEOF_SHORT   = 2
	SIZEOF_INT     = 4
	SIZEOF_LONG    = 8
	NULL_BYTE      = 0x80
	TRUE_BYTE      = 1
	FALSE_BYTE     = 0
	MAX_SHIFT      = 7
	NULL_SHIFT     = 6
	BYTE_MASK      = 255
	NEXT_MASK      = -128
	NULL_NEXT_MASK = 64
	LAST_MASK      = 0
	NULL_LSB_MASK  = 63
	LSB_MASK       = 127
	TEMP_CAPACITY  = 256
)

const AgeDelta = 621355968000000

func ParseBytes(r io.Reader, data []byte) error {

	if len(data) == 0 {
		var err error
		data, err = io.ReadAll(r)
		if err != nil {
			return err
		}
		return nil
	}

	readLength := 0
	n := 0
	var err error

	for readLength < len(data) {

		n, err = r.Read(data[readLength:])
		readLength += n

		if err != nil {
			return err
		}
	}

	return nil
}

func ParseUUID(r io.Reader, into interface{}) error {

	buf := make([]byte, 16)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"uuid",
			err.Error(),
		}
	}

	u, err := uuid.FromBytes(buf)
	if err != nil {
		return &ParseError{
			"uuid",
			err.Error(),
		}
	}

	switch typed := into.(type) {
	case []byte:
		copy(typed, buf)
	case *[]byte:
		*typed = buf
	case *string:
		*typed = u.String()
	case *uuid.UUID:
		*typed = u
	default:
		return &ParseError{"uuid",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}

	return nil
}

func ParseTime(r io.Reader, into interface{}) error {

	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"time",
			err.Error(),
		}
	}

	val := binary.BigEndian.Uint64(buf)
	ticks := int64(val)
	timeT := (ticks - AgeDelta) / 10

	timestamp := time.Unix(0, timeT*int64(time.Millisecond)).UnixNano()

	switch typed := into.(type) {
	case *uint64:
		*typed = uint64(timestamp)
	case *int64:
		*typed = timestamp
	case *time.Time:
		*typed = time.Unix(0, timestamp)
	case *pb.Timestamp:
		*typed = *pb.New(time.Unix(0, timestamp))
	default:
		return &ParseError{"time",
			fmt.Sprintf("Parse time to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseType(r io.Reader, into interface{}) error {

	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
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
		return &ParseError{"type",
			fmt.Sprintf("Parse type to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseByte(r io.Reader, into interface{}) error {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"byte",
			err.Error(),
		}
	}

	b1 := buf[0]

	switch typed := into.(type) {
	case *byte:
		*typed = b1
	case *int8:
		*typed = int8(b1)
	case *int32:
		*typed = int32(b1)
	case *int:
		*typed = int(b1)
	default:
		return &ParseError{"byte",
			fmt.Sprintf("Parse byte to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseBool(r io.Reader, into interface{}) error {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
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
		return &ParseError{"bool",
			fmt.Sprintf("Parse byte to <%s> unsupporsed", typed)}
	}
	return nil

}

func ParseShort(r io.Reader, into interface{}) error {

	buf := make([]byte, SIZEOF_SHORT)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"short",
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
		return &ParseError{"uint16",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil

}

func ParseInt(r io.Reader, into interface{}) error {
	buf := make([]byte, SIZEOF_INT)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"int32",
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
	case bool:
		typed = false
		if val == 1 {
			typed = true
		}
	case *bool:
		*typed = false
		if val == 1 {
			*typed = true
		}
	default:
		return &ParseError{"int32",
			fmt.Sprintf("convert to <%s> unsupporsed", reflect.TypeOf(typed))}
	}
	return nil

}

func ParseLong(r io.Reader, into interface{}) error {
	buf := make([]byte, SIZEOF_LONG)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"long",
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
		return &ParseError{"uint64",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil

}

func ParseFloat(r io.Reader, into interface{}) error {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
			"float32",
			err.Error(),
		}
	}

	val := math.Float32frombits(binary.BigEndian.Uint32(buf))

	switch typed := into.(type) {
	case *float32:
		*typed = val
	case *float64:
		*typed = float64(val)
	default:
		return &ParseError{"float32",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseDouble(r io.Reader, into interface{}) error {

	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		return &ParseError{
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
		return &ParseError{"float64",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseString(r io.Reader, into interface{}) error {

	var size int

	err := ParseNullable(r, &size)
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	_, err = r.Read(buf)
	if err != nil {
		return &ParseError{"string",
			err.Error()}
	}

	switch typed := into.(type) {
	case *string:
		*typed = string(buf)
	case *[]byte:
		*typed = buf
	case []byte:
		copy(typed, buf)
	default:
		return &ParseError{"string",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}

	return nil
}

func ParseNullable(r io.Reader, into interface{}) error {

	readByte := func(fnName string) (int, byte, error) {
		buf := make([]byte, 1)
		n, err := r.Read(buf)
		if err != nil {
			return n, 0, &ParseError{
				fnName,
				err.Error(),
			}
		}
		b1 := buf[0]
		return n, b1, err
	}

	size := 0
	_, b1, err := readByte("nullable")

	if err != nil {
		return err
	}

	cur := int(b1 & 0xFF)
	if (cur & 0xFFFFFF80) == 0x0 {
		size = cur & 0x3F
		if cur&0x40 == 0x0 {
			return applyNullableS(size, into)
		}

		shift := NULL_SHIFT
		_, b1, err := readByte("nullable")
		if err != nil {
			return err
		}
		cur := int(b1 & 0xFF)
		size += (cur & 0x7F) << NULL_SHIFT
		shift += MAX_SHIFT

		for (cur & 0xFFFFFF80) != 0x0 {

			_, b1, err := readByte("nullable")
			if err != nil {
				return err
			}

			cur = int(b1 & 0xFF)
			size += (cur & 0x7F) << shift
			shift += MAX_SHIFT

		}
		return applyNullableS(size, into)
	}

	if (cur & 0x7F) != 0x0 {
		return &ParseError{
			"nullable",
			"null expected",
		}
	}

	return applyNullableS(size, into)
}

func applyNullableS(val int, into interface{}) error {
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
		return &ParseError{"nullable",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

func ParseSize(r io.Reader, into interface{}) error {

	readByte := func(fnName string) (int, byte, error) {
		buf := make([]byte, 1)
		n, err := r.Read(buf)
		if err != nil {
			return n, 0, &ParseError{
				fnName,
				err.Error(),
			}
		}
		b1 := buf[0]
		return n, b1, err
	}
	ff := 0xFFFFFF80
	_, b1, err := readByte("size")
	if err != nil {
		return err
	}

	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := MAX_SHIFT; (cur & ff) != 0x0; {

		_, b1, err = readByte("size")
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
		return &ParseError{"size",
			fmt.Sprintf("convert to <%s> unsupporsed", typed)}
	}
	return nil
}

type ParseError struct {
	Mame string
	Msg  string
}

func (e *ParseError) Error() string {
	return "ras: (ParserFunc " + e.Mame + ") " + e.Msg + ""
}
