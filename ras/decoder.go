package ras

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"time"
)

type Decoder struct {
	r   io.Reader
	err error
}

// NewDecoder create new encoderFunc for version
//
func NewDecoder(r io.Reader) *Decoder {

	return &Decoder{
		r: r,
	}

}

func NewDecoderFromBytes(b []byte) *Decoder {

	return NewDecoder(bytes.NewReader(b))

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

func Decode(data []byte, v interface{}, version int) error {

	decoder := NewDecoderFromBytes(data)

	return decoder.Decode(v, version)

}

func (dec *Decoder) Decode(val interface{}, version int) error {

	if dec.err != nil {
		return dec.err
	}

	if val == nil || reflect.ValueOf(val).Kind() != reflect.Ptr || reflect.ValueOf(val).IsNil() {
		return &InvalidDecodeError{reflect.TypeOf(val)}
	}

	rValue := reflect.ValueOf(val)

	return dec.decodeValue(rValue, version)

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
		case *time.Time:
			err := decodeTime(dec.r, iFace)
			if err != nil {
				return err
			}

			return nil
		}
	}

	rKind := rType.Kind()

	if rType.Implements(marshalerType) {

		panic("FIXME")

		return err
	}

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

		err := decodeString(dec.r, iFace)
		if err != nil {
			return err
		}

	case reflect.Bool:

		err := decodeBool(dec.r, iFace)
		if err != nil {
			return err
		}

	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32:
		err := decodeUint32(dec.r, iFace)
		if err != nil {
			return err
		}
	case reflect.Int16, reflect.Uint16:
		err := decodeUint16(dec.r, iFace)
		if err != nil {
			return err
		}
	case reflect.Int64, reflect.Uint64:
		err := decodeUint64(dec.r, iFace)
		if err != nil {
			return err
		}
	case reflect.Int8, reflect.Uint8:
		err := decodeByte(dec.r, iFace)
		if err != nil {
			return err
		}
	case reflect.Float32:

		err := decodeFloat32(dec.r, iFace)
		if err != nil {
			return err
		}

	case reflect.Float64:
		err := decodeFloat32(dec.r, iFace)
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

				err := typeDecoderFunc(dec.r, iFace)
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

	valType := value.Type()
	valElemType := valType.Elem()

	if value.CanSet() {
		realVal := reflect.New(valElemType)

		if err := dec.decodeValue(reflect.Indirect(realVal), version); err != nil {
			return err
		}

		value.Set(realVal)
	} else {
		if err := dec.decodeValue(reflect.Indirect(value), version); err != nil {
			return err
		}
	}
	return nil
}

func (dec *Decoder) decodeSlice(value reflect.Value, version int) error {

	var size int
	err := decodeSize(dec.r, &size)
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

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}
