package channels

type Result[T any] struct {
	value T
	err   error
}

func NewResult[T any](value T, err error) *Result[T] {
	return &Result[T]{value: value, err: err}
}

func NewErrResult[T any](err error) *Result[T] {
	return &Result[T]{err: err}
}

func (r *Result[T]) Unpack() (T, error) {
	return r.value, r.err
}

type PairResult[T any, U any] struct {
	first  T
	second U
	err    error
}

func NewPairResult[T any, U any](first T, second U, err error) *PairResult[T, U] {
	return &PairResult[T, U]{first: first, second: second, err: err}
}

func NewErrPairResult[T any, U any](err error) *PairResult[T, U] {
	return &PairResult[T, U]{err: err}
}

func (r *PairResult[T, U]) Unpack() (T, U, error) {
	return r.first, r.second, r.err
}
