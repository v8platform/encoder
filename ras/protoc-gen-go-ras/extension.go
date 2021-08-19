package main

import (
	"github.com/ImSingee/protoc-gen-starlark-go/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MessageExtension struct {
	*options.MessageOption
}

func GetMessageExtensionFor(messageDesc protoreflect.MessageDescriptor) *MessageExtension {
	opts := messageDesc.Options().(*descriptorpb.MessageOptions)
	if opts == nil || !proto.HasExtension(opts, options.E_MessageOption) {
		return nil
	}

	ext := proto.GetExtension(opts, options.E_MessageOption).(*options.MessageOption)

	return &MessageExtension{ext}
}
