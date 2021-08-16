package ras

type TypeInterface byte

//goland:noinspection ALL
const (
	BOOLEAN       TypeInterface = 1
	BYTE                        = 2
	SHORT                       = 3
	INT                         = 4
	LONG                        = 5
	FLOAT                       = 6
	DOUBLE                      = 7
	SIZE                        = 8
	NULLABLE_SIZE               = 9
	STRING                      = 10
	UUID                        = 11
	TYPE                        = 12
	ENDPOINT_ID                 = 13
)

func (t TypeInterface) raw() byte {
	return byte(t)
}
func (t TypeInterface) Type() byte {
	return byte(t)
}

type CodecKind int8 // keep exact type opaque as the int type may change

// Constants as defined by the google.protobuf.CodecField.Kind enumeration.
const (
	TypeKind CodecKind = iota
	BoolKind
	Int32Kind
	Uint32Kind
	Int64Kind
	Uint64Kind
	Fixed32Kind
	FloatKind
	Fixed64Kind
	DoubleKind
	StringKind
	BytesKind
	MessageKind
	TimestampKind
)
