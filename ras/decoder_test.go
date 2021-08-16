package ras

import (
	"bytes"
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
	pb "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
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
			dec := NewDecoderFromBytes(tt.buf)
			if err := dec.Decode(tt.value, 1); (err != nil) != tt.wantErr {
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
	codec.Uuid(uuid.NewV1(), buf)
	codec.Int32(1, buf)
	codec.String("Блокировка 1", buf)
	codec.Uuid(uuid.NewV1(), buf)
	codec.Int32(2, buf)
	codec.String("Блокировка 2", buf)
	codec.Time(time.Now().AddDate(0, -1, 9), buf)

	return buf.Bytes()
}

type Message struct {
	Type  int           `rac:",1"`
	Kind  *int64        `rac:"int64,2"`
	Locks []*Lock       `rac:",3"`
	Time  *pb.Timestamp `rac:"time,4"`
}

type Lock struct {
	UUID string `rac:"uuid,1"`
	ID   int    `rac:",2"`
	Msg  string `rac:",3"`
	// Time time.Time `rac:"time,3"`
}
