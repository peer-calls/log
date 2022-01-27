# logger

A versatile logger configurable via environment variables.

```go
package main

import "github.com/peer-calls/log"

func main() {
  cfg := NewConfig(log.ConfigMap{
    "test1": log.LevelWarn,
    "test2": log.LevelInfo,
  })

  factory := log.New().WithConfig(cfg)

  logCtx := log.Ctx{
    "a": "1000",
    "b": 20,
  }

  l1 := l.WithNamespace("test1")
  l2 := log.New()

  l1.Info("test", logCtx)
  l2.Warn("test", logCtx)
  l3.Error("test", logCtx)

  l1.Info("test", logCtx)
  l2.Warn("test", logCtx)
  l3.Error("test", logCtx)
}
```

If you're using `.golangci-lint.yml` and don't want to handle write errors,
you should add this to your linter settings:

```yaml
linters-settings:
  errcheck:
    exclude: .golangci-errcheck-exclude.txt
```

And then to `.golangci-errcheck-exclude.txt`:

```
(github.com/peer-calls/peer-calls/server/logger.Logger).Error
(github.com/peer-calls/peer-calls/server/logger.Logger).Warn
(github.com/peer-calls/peer-calls/server/logger.Logger).Info
(github.com/peer-calls/peer-calls/server/logger.Logger).Debug
(github.com/peer-calls/peer-calls/server/logger.Logger).Trace
```

# License

Apache v2
