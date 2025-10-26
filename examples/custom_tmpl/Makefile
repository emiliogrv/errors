.PHONY: build run clean

generate:
	@go tool -modfile=go.tool.mod errors_generator \
		-input-dir ./tmpl \
		-output-dir generated \
		-formats somelogger,zap
