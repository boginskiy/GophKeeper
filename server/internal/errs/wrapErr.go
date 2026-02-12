package errs

type ErrWrapper struct {
	msg string
	err error
}

func NewErrWrapper(mess string) *ErrWrapper {
	return &ErrWrapper{msg: mess}
}

func (e *ErrWrapper) Error() string {
	return e.msg
}

func (e *ErrWrapper) Wrap(err error) error {
	e.err = err
	return e
}
