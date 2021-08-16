package ras

import (
	"encoding/binary"
	uuid "github.com/satori/go.uuid"
	"io"
	"math"
	"time"
)

var _ CodecReader = (*reader)(nil)

// NewReader returns a new decoderFunc that reads from r.
//
// The decoderFunc introduces its own buffering and may
// read data from r beyond the JSON values requested.
func NewReader() CodecReader {
	return &reader{}
}

type reader struct {
	PanicOnError bool
}

func (e *reader) EndpointIdPtr(ptr *int, r io.Reader) {
	*ptr = e.EndpointId(r)
}

func (e *reader) EndpointId(r io.Reader) int {
	id := e.NullableSize(r)
	return id
}

func (e *reader) TimePtr(ptr *time.Time, r io.Reader) {
	*ptr = e.Time(r)
}

func (e *reader) Time(r io.Reader) time.Time {

	ticks := e.Long(r)
	return dateFromTicks(ticks)

}

func (e *reader) BoolPtr(ptr *bool, r io.Reader) {

	val, ok := e.Bool(r)

	if ok {
		*ptr = val
	}

}

func (e *reader) Bool(r io.Reader) (bool, bool) {
	b := e.readByte("Bool", r)

	switch b {

	case TRUE_BYTE:
		return true, true
	case FALSE_BYTE:
		return false, true
	}

	return false, false
}

func (e *reader) BytePtr(ptr *byte, r io.Reader) {

	*ptr = e.Byte(r)

}

func (e *reader) Byte(r io.Reader) byte {
	b := e.readByte("Byte", r)
	return b
}

func (e *reader) CharPtr(ptr *int16, r io.Reader) {
	e.Int16Ptr(ptr, r)
}

func (e *reader) Char(r io.Reader) int16 {
	return e.Int16(r)
}

func (e *reader) ShortPtr(ptr *int16, r io.Reader) {
	e.Int16Ptr(ptr, r)
}

func (e *reader) Short(r io.Reader) int16 {
	return e.Int16(r)
}

func (e *reader) IntPtr(ptr *int, r io.Reader) {
	*ptr = e.Int(r)
}

func (e *reader) Int(r io.Reader) int {
	return int(e.Uint32(r))
}

func (e *reader) Uint16Ptr(ptr *uint16, r io.Reader) {
	*ptr = e.Uint16(r)
}

func (e *reader) Uint16(r io.Reader) uint16 {

	buf := make([]byte, 2)
	e.read("Uint16", r, buf)

	val := binary.BigEndian.Uint16(buf)
	return val
}

func (e *reader) Int16Ptr(ptr *int16, r io.Reader) {
	*ptr = e.Int16(r)
}

func (e *reader) Int16(r io.Reader) int16 {

	buf := make([]byte, 2)
	e.read("Int16", r, buf)
	//buf = buf[:n]

	val := int16(binary.BigEndian.Uint16(buf))
	return val
}

func (e *reader) UintPtr(ptr *uint, r io.Reader) {
	*ptr = e.Uint(r)
}

func (e *reader) Uint(r io.Reader) uint {
	return uint(e.Uint32(r))
}

func (e *reader) Int32Ptr(ptr *int32, r io.Reader) {
	*ptr = e.Int32(r)
}

func (e *reader) Int32(r io.Reader) int32 {
	return int32(e.Uint32(r))
}

func (e *reader) Uint32Ptr(ptr *uint32, r io.Reader) {
	*ptr = e.Uint32(r)
}

func (e *reader) Uint32(r io.Reader) uint32 {

	buf := make([]byte, 4)
	e.read("Uint32", r, buf)
	val := binary.BigEndian.Uint32(buf)
	return val
}

func (e *reader) Int64Ptr(ptr *int64, r io.Reader) {
	*ptr = e.Int64(r)
}

func (e *reader) Int64(r io.Reader) int64 {
	return int64(e.Uint64(r))
}

func (e *reader) Uint64Ptr(ptr *uint64, r io.Reader) {
	*ptr = e.Uint64(r)
}

func (e *reader) Uint64(r io.Reader) uint64 {
	buf := make([]byte, 8)
	e.read("Uint64", r, buf)
	val := binary.BigEndian.Uint64(buf)
	return val
}

func (e *reader) LongPtr(ptr *int64, r io.Reader) {
	*ptr = e.Long(r)
}

func (e *reader) Long(r io.Reader) int64 {
	return e.Int64(r)
}

func (e *reader) Float32Ptr(ptr *float32, r io.Reader) {
	*ptr = e.Float32(r)
}

