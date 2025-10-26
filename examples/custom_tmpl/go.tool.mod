module example.com

go 1.24.5

replace github.com/emiliogrv/errors => ../../

tool github.com/emiliogrv/errors/cmd/errors_generator

require github.com/emiliogrv/errors v0.0.1 // indirect
