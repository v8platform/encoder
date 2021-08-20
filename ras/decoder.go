package ras

import (
	"bytes"
	"fmt"
	pb "google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"reflect"
	"time"
)

var (
	parserType      = reflect.TypeOf((*Parser)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

type Parser interface {
	ParseRAS(data []byte, version int) (n int, err error)
}

type Unmarshaler interface {
	UnmarshalRAS(reader io.Reader, version int) (n int, err error)
}

type Decoder struct {
	buf *bytes.Buffer
	err error
	n   int // bytes decoded
}

// NewDecoder create new encoderFunc for version
//
func NewDecoderFromReader(r io.Reader) *Decoder {

	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil
	}

	return &Decoder{
		buf: buf,
	}

}

func NewDecoder(b []byte) *Decoder {

	return &Decoder{
		buf: bytes.NewBuffer(b),
	}

}

// An InvalidEncodeError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidDecodeError struct {
	Type reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.Type == nil {
		return "ras: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "ras: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "ras: Decode(nil " + e.Type.String() + ")"
}

func Decode(data []byte, v interface{}, version int) (int, error) {

	decoder := NewDecoder(data)

	return decoder.Decode(v, version)

}

func (dec *Decoder) Decode(val interface{}, version int) (int, error) {

	dec.n = 0

	if dec.err != nil {
		return dec.n, dec.err
	}

	if val == nil || reflect.ValueOf(val).Kind() != reflect.Ptr || reflect.ValueOf(val).IsNil() {
		return dec.n, &InvalidDecodeError{reflect.TypeOf(val)}
	}

	rValue := reflect.ValueOf(val)

	return dec.n, dec.decodeValue(rValue, version)

}

func (dec *Decoder) decodeValue(rValue reflect.Value, version int) error {

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
		case *time.Time, *pb.Timestamp:
			n, err := decodeTime(dec.buf, iFace)
			dec.n += n
			if err != nil {
				return err
			}

			return nil
		}
	}

	rKind := rType.Kind()

	switch rKind {
	case reflect.Struct:
		err = dec.decodeStruct(rType, rValue, version)
	case reflect.Slice:
		err = dec.decodeSlice(rValue, version)
	case reflect.Ptr:
		err = dec.decodePtr(rValue, version)
	default:
		err = dec.decodeBasic(rType, rValue)
	}
	return err
}

func (dec *Decoder) decodeCustom(v reflect.Value, decodeFn func() interface{}) error {

	value := decodeFn()
	if v.CanSet() {
		v.Set(reflect.ValueOf(value))
	}

	return nil
}

func (dec *Decoder) decodeBasic(rType reflect.Type, v reflect.Value) error {

	rKind := rType.Kind()

	if !v.CanAddr() {
		return fmt.Errorf("ras: cannot addr for value: %s", v.String())
	}

	iFace := v.Addr().Interface()

	switch rKind {

	case reflect.String:

		n, err := decodeString(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}

	case reflect.Bool:

		n, err := decodeBool(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}

	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32:
		n, err := decodeUint32(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}
	case reflect.Int16, reflect.Uint16:
		n, err := decodeUint16(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}
	case reflect.Int64, reflect.Uint64:
		n, err := decodeUint64(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}
	case reflect.Int8, reflect.Uint8:
		n, err := decodeByte(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}
	case reflect.Float32:

		n, err := decodeFloat32(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}

	case reflect.Float64:
		n, err := decodeFloat32(dec.buf, iFace)
		dec.n += n
		if err != nil {
			return err
		}

	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("ras: unsupported type: %s", rKind)

	}
	return nil
}

func (dec *Decoder) decodeStruct(rType reflect.Type, rValue reflect.Value, version int) error {

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

			if typeDecoderFunc, ok := decoderFunc[codecField.codec]; ok {
				var iFace interface{}

				if f.Kind() == reflect.Ptr {
					valType := f.Type()
					valElemType := valType.Elem()
					val := reflect.New(valElemType)
					iFace = val.Interface()
				} else {
					iFace = f.Addr().Interface()
				}

				n, err := typeDecoderFunc(dec.buf, iFace)
				dec.n += n
				if err != nil {
					return err
				}

				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(iFace))
				}

				continue
			}

			return &TypeDecodeError{codecField.codec, "not found codec func"}

		}

		err := dec.decodeValue(f, version)
		if err != nil {
			return err
		}

	}

	return nil
}

func (dec *Decoder) decodePtr(value reflect.Value, version int) error {

	un, _, rValue := indirect(value, false)

	if un != nil {

		n, err := un.UnmarshalRAS(dec.buf, version)
		dec.n += n
		if err != nil {
			return err
		}
		return nil

	}

	//
	// valType := value.Type()
	// valElemType := valType.Elem()
	//
	// if value.CanSet() {
	// 	realVal := reflect.New(valElemType)
	//
	// 	if err := dec.decodeValue(reflect.Indirect(realVal), version); err != nil {
	// 		return err
	// 	}
	//
	// 	value.Set(realVal)
	// } else {
	// 	if err := dec.decodeValue(reflect.Indirect(value), version); err != nil {
	// 		return err
	// 	}
	// }
	return dec.decodeValue(rValue, version)
}

func (dec *Decoder) decodeSlice(value reflect.Value, version int) error {

	var size int
	n, err := decodeSize(dec.buf, &size)
	dec.n += n
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		elem := reflectAlloc(value.Type().Elem())

		err := dec.decodeValue(elem, version)
		if err != nil {
			return err
		}

		value.Set(reflect.Append(value, elem))
	}

	return nil
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// If it encounters an Unmarshaler, indirect stops and returns that.
// If decodingNull is true, indirect stops at the first settable pointer so it
// can be set to nil.
func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, Parser, reflect.Value) {

	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if !decodingNull {
				if u, ok := v.Interface().(Parser); ok {
					return nil, u, reflect.Value{}
				}
			}
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, nil, v
}

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}
