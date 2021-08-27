package ras

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

var _ Codec = (*codec)(nil)
var _ CodecWriter = (*codec)(nil)
var _ CodecReader = (*codec)(nil)

func NewCodecWriter() CodecWriter {
	return &codec{}
}
func NewCodecReader() CodecReader {
	return &codec{}
}

func NewCodec() Codec {
	return &codec{}
}

//goland:noinspection ALL
type codec struct{}

func (c *codec) ReadBoolPtr(val *bool, reader io.Reader) (n int, err error) {
	return decodeBool(reader, val)
}

func (c *codec) ReadBool(reader io.Reader) (val bool, n int, err error) {
	n, err = c.ReadBoolPtr(&val, reader)
	return
}

func (c *codec) ReadBytePtr(val *byte, reader io.Reader) (n int, err error) {
	return decodeByte(reader, val)
}

func (c *codec) ReadByte(reader io.Reader) (val byte, n int, err error) {
	n, err = c.ReadBytePtr(&val, reader)
	return
}

func (c *codec) ReadIntPtr(val *int, reader io.Reader) (n int, err error) {
	return decodeUint32(reader, val)
}

func (c *codec) ReadInt(reader io.Reader) (val int, n int, err error) {
	n, err = c.ReadIntPtr(&val, reader)
	return
}

func (c *codec) ReadUintPtr(val *uint, reader io.Reader) (n int, err error) {
	return decodeUint32(reader, val)
}

func (c *codec) ReadUint(reader io.Reader) (val uint, n int, err error) {
	n, err = c.ReadUintPtr(&val, reader)
	return
}
func (c *codec) ReadUint16Ptr(val *uint16, reader io.Reader) (n int, err error) {
	return decodeUint16(reader, val)
}

func (c *codec) ReadUint16(reader io.Reader) (val uint16, n int, err error) {
	n, err = c.ReadUint16Ptr(&val, reader)
	return
}

func (c *codec) ReadInt32Ptr(val *int32, reader io.Reader) (n int, err error) {
	return decodeUint32(reader, val)
}

func (c *codec) ReadInt32(reader io.Reader) (val int32, n int, err error) {
	n, err = c.ReadInt32Ptr(&val, reader)
	return
}

func (c *codec) ReadUint32Ptr(val *uint32, reader io.Reader) (n int, err error) {
	return decodeUint32(reader, val)
}

func (c *codec) ReadUint32(reader io.Reader) (val uint32, n int, err error) {
	n, err = c.ReadUint32Ptr(&val, reader)
	return
}

func (c *codec) ReadInt64Ptr(val *int64, reader io.Reader) (n int, err error) {
	return decodeUint64(reader, val)
}

func (c *codec) ReadInt64(reader io.Reader) (val int64, n int, err error) {
	n, err = c.ReadInt64Ptr(&val, reader)
	return
}

func (c *codec) ReadUint64Ptr(val *uint64, reader io.Reader) (n int, err error) {
	return decodeUint64(reader, val)
}

func (c *codec) ReadUint64(reader io.Reader) (val uint64, n int, err error) {
	n, err = c.ReadUint64Ptr(&val, reader)
	return
}

func (c *codec) ReadFloat32Ptr(val *float32, reader io.Reader) (n int, err error) {
	return decodeFloat32(reader, val)
}

func (c *codec) ReadFloat32(reader io.Reader) (val float32, n int, err error) {
	n, err = c.ReadFloat32Ptr(&val, reader)
	return
}

func (c *codec) ReadFloat64Ptr(val *float64, reader io.Reader) (n int, err error) {
	return decodeFloat64(reader, val)
}

func (c *codec) ReadFloat64(reader io.Reader) (val float64, n int, err error) {
	n, err = c.ReadFloat64Ptr(&val, reader)
	return
}

func (c *codec) ReadStringPtr(val *string, reader io.Reader) (n int, err error) {
	return decodeString(reader, val)
}

func (c *codec) ReadString(reader io.Reader) (val string, n int, err error) {
	n, err = c.ReadStringPtr(&val, reader)
	return
}

