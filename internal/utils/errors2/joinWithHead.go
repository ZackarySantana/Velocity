package errors2

import "strings"

type joinWithHeadError struct {
	errs []error
}

func JoinWithHead(head error, errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinWithHeadError{
		errs: make([]error, 0, n+1),
	}
	e.errs = append(e.errs, head)
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

func (e *joinWithHeadError) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i == 0 {
			b = append(b, err.Error()...)
			continue
		}
		e := err.Error()
		p := strings.Split(e, "\n")
		for i, l := range p {
			b = append(b, '\n')
			b = append(b, '\t')
			if i == 0 {
				b = append(b, "- "...)
			} else {
				b = append(b, "  "...)
			}
			b = append(b, l...)
		}
	}
	return string(b)
}

func (e *joinWithHeadError) Unwrap() []error {
	return e.errs
}
