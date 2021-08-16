package ras

import (
	"github.com/k0kubun/pp"
	pb "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		data    []byte
		v       interface{}
		version int
	}

	var ptrInt64 int64 = 222
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"simple",
			args{
				v: Message{
					Type: 111,
					Kind: &ptrInt64,
					Locks: []*Lock{
						{
							1,
							"msg1",
						},
						{
							2,
							"msg2",
						},
					},
					Time: pb.Now(),
				},
				version: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var err error

			tt.args.data, err = Encode(tt.args.v, tt.args.version)
			if err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		pp.Println(tt.args.data)

		var msg Message

		err := Decode(tt.args.data, &msg, 1)
		if err != nil {
			t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
		}

		pp.Println(msg)
		pp.Println(msg.Time.AsTime().String())

	}
}