func (e *reader) Float32(r io.Reader) float32 {
	b := e.Uint32(r)

	return math.Float32frombits(b)
}

func (e *reader) Float64Ptr(ptr *float64, r io.Reader) {
	*ptr = e.Float64(r)
}

func (e *reader) Float64(r io.Reader) float64 {
	b := e.Uint64(r)
	return math.Float64frombits(b)
}

func (e *reader) DoublePtr(ptr *float64, r io.Reader) {
	*ptr = e.Double(r)
}

func (e *reader) Double(r io.Reader) float64 {
	return e.Float64(r)
}

func (e *reader) Null(_ io.Reader) {
	panic("implement me")
}

func (e *reader) StringPtr(ptr *string, r io.Reader) {
	*ptr = e.String(r)
}

func (e *reader) String(r io.Reader) string {
	size := e.NullableSize(r)
	if size == 0 {
		return ""
	}
	buf := make([]byte, size)
	e.read("String", r, buf)

	return string(buf)
}

func (e *reader) TypedValue(_ interface{}, _ io.Reader) {
	panic("implement me")
}

func (e *reader) UuidPtr(ptr *uuid.UUID, r io.Reader) {
	*ptr = e.Uuid(r)
}

func (e *reader) Uuid(r io.Reader) uuid.UUID {

	buf := make([]byte, 16)
	e.read("Uuid", r, buf)
	u, _ := uuid.FromBytes(buf)

	return u
}

func (e *reader) Size(r io.Reader) int {

	ff := 0xFFFFFF80
	b1 := e.readByte("Size", r)
	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := MAX_SHIFT; (cur & ff) != 0x0; {

		b1 = e.readByte("Size", r)
		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += MAX_SHIFT
	}

	return size
}

func (e *reader) NullableSize(r io.Reader) int {
	size := 0
	//ff := 0xFFFFFF80
	b1 := e.readByte("NullableSize", r)
	cur := int(b1 & 0xFF)
	if (cur & 0xFFFFFF80) == 0x0 {
		size = cur & 0x3F
		if cur&0x40 == 0x0 {
			return size
		}

		shift := NULL_SHIFT
		b1 := e.readByte("NullableSize", r)
		cur := int(b1 & 0xFF)
		size += (cur & 0x7F) << NULL_SHIFT
		shift += MAX_SHIFT

		for (cur & 0xFFFFFF80) != 0x0 {

			b1 := e.readByte("NullableSize", r)
			cur = int(b1 & 0xFF)
			size += (cur & 0x7F) << shift
			shift += MAX_SHIFT

		}
		return size
	}

	if (cur & 0x7F) != 0x0 {
		panic("null expected")
	}

	return size
}

func (e *reader) Type(r io.Reader) byte {
	b1 := e.readByte("Type", r)
	cur := b1 & 0xFF

	return cur
}

func (e *reader) Bytes(ptr []byte, r io.Reader) {

	e.read("Bytes", r, ptr)

}

func (e *reader) Value(ptr interface{}, r io.Reader) {

	switch typed := ptr.(type) {

	case *bool:
		e.BoolPtr(typed, r)
	case *float32:
		e.Float32Ptr(typed, r)
	case *float64:
		e.Float64Ptr(typed, r)
	case *int:
		e.IntPtr(typed, r)
	case *uint:
		e.UintPtr(typed, r)
	case *int16:
		e.Int16Ptr(typed, r)
	case *uint16:
		e.Uint16Ptr(typed, r)
	case *int32:
		e.Int32Ptr(typed, r)
	case *uint32:
		e.Uint32Ptr(typed, r)
	case *int64:
		e.Int64Ptr(typed, r)
	case *uint64:
		e.Uint64Ptr(typed, r)
	case *string:
		e.StringPtr(typed, r)
	case *byte:
		e.BytePtr(typed, r)
	default:
		panic("error decode value")
	}
}

func (e *reader) panicOnError(fnName string, p []byte, n int, err error) {

	if err != nil && e.PanicOnError {
		panic(&DecoderError{
			fn:        fnName,
			needBytes: len(p),
			readBytes: n,
			err:       err,
		})
	}
}

func (e *reader) read(fnName string, r io.Reader, p []byte) {
	n, err := r.Read(p)
	e.panicOnError(fnName, p, n, err)
}

func (e *reader) readByte(fnName string, r io.Reader) byte {

	buf := make([]byte, 1)
	e.read(fnName, r, buf)
	return buf[0]
}