func (c *codec) ReadUuidPtr(val interface{}, reader io.Reader) (n int, err error) {
	return decodeUUID(reader, val)
}

func (c *codec) ReadUuid(reader io.Reader) (val uuid.UUID, n int, err error) {
	n, err = c.ReadUuidPtr(&val, reader)
	return
}
func (c *codec) ReadSizePtr(val interface{}, reader io.Reader) (n int, err error) {
	return decodeSize(reader, val)
}

func (c *codec) ReadSize(reader io.Reader) (val int, n int, err error) {
	n, err = c.ReadUuidPtr(&val, reader)
	return
}

func (c *codec) ReadNullableSizePtr(val interface{}, reader io.Reader) (n int, err error) {
	return decodeNullableSize(reader, val)
}
func (c *codec) ReadNullableSize(reader io.Reader) (val int, n int, err error) {
	n, err = c.ReadNullableSizePtr(&val, reader)
	return
}

func (c *codec) ReadTypePtr(val *byte, reader io.Reader) (n int, err error) {
	return decodeType(reader, val)
}

func (c *codec) ReadType(reader io.Reader) (val byte, n int, err error) {
	n, err = c.ReadTypePtr(&val, reader)
	return
}

func (c *codec) ReadTimePtr(val interface{}, reader io.Reader) (n int, err error) {
	return decodeTime(reader, val)
}

func (c *codec) ReadTime(reader io.Reader) (val time.Time, n int, err error) {
	n, err = c.ReadTimePtr(&val, reader)
	return
}

func (c *codec) WriteBool(val bool, writer io.Writer) (n int, err error) {
	return encodeBool(writer, val)
}

func (c *codec) WriteByte(val byte, writer io.Writer) (n int, err error) {
	return encodeByte(writer, val)
}

func (c *codec) WriteInt(val int, writer io.Writer) (n int, err error) {
	return encodeUint32(writer, val)
}

func (c *codec) WriteUint(val uint, writer io.Writer) (n int, err error) {
	return encodeUint32(writer, val)
}

func (c *codec) WriteInt16(val int16, writer io.Writer) (n int, err error) {
	return encodeUint16(writer, val)
}

func (c *codec) WriteUint16(val uint16, writer io.Writer) (n int, err error) {
	return encodeUint16(writer, val)
}

func (c *codec) WriteInt32(val int32, writer io.Writer) (n int, err error) {
	return encodeUint32(writer, val)
}

func (c *codec) WriteUint32(val uint32, writer io.Writer) (n int, err error) {
	return encodeUint32(writer, val)
}

func (c *codec) WriteInt64(val int64, writer io.Writer) (n int, err error) {
	return encodeUint64(writer, val)
}

func (c *codec) WriteUint64(val uint64, writer io.Writer) (n int, err error) {
	return encodeUint64(writer, val)
}

func (c *codec) WriteFloat32(val float32, writer io.Writer) (n int, err error) {
	return encodeFloat32(writer, val)
}

func (c *codec) WriteFloat64(val float64, writer io.Writer) (n int, err error) {
	return encodeFloat64(writer, val)
}

func (c *codec) WriteNull(writer io.Writer) (n int, err error) {
	return writeNull(writer)
}

func (c *codec) WriteString(val string, writer io.Writer) (n int, err error) {
	return encodeString(writer, val)
}

func (c *codec) WriteUuid(val interface{}, writer io.Writer) (n int, err error) {
	return EncodeUuid(writer, val)
}

func (c *codec) WriteSize(val int, writer io.Writer) (n int, err error) {
	return encodeSize(writer, val)
}

func (c *codec) WriteNullableSize(val int, writer io.Writer) (n int, err error) {
	return encodeNullableSize(writer, val)
}

func (c *codec) WriteType(val byte, writer io.Writer) (n int, err error) {
	return encodeType(writer, val)
}

func (c *codec) WriteTime(val interface{}, writer io.Writer) (n int, err error) {
	return encodeTime(writer, val)
}

func (c *codec) Version() int {
	return 256 // Версия кодека у 1С
}

