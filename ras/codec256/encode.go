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

func FormatUuid(r io.Writer, value interface{}) error {

	switch val := value.(type) {
	case []byte:
		return writeBuf("uuid", r, val)
	case *[]byte:
		return writeBuf("uuid", r, *val)
	case *uuid.UUID:
		return writeBuf("uuid", r, val.Bytes())
	case uuid.UUID:
		return writeBuf("uuid", r, val.Bytes())
	case string:
		return writeBuf("uuid", r, uuid.FromStringOrNil(val).Bytes())
	case *string:
		return writeBuf("uuid", r, uuid.FromStringOrNil(*val).Bytes())
	default:
		return &TypeEncoderError{"uuid", "unknown uuid type"}
	}

}

func FormatBytes(r io.Writer, value interface{}) error {

	switch val := value.(type) {
	case []byte:
		return writeBuf("bytes", r, val)
	case *[]byte:
		return writeBuf("bytes", r, *val)
	case *uuid.UUID:
		return writeBuf("bytes", r, val.Bytes())
	case uuid.UUID:
		return writeBuf("bytes", r, val.Bytes())
	case string:
		return writeBuf("bytes", r, []byte(val))
	case *string:
		return writeBuf("bytes", r, []byte(*val))
	default:
		return &TypeEncoderError{"bytes", "unknown bytes type"}
	}

}

func FormatTime(w io.Writer, value interface{}) error {
	var val int64

	switch tVal := value.(type) {
	case int64:
		val = int64(tVal)
	case uint64:
		val = int64(tVal)
	case *int64:
		val = int64(*tVal)
	case *uint64:
		val = int64(*tVal)
	case time.Time:
		val = tVal.UnixNano()
	case *time.Time:
		val = tVal.UnixNano()
	case pb.Timestamp:
		val = tVal.AsTime().UnixNano()
	case *pb.Timestamp:
		val = tVal.AsTime().UnixNano()
	default:
		return &TypeEncoderError{"time", "TODO"}
	}
	ticks := val / int64(time.Millisecond)
	ticks = ticks*10 + AgeDelta

	return FormatLong(w, ticks)

}

func FormatShort(w io.Writer, value interface{}) error {
	var val uint16

	switch tVal := value.(type) {
	case int16:
		val = uint16(tVal)
	case uint16:
		val = uint16(tVal)
	case *int16:
		val = uint16(*tVal)
	case *uint16:
		val = uint16(*tVal)
	default:
		return &TypeEncoderError{"short", "TODO"}
	}
	buf := make([]byte, SIZEOF_SHORT)
	binary.BigEndian.PutUint16(buf, val)
	return writeBuf("short", w, buf)

}

func FormatInt(w io.Writer, value interface{}) error {
	var val uint32

	switch tVal := value.(type) {
	case int:
		val = uint32(tVal)
	case uint:
		val = uint32(tVal)
	case *int:
		val = uint32(*tVal)
	case *uint:
		val = uint32(*tVal)
	case int32:
		val = uint32(tVal)
	case uint32:
		val = uint32(tVal)
	case *int32:
		val = uint32(*tVal)
	case *uint32:
		val = uint32(*tVal)
	default:
		return &TypeEncoderError{"int", "TODO"}
	}
	buf := make([]byte, SIZEOF_INT)
	binary.BigEndian.PutUint32(buf, val)
	return writeBuf("int", w, buf)

}

func FormatLong(w io.Writer, value interface{}) error {
	var val uint64

	switch tVal := value.(type) {
	case int64:
		val = uint64(tVal)
	case uint64:
		val = uint64(tVal)
	case *int64:
		val = uint64(*tVal)
	case *uint64:
		val = uint64(*tVal)
	default:
		return &TypeEncoderError{"long", fmt.Sprintf("%s", reflect.TypeOf(tVal))}
	}
	buf := make([]byte, SIZEOF_LONG)
	binary.BigEndian.PutUint64(buf, val)
	return writeBuf("long", w, buf)

}

func FormatFloat(w io.Writer, value interface{}) error {
	var val float32

	switch tVal := value.(type) {
	case float32:
		val = tVal
	case *float32:
		val = *tVal
	default:
		return &TypeEncoderError{"float", "TODO"}
	}
	return FormatInt(w, math.Float32bits(val))
}

func FormatDouble(w io.Writer, value interface{}) error {
	var val float64

	switch tVal := value.(type) {
	case float64:
		val = tVal
	case *float64:
		val = *tVal
	default:
		return &TypeEncoderError{"double", "TODO"}
	}
	return FormatLong(w, math.Float64bits(val))

}

func FormatString(w io.Writer, value interface{}) error {
	var val []byte

	switch tVal := value.(type) {
	case []byte:
		val = tVal
	case *[]byte:
		val = *tVal
	case string:
		val = []byte(tVal)
	case *string:
		val = []byte(*tVal)
	default:
		return &TypeEncoderError{"string", "TODO"}
	}

	if len(val) == 0 {
		err := writeNull(w)
		if err != nil {
			return err
		}
		return nil
	}

	size := len(val)
	err := FormatNullable(w, size)
	if err != nil {
		return err
	}

	err = writeBuf("string", w, val)
	if err != nil {
		return err
	}

	return nil

}

