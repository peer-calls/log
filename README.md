# github.com/peer-calls/log

A versatile logger configurable via environment variables. It was first built
for [peer-calls](https://github.com/peer-calls/peer-calls), an open-source,
distributed peer to peer conferencing software, but I found it useful in a lot
of other projects so I'm importing it here.

```go
package main

import (
	"errors"

	"github.com/peer-calls/log"
)

func main() {
	cfg := log.NewConfig(log.ConfigMap{
		"test1": log.LevelWarn,
		"test2": log.LevelInfo,
	})

	factory := log.New().WithConfig(cfg)

	logCtx := log.Ctx{
		"a": "1000",
		"b": 20,
	}

	l1 := factory.WithNamespace("test1")
	l2 := factory.WithNamespace("test2")

	l1.Info("test", logCtx)
	l2.Warn("test", logCtx)
	l2.Error("test", errors.New("err1"), logCtx)

	l2.Info("test", logCtx)
	l2.Warn("test", logCtx)
	l2.Error("test", errors.New("err2"), logCtx)
}
```

After running the script above you should see the output below:

```go
$ go run ./main.go
2022-01-27T09:38:46.952896+01:00  warn [               test2] test a=1000 b=20
2022-01-27T09:38:46.953052+01:00 error [               test2] test: err1 a=1000 b=20
2022-01-27T09:38:46.953068+01:00  info [               test2] test a=1000 b=20
2022-01-27T09:38:46.953080+01:00  warn [               test2] test a=1000 b=20
2022-01-27T09:38:46.953093+01:00 error [               test2] test: err2 a=1000 b=20
```

This is only the default logging format. It is configurable. For more details
see the docs: https://pkg.go.dev/github.com/peer-calls.

## Golang CI Lint

If you're using `.golangci-lint.yml` and don't want to handle write errors,
you should add this to your linter settings:

```yaml
linters-settings:
  errcheck:
    exclude: .golangci-errcheck-exclude.txt
```

And then to `.golangci-errcheck-exclude.txt`:

```
(github.com/peer-calls/peer-calls/log.Logger).Error
(github.com/peer-calls/peer-calls/log.Logger).Warn
(github.com/peer-calls/peer-calls/log.Logger).Info
(github.com/peer-calls/peer-calls/log.Logger).Debug
(github.com/peer-calls/peer-calls/log.Logger).Trace
```

## License

Apache v2
