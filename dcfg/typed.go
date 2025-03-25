package dcfg

import "context"

type TypedValue[T any] struct {
	backend Backend
	type_   Type
	key     Key
}

func Typed[T any](b Backend, key Key) TypedValue[T] {
	var t T
	return TypedValue[T]{
		backend: b,
		type_:   TypeOf(t),
		key:     key,
	}
}

func (tv *TypedValue[T]) Load(ctx context.Context) (out T, err error) {
	err = tv.backend.Load(ctx, tv.type_, tv.key, &out)
	return
}

func (tv *TypedValue[T]) Store(ctx context.Context, value T) error {
	return tv.backend.Store(ctx, tv.type_, tv.key, value)
}
