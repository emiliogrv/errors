# Full Import Example

This example demonstrates importing the full package which includes all logging framework integrations. By importing the
full package, the following logging frameworks and their dependencies will be added to your `go.mod`:

- [Zap](https://github.com/uber-go/zap)
- [Logrus](https://github.com/sirupsen/logrus)
- [Zerolog](https://github.com/rs/zerolog)

Note that this will increase your project's dependency footprint. If you only need specific logging framework
integration, consider using the partial import approach instead (see the `partial_import` or `custom_tpml` example).