func FormatType(w io.Writer, value interface{}) error {
	var val byte

	switch tVal := value.(type) {
	case int8:
		val = byte(tVal)
	case uint8:
		val = byte(tVal)
	case *int8:
		val = byte(*tVal)
	case *uint8:
		val = byte(*tVal)
	default:
		return &TypeEncoderError{"type", "TODO"}
	}

	if val == NULL_BYTE {
		return writeNull(w)
	}
	return writeBuf("type", w, []byte{val})

}

func FormatBool(w io.Writer, value interface{}) error {
	var val byte

	switch tVal := value.(type) {
	case int:
		val = byte(tVal)
	case *int:
		val = byte(*tVal)
	case bool:
		val = FALSE_BYTE
		if tVal {
			val = TRUE_BYTE
		}
	case *bool:
		val = FALSE_BYTE
		if *tVal {
			val = TRUE_BYTE
		}
	default:
		return &TypeEncoderError{"bool", "TODO"}
	}

	return writeBuf("bool", w, []byte{val})

}

func FormatByte(w io.Writer, value interface{}) error {
	var val byte

	switch tVal := value.(type) {
	case int8:
		val = byte(tVal)
	case uint8:
		val = byte(tVal)
	case *int8:
		val = byte(*tVal)
	case *uint8:
		val = byte(*tVal)
	case int:
		val = byte(tVal)
	case uint:
		val = byte(tVal)
	case *int:
		val = byte(*tVal)
	case *uint:
		val = byte(*tVal)
	case int32:
		val = byte(tVal)
	case uint32:
		val = byte(tVal)
	case *int32:
		val = byte(*tVal)
	case *uint32:
		val = byte(*tVal)
	default:
		return &TypeEncoderError{"byte", "TODO"}
	}

	if val == NULL_BYTE {
		return writeNull(w)
	}
	return writeBuf("byte", w, []byte{val})

}

func FormatSize(w io.Writer, value interface{}) error {

	val, err := castToInt("size", value)
	if err != nil {
		return err
	}

	var b1 int

	msb := val >> MAX_SHIFT
	if msb != 0 {
		b1 = NEXT_MASK
	} else {
		b1 = 0
	}

	err = writeBuf("size", w, []byte{byte(b1 | (val & 0x7F))})
	if err != nil {
		return err
	}
	for val = msb; val > 0; val = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = NEXT_MASK
		} else {
			b1 = 0
		}

		err := writeBuf("size", w, []byte{byte(b1 | (val & 0x7F))})
		if err != nil {
			return err
		}
	}

	return err
}

func FormatNullable(w io.Writer, value interface{}) error {

	val, err := castToInt("nullable", value)
	if err != nil {
		return err
	}

	var b1 int

	msb := val >> NULL_SHIFT
	if msb != 0 {
		b1 = NULL_NEXT_MASK
	} else {
		b1 = 0
	}

	if err = writeBuf("nullable", w, []byte{byte(b1 | (val & 0x7F))}); err != nil {
		return err
	}

	for val = msb; val > 0; val = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = NEXT_MASK
		} else {
			b1 = 0
		}

		if err := writeBuf("null-size", w, []byte{byte(b1 | (val & 0x7F))}); err != nil {
			return err
		}
	}

	return nil
}

func writeNull(w io.Writer) error {
	return writeBuf("write null", w, []byte{0x00})
}

func writeBuf(fnName string, w io.Writer, buf []byte) error {

	_, err := w.Write(buf)
	if err != nil {
		return &EncoderWriteError{fnName, err}
	}

	return nil
}

func castToInt(fnName string, value interface{}) (int, error) {
	var val int

	switch tVal := value.(type) {
	case int:
		val = int(tVal)
	case uint:
		val = int(tVal)
	case *int:
		val = int(*tVal)
	case *uint:
		val = int(*tVal)
	case int32:
		val = int(tVal)
	case uint32:
		val = int(tVal)
	case *int32:
		val = int(*tVal)
	case *uint32:
		val = int(*tVal)
	case int64:
		val = int(tVal)
	case uint64:
		val = int(tVal)
	case *int64:
		val = int(*tVal)
	case *uint64:
		val = int(*tVal)
	default:
		return 0, &TypeEncoderError{fnName, "TODO"}
	}

	return val, nil
}

type EncoderWriteError struct {
	Mame string
	err  error
}

func (e *EncoderWriteError) Error() string {

	return "ras: (FormatFunc " + e.Mame + ") write" + e.err.Error() + ""
}

type TypeEncoderError struct {
	Mame string
	Msg  string
}

func (e *TypeEncoderError) Error() string {
	// if e.Type == nil {
	// 	return "ras: Decode(nil)"
	// }

	// if e.Type.Kind() != reflect.Ptr {
	// 	return "ras: Decode(non-pointer " + e.Type.String() + ")"
	// }
	return "ras: (FormatFunc " + e.Mame + ") " + e.Msg + ""
}
