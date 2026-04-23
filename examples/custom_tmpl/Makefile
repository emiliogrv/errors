.PHONY: build run clean

generate:
	@rm -rf ./generated/*
	@cd tools && \
		go mod tidy && \
		go tool errors_generator \
		-input-dir ../tmpl \
		-output-dir ../generated \
		-formats somelogger,zap
