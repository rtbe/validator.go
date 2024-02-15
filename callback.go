package validator

import (
	"context"

	"github.com/pkg/errors"
)

type CallbackFunc[T any] func(ctx context.Context, value T) error
type Callback[T any] struct {
	f CallbackFunc[T]
}

func NewCallback[T any](f CallbackFunc[T]) Callback[T] {
	return Callback[T]{
		f: f,
	}
}

func (c Callback[T]) ValidateValue(ctx context.Context, value any) error {
	v, ok := value.(T)
	if !ok {
		var v T
		return errors.Wrapf(CallbackUnexpectedValueTypeError, "got %T want %T", value, v)
	}

	if err := c.f(ctx, v); err != nil {
		var vErr *ValidationError
		if errors.As(err, &vErr) {
			return NewResult().WithError(vErr)
		}

		var result *Result
		if errors.As(err, &result) {
			return result
		}

		return err
	}

	return nil
}
