package ras

import (
	"bytes"
	"github.com/k0kubun/pp"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {

	tests := []struct {
		name    string
		buf     []byte
		value   interface{}
		wantErr bool
	}{
		{
			"Test decode int",
			getTestData(),
			&Message{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := NewDecoderFromBytes(tt.buf, 1)
			if err := dec.Decode(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			pp.Println(tt.value)

		})
	}
}

func getTestData() []byte {

	codec := &encoder{}
	buf := bytes.NewBuffer([]byte{})
	codec.Int32(111, buf)
	codec.Uint64(222, buf)
	codec.Size(2, buf)
	codec.Int32(1, buf)
	codec.String("Блокировка 1", buf)
	// codec.Time(time.Now().AddDate(0, -1, 9), buf)
	codec.Int32(2, buf)
	codec.String("Блокировка 2", buf)
	// codec.Time(time.Now().AddDate(0, 0, 30), buf)

	return buf.Bytes()
}

type Message struct {
	Type  int     `rac:",1"`
	Kind  int64   `rac:"uint64,2"`
	Locks []*Lock `rac:",3"`
}

type Lock struct {
	ID  int    `rac:",1"`
	Msg string `rac:",2"`
	// Time time.Time `rac:"time,3"`
}
