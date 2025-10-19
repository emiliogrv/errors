package main

import (
	"go.uber.org/zap"

	errors "github.com/emiliogrv/errors/examples/customtmpl/generated"
)

func main() {
	err := errors.New("this is an example error for custom marshalers and user templates").
		WithAttrs(
			errors.Bool("overrode", errors.MaxDepthMarshal() == 10), // true
			errors.Int("max_depth", errors.MaxDepthMarshal()),
		)

	err.MarshalSomeLogger() // here just to show that it was generated.

	logZap(err)
}

func logZap(err error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	logger.Error("zap json", zap.Any("err", err))
}
