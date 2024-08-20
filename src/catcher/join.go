package catcher

func Join(errs ...error) error {
	return &combinedError{errs}
}

type combinedError struct {
	errs []error
}

func (e *combinedError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}
	s := e.errs[0].Error()
	for _, err := range e.errs[1:] {
		s += ": " + err.Error()
	}
	return s
}
