.PHONY: build run clean

generate:
	@go run github.com/emiliogrv/errors/cmd/errors_generator \
		-input ./tmpl \
		-output generated \
		-formats somelogger,zap
