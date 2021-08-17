package ras

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"time"
)

var (
	formatterType  = reflect.TypeOf((*Formatter)(nil)).Elem()
	marshallerType = reflect.TypeOf((*Marshaller)(nil)).Elem()
)

type Formatter interface {
	FormatRAS(version int) ([]byte, error)
}

type Marshaller interface {
	MarshalRAS(writer io.Writer, version int) (int, error)
}

type Encoder struct {
	writer io.Writer
	err    error
}

// NewDecoder create new encoderFunc for version
//
func NewEncoder(r io.Writer) *Encoder {

	return &Encoder{
		writer: r,
	}

}

// An InvalidEncodeError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidEncodeError struct {
	Type reflect.Type
}

func (e *InvalidEncodeError) Error() string {
	if e.Type == nil {
		return "ras: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "ras: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "ras: Encode(nil " + e.Type.String() + ")"
}

func Encode(v interface{}, version int) ([]byte, error) {

	buf := bytes.NewBuffer([]byte{})

	encoder := NewEncoder(buf)
	err := encoder.Encode(v, version)
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil

}

func (dec *Encoder) Encode(val interface{}, version int) error {

	if dec.err != nil {
		return dec.err
	}

	if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		return &InvalidEncodeError{reflect.TypeOf(val)}
	}

	rValue := reflect.ValueOf(val)

	return dec.encode(rValue, version)

}

func (dec *Encoder) encode(rValue reflect.Value, version int) error {

	var err error

	if dec.err != nil {
		return dec.err
	}

	rType := rValue.Type()
	if _, ok := rType.(reflect.Type); !ok {
		rType = rType.Elem()
	}

	if rValue.CanAddr() {
		iFace := rValue.Addr().Interface()

		switch iFace.(type) {
		case *time.Time, time.Time:
			_, err := encodeTime(dec.writer, iFace)
			if err != nil {
				return err
			}

			return nil
		}
	}

	rKind := rType.Kind()

	if rType.Implements(marshallerType) {

		panic("FIXME")

		return err
	}

	switch rKind {
	case reflect.Struct:
		err = dec.encodeStruct(rType, rValue, version)
	case reflect.Slice:
		err = dec.encodeSlice(rValue, version)
	case reflect.Ptr:
		err = dec.encodePtr(rValue, version)
	default:
		err = dec.encodeBasic(rType, rValue)
	}
	return err
}

func (dec *Encoder) decodeCustom(v reflect.Value, decodeFn func() interface{}) error {

	value := decodeFn()
	if v.CanSet() {
		v.Set(reflect.ValueOf(value))
	}

	return nil
}

func (dec *Encoder) encodeBasic(rType reflect.Type, v reflect.Value) error {

	rKind := rType.Kind()

	iFace := v.Interface()

	switch rKind {

	case reflect.String:
		_, err := encodeString(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Bool:
		_, err := encodeBool(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32:
		_, err := encodeUint32(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Int16, reflect.Uint16:
		_, err := encodeUint16(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Int64, reflect.Uint64:
		_, err := encodeUint64(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Int8, reflect.Uint8:
		_, err := encodeByte(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Float32:
		_, err := encodeFloat32(dec.writer, iFace)
		if err != nil {
			return err
		}
	case reflect.Float64:
		_, err := encodeFloat64(dec.writer, iFace)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("ras: unsupported type: %s", rKind)
	}
	return nil
}

func (dec *Encoder) encodeStruct(rType reflect.Type, rValue reflect.Value, version int) error {

	fields := getCodecFields(rType)

	for _, codecField := range fields {
		if codecField.Ignore {
			continue
		}

		if codecField.Version > version {
			continue
		}

		f := rValue.Field(codecField.fieldIdx)

		if codecField.codec != "" {

			if fn, ok := encoderFunc[codecField.codec]; ok {

				iFace := f.Interface()

				_, err := fn(dec.writer, iFace)
				if err != nil {
					return err
				}
				continue
			}

			return &TypeDecodeError{codecField.codec, "not found codec func"}

		}
		err := dec.encode(f, version)
		if err != nil {
			return err
		}

	}

	return nil
}

func (dec *Encoder) encodePtr(value reflect.Value, version int) error {

	elem := value.Elem()
	if err := dec.encode(elem, version); err != nil {
		return err
	}

	return nil
}

func (dec *Encoder) encodeSlice(value reflect.Value, version int) error {

	var size int

	size = value.Len()

	_, err := encodeSize(dec.writer, size)
	if err != nil {
		return err
	}

	for i := 0; i < size; i++ {

		elem := value.Index(i)

		err := dec.encode(elem, version)
		if err != nil {
			return err
		}

	}

	return nil
}
