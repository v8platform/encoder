package ras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"math"
	"reflect"
	"sort"
	"time"
)

var defaultCodecReader = NewReader()

type Decoder struct {
	r            io.Reader
	err          error
	codecVersion int
	PanicOnError bool
}

// NewDecoder create new decoderFunc for version
//
func NewDecoder(r io.Reader, version int) *Decoder {

	return &Decoder{
		r:            r,
		codecVersion: version,
	}

}

func NewDecoderFromBytes(b []byte, version int) *Decoder {

	return NewDecoder(bytes.NewReader(b), version)

}

// An InvalidDecodeError describes an invalid argument passed to Unmarshal.
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

func Decode(data []byte, v interface{}, version int) error {
	return nil

}

func (dec *Decoder) Decode(val interface{}) error {

	if dec.err != nil {
		return dec.err
	}

	if val == nil || reflect.ValueOf(val).Kind() != reflect.Ptr || reflect.ValueOf(val).IsNil() {
		return &InvalidDecodeError{reflect.TypeOf(val)}
	}

	rValue := reflect.ValueOf(val)

	return dec.decodeValue(rValue)

}

func (dec *Decoder) decodeValue(rValue reflect.Value) error {

	if dec.err != nil {
		return dec.err
	}

	rType := rValue.Type()
	if _, ok := rType.(reflect.Type); !ok {
		rType = rType.Elem()
	}

	var err error

	rKind := rType.Kind()

	if rType.Implements(marshalerType) {

		panic("FIXME")

		return err
	}

	switch rKind {
	case reflect.Struct:
		err = dec.decodeStruct(rType, rValue)
	case reflect.Slice:
		err = dec.decodeSlice(rType, rValue)
	case reflect.Array:
		panic("FIXME")
		// err = dec.decodeArray(input, outVal)
	case reflect.Ptr:
		err = dec.decodePtr(rValue)
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
	val := reflect.New(v.Type())
	iFace := val.Interface()

	switch rKind {

	case reflect.String:

		err := decodeString(dec.r, &iFace)
		if err != nil {
			return err
		}

		if v.CanSet() {
			v.SetString(val)
		}

	case reflect.Bool:
		var val bool

		err := decodeBool(dec.r, &val)
		if err != nil {
			return err
		}

		if v.CanSet() {
			v.SetBool(val)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return dec.decodeBasicInt(rKind, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return dec.decodeBasicUint(rKind, v)
	case reflect.Float32:
		var val float64

		err := decodeFloat32(dec.r, &val)
		if err != nil {
			return err
		}

		if v.CanSet() {
			v.SetFloat(val)
		}
	case reflect.Float64:
		var val float64

		err := decodeFloat32(dec.r, &val)
		if err != nil {
			return err
		}

		if v.CanSet() {
			v.SetFloat(val)
		}
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("ras: unsupported type: %s", rKind)
	}
	return nil
}

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}

func (dec *Decoder) decodeBasicInt(kind reflect.Kind, value reflect.Value) error {

	val := reflect.New(value.Type())
	iFace := val.Interface()

	switch kind {
	case reflect.Int, reflect.Int32, reflect.Uint32:
		err := decodeUint32(dec.r, &iFace)
		if err != nil {
			return err
		}
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int64:
	default:
		return &InvalidDecodeError{value.Type()}
	}

	if value.CanSet() {
		value.Set(val)
	}

	return nil
}

func (dec *Decoder) decodeBasicUint(kind reflect.Kind, value reflect.Value) error {
	switch kind {
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	default:
		return &InvalidDecodeError{value.Type()}
	}
	return nil
}

func (dec *Decoder) decodeStruct(rType reflect.Type, rValue reflect.Value) error {

	fields := getCodecFields(rType)

	for _, codecField := range fields {
		if codecField.Ignore {
			continue
		}

		f := rValue.Field(codecField.fieldIdx)

		if codecField.decoder != "" {

			if fn, ok := decoderFunc[codecField.decoder]; ok {

				iFace := f.Addr().Interface()

				err := fn(dec.r, iFace)
				if err != nil {
					return err
				}

				return nil
			} else {
				return &TypeDecodeError{codecField.decoder, "not found decoder func"}
			}

		} else {
			err := dec.decodeValue(f)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (dec *Decoder) decodePtr(value reflect.Value) error {

	valType := value.Type()
	valElemType := valType.Elem()

	if value.CanSet() {
		realVal := reflect.New(valElemType)

		if err := dec.decodeValue(reflect.Indirect(realVal)); err != nil {
			return err
		}

		value.Set(realVal)
	} else {
		if err := dec.decodeValue(reflect.Indirect(value)); err != nil {
			return err
		}
	}
	return nil
}

func (dec *Decoder) decodeSlice(rType reflect.Type, value reflect.Value) error {

	size := dec.decodeSize()

	for i := 0; i < size; i++ {
		elem := reflectAlloc(value.Type().Elem())

		err := dec.decodeValue(elem)
		if err != nil {
			return err
		}

		value.Set(reflect.Append(value, elem))
	}

	return nil
}