type Codec interface {
	CodecWriter
	CodecReader
	Version() int
}

type CodecWriter interface {
	WriteBool(val bool, writer io.Writer) (n int, err error)
	WriteByte(val byte, writer io.Writer) (n int, err error)
	WriteInt(val int, writer io.Writer) (n int, err error)
	WriteUint(val uint, writer io.Writer) (n int, err error)
	WriteInt16(val int16, writer io.Writer) (n int, err error)
	WriteUint16(val uint16, writer io.Writer) (n int, err error)
	WriteInt32(val int32, writer io.Writer) (n int, err error)
	WriteUint32(val uint32, writer io.Writer) (n int, err error)
	WriteInt64(val int64, writer io.Writer) (n int, err error)
	WriteUint64(val uint64, writer io.Writer) (n int, err error)
	WriteFloat32(val float32, writer io.Writer) (n int, err error)
	WriteFloat64(val float64, writer io.Writer) (n int, err error)

	WriteNull(writer io.Writer) (n int, err error)
	WriteString(val string, writer io.Writer) (n int, err error)

	WriteUuid(val interface{}, writer io.Writer) (n int, err error)
	WriteSize(val int, writer io.Writer) (n int, err error)
	WriteNullableSize(val int, writer io.Writer) (n int, err error)
	WriteType(val byte, writer io.Writer) (n int, err error)
	WriteTime(val interface{}, writer io.Writer) (n int, err error)
}

type CodecReader interface {
	ReadBoolPtr(val *bool, reader io.Reader) (n int, err error)
	ReadBool(reader io.Reader) (val bool, n int, err error)

	ReadBytePtr(val *byte, reader io.Reader) (n int, err error)
	ReadByte(reader io.Reader) (val byte, n int, err error)

	ReadIntPtr(val *int, reader io.Reader) (n int, err error)
	ReadInt(reader io.Reader) (val int, n int, err error)

	ReadUintPtr(val *uint, reader io.Reader) (n int, err error)
	ReadUint(reader io.Reader) (val uint, n int, err error)

	ReadUint16(reader io.Reader) (val uint16, n int, err error)
	ReadUint16Ptr(ptr *uint16, reader io.Reader) (n int, err error)

	ReadInt32Ptr(val *int32, reader io.Reader) (n int, err error)
	ReadInt32(reader io.Reader) (val int32, n int, err error)

	ReadUint32Ptr(val *uint32, reader io.Reader) (n int, err error)
	ReadUint32(reader io.Reader) (val uint32, n int, err error)

	ReadInt64Ptr(val *int64, reader io.Reader) (n int, err error)
	ReadInt64(reader io.Reader) (val int64, n int, err error)

	ReadUint64Ptr(val *uint64, reader io.Reader) (n int, err error)
	ReadUint64(reader io.Reader) (val uint64, n int, err error)

	ReadFloat32Ptr(val *float32, reader io.Reader) (n int, err error)
	ReadFloat32(reader io.Reader) (val float32, n int, err error)

	ReadFloat64Ptr(val *float64, reader io.Reader) (n int, err error)
	ReadFloat64(reader io.Reader) (val float64, n int, err error)

	ReadStringPtr(val *string, reader io.Reader) (n int, err error)
	ReadString(reader io.Reader) (val string, n int, err error)

	ReadUuidPtr(val interface{}, reader io.Reader) (n int, err error)
	ReadUuid(reader io.Reader) (val uuid.UUID, n int, err error)

	ReadSizePtr(val interface{}, reader io.Reader) (n int, err error)
	ReadSize(reader io.Reader) (val int, n int, err error)

	ReadNullableSizePtr(val interface{}, reader io.Reader) (n int, err error)
	ReadNullableSize(reader io.Reader) (int, n int, err error)

	ReadTypePtr(val *byte, reader io.Reader) (n int, err error)
	ReadType(reader io.Reader) (val byte, n int, err error)

	ReadTimePtr(ptr interface{}, reader io.Reader) (n int, err error)
	ReadTime(reader io.Reader) (val time.Time, n int, err error)
}
