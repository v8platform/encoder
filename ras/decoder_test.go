package ras

import (
	"bytes"
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
	pb "google.golang.org/protobuf/types/known/timestamppb"
	"io"
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
			dec := NewDecoder(tt.buf)
			if _, err := dec.Decode(tt.value, 1); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			pp.Println(tt.value)

		})
	}
}

func getTestData() []byte {

	codec := NewCodec()
	buf := bytes.NewBuffer([]byte{})
	codec.WriteInt32(111, buf)
	codec.WriteUint64(222, buf)
	codec.WriteSize(2, buf)
	codec.WriteUuid(uuid.NewV1(), buf)
	codec.WriteInt32(1, buf)
	codec.WriteString("Блокировка 1", buf)
	codec.WriteUuid(uuid.NewV1(), buf)
	codec.WriteInt32(2, buf)
	codec.WriteString("Блокировка 2", buf)
	codec.WriteTime(time.Now().AddDate(0, -1, 9), buf)

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

func (l *Lock) UnmarshalRAS(reader io.Reader, version int) (n int, err error) {

	c := NewCodecReader()

	n, err = c.ReadUuidPtr(&l.UUID, reader)
	if err != nil {
		return n, err
	}

	n, err = c.ReadIntPtr(&l.ID, reader)
	if err != nil {
		return n, err
	}

	n, err = c.ReadStringPtr(&l.Msg, reader)
	if err != nil {
		return n, err
	}

	return n, nil

}
