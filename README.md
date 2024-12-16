# TgLoggerApi

## Using TgLoggerApi

TgLoggerApi is a tool, that can connect to [TgLogger](https://github.com/vandi37/TgLogger). 

The realization is an io.Writer, so it could be used in any loggers.

## Example

```go
w := New("my token", 0) // NewWithUrl("my token", 0, "url")

log.SetOutput(w) // Example 

w.Write([]byte("my log")) // Writing
```

## License 

[MIT](LISENSE)