package storage

import "context"

type RContext interface {
	context.Context
	Reader
}

func NewRContext(ctx context.Context, r Reader) RContext{
	return struct{
		context.Context
		Reader
	}{
		Context: ctx,
		Reader:  r,
	}
}

type RWContext interface {
	context.Context
	ReaderWriter
}

func NewRWContext(ctx context.Context, rw ReaderWriter) RWContext{
	return struct{
		context.Context
		ReaderWriter
	}{
		Context: ctx,
		ReaderWriter:  rw,
	}
}