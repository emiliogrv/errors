module github.com/emiliogrv/errors/examples/partialimport

go 1.24.5

replace github.com/emiliogrv/errors => ../../

require (
	github.com/emiliogrv/errors v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect
