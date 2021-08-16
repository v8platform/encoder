package ras

import (
	"fmt"
)

type DecoderError struct {
	fn        string
	err       error
	needBytes int
	readBytes int
}

func (e *DecoderError) Error() string {

	return fmt.Sprintf("decoderFunc: fn<%s> err<%s>", e.fn, e.err.Error())

}
