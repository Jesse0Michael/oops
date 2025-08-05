# Oops, Errors

[![Go Reference](https://pkg.go.dev/badge/github.com/jesse0michael/oops.svg)](https://pkg.go.dev/github.com/jesse0michael/oops)

Error handling for adding attributes and source location to errors when logging them. 

Inspired by [https://github.com/samber/oops](https://github.com/samber/oops) with a minimal API.

Works with [log/slog](https://pkg.go.dev/log/slog) through a [LogValue](https://pkg.go.dev/log/slog#LogValuer) that decorates the logged error attribute with the attributes and source location to any `oops.Error`.

## Example

```go
func (s *Service)process(ctx context.Context, id string) error {
    now := time.Now()
    if id == "" {
        return oops.New("process failed: bad input").With("code", 400)
    }

    results, err := s.Run(ctx, id)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return oops.Errorf("process failed: deadline exceeded %w", err).With("duration", time.Since(now))
        }
        return oops.Wrap(err).With("id", id, "method", "Run")
    }

    return results, nil
}

func main() {
    ctx := context.Background()
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    s := &Service{}

    _, err := s.process(ctx, "oops")
    if err != nil {
        logger.Error("process failed", "error", err)
    }
}
```

```json
{
   "time":"2025-02-01T21:59:54.959869-07:00",
   "level":"ERROR",
   "msg":"process failed",
   "error":{
      "err":"process not found",
      "id":"oops",
      "method":"Run",
      "source":{
         "file":"#####/github.com/jesse0michael/oops/README.md",
         "function":"github.com/jesse0michael/oops.process",
         "line":30
      }
   }
}
```

This example highlights some of the different ways to create and decorate errors by dropping in the oops package. Create a new error with `oops.New` or wrap an existing error with `oops.Wrap` or use `oops.Errorf` to format an error message.

Add attributes to the error with the `With` method. The `With` method can take any number of key-value pairs or slog.Attr's.

Wrapping or formatting an oops.Error with `%w` will append the source location with the new error and preserve attributes of any prior error.

> [!NOTE]  
> runtime source location can be disabled by overwriting the `oops.WithSource` package variable.
> ```go
> oops.WithSource = false 
> ```
>
