package errors

// Join returns an error that wraps the given errors, any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats depending on logging format otherwise as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	count := zero

	for _, err := range errs {
		if err != nil {
			count++
		}
	}

	if count == zero {
		return nil
	}

	_err := &StructuredError{
		joined: true,
	}

	for _, err := range errs {
		if err != nil {
			_err.Errors = append(_err.Errors, err)
		}
	}

	return _err
}

// JoinIf is similar to Join, but it will only join the errors if the first error is not nil.
// If the first error is nil, it will return nil, otherwise it will join all the errors.
// The error formats depending on logging format otherwise as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by JoinIf implements the Unwrap() []error method.
func JoinIf(errs ...error) error {
	if len(errs) == zero {
		return nil
	}

	if errs[zero] != nil {
		return Join(errs...)
	}

	return nil
}
