package errors

type pubErr struct {
	err error
	msg string
}

func Public(err error, msg string) error {
	return pubErr{err, msg}
}

func (pe pubErr) Error() string {
	return pe.err.Error()
}

func (pe pubErr) Public() string {
	return pe.msg
}

func (pe pubErr) Unwrap() error {
	return pe.err
}
