package ras

type CodecOptions struct {
	Reader  CodecReader
	Writer  CodecWriter
	Version int
}

type Option func(o *CodecOptions)

func WithCodecReader(reader CodecReader) Option {
	return func(o *CodecOptions) {
		o.Reader = reader
	}
}

func WithCodecWriter(writer CodecWriter) Option {
	return func(o *CodecOptions) {
		o.Writer = writer
	}
}

func WithCodecVersion(version int) Option {
	return func(o *CodecOptions) {
		o.Version = version
	}
}
