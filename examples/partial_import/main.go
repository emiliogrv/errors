package main

import (
	stderrors "errors"
	"fmt"
	"runtime/debug"
	"time"

	"go.uber.org/zap"

	errors "github.com/emiliogrv/errors/pkg/zap"
)

func main() {
	err := anError()

	logZap(err)
}

func anError() error {
	var nilError *errors.StructuredError

	attrs := []errors.Attr{
		errors.Int("code", 100),
		errors.Ints("numbers", 10, 20, 30, 40),
		errors.String("text", "example text"),
		errors.Bool("bool", true),
		errors.Bools("bools", true, false, true),
		errors.Duration("duration", time.Second), // be aware that each marshaler has its own way of handling this
		errors.Time("time", time.Now()),          // be aware that each marshaler has its own way of handling this
	}

	// err := errors.New("this is an example error message with all marshalers") // this is the same as the next line
	var err error = errors.New("this is an example error message with all marshalers").
		WithAttrs(attrs...).
		WithErrors(
			errors.New(""). // empty error message => !NILVALUE
				WithErrors(
					nilError, // nil error => !NILVALUE
					fmt.Errorf("fmt error"),
					errors.New("errors new nested 1"),
					errors.New("error with more nested errors").
						WithErrors(
							nilError, // nil error => !NILVALUE
							fmt.Errorf("fmt error"),
							errors.New("errors new nested 2"),
						),
				),
			fmt.Errorf(
				"text example %w", // "text example" will be lost
				errors.New("the 'text example' part will be lost"),
			),
			fmt.Errorf(
				"%w, %w",
				errors.New("wrapped 1"),
				errors.New("wrapped 2"),
			),
			stderrors.New("std error"),
			fmt.Errorf("fmt error"),
			nil,               // nil  => !NILVALUE
			nilError,          // nil error => !NILVALUE
			stderrors.New(""), // empty error message => !NILVALUE
			stderrors.Join(
				stderrors.New("std joined error"),
				stderrors.Join(
					stderrors.New("std joined nested error"),
				),
			),
			errors.Join(
				stderrors.New("std joined error"),
				stderrors.Join(
					stderrors.New("std joined nested error"),
					errors.New("errors joined nested error").
						WithAttrs(attrs...).
						WithErrors(
							stderrors.New("std nested error 1"),
							errors.New(""),
						),
				),
			),
		).
		WithTags("tag1", "tag2", "tag3").
		WithStack(debug.Stack())

	return err
}

func logZap(err error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	logger.Error("zap json", zap.Any("err", err))
}
