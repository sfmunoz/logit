# logit

logit: Simple Golang logging/slog library

- [Features](#features)
- [TL;DR](#tldr)
- [References](#references)
- [Others](#others)

> [!CAUTION]
> Even though it's ready to be used the repository is evolving at a fast pace these days so it's advisable to stick to a fixed version to prevent compilation from breaking because of API changes.\
> E.g.:
> - OK → `go get -u github.com/sfmunoz/logit@v0.2.0`
> - KO → `go get -u github.com/sfmunoz/logit`

## Features

- No dependencies
- Fluent configuration (one line)
  - Simple and easy
  - Compact
- Extra log levels
  - Usual: debug, info, warn, error
  - Added: trace, notice, fatal, panic
- Color support
- slog based so it requires Golang 1.21+
- slog.Logger + slog.Handler compatible
- Opinionated

## TL;DR

Init module within an empty folder:
```
$ go mod init example.com/demo

$ cat go.mod 
module example.com/demo

go 1.24.5
```
Download/install **logit**:
```
$ go get -u github.com/sfmunoz/logit
```
Create **main.go** file:
```go
package main

import "github.com/sfmunoz/logit"

var log = logit.Logit().
	WithLevel(logit.LevelNotice).
	With("app", "my-app")

func main() {
	log.Trace("trace-msg")
	log.Debug("debug-msg")
	log.Info("info-msg")
	log.Notice("notice-msg")
	log.Warn("warn-msg")
	log.Error("error-msg")
}
```
Run it using `go run main.go`:

![20250714_151107.png](https://github.com/sfmunoz/logit/blob/assets/20250714_151107.png)

Detailed configuration:
```go
package main

import (
	"log/slog"
	"os"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	With("app", "my-app").
	WithWriter(os.Stderr).
	WithTpl(logit.TplTime, logit.TplUptime, logit.TplLevel, logit.TplSource, logit.TplMessage, logit.TplAttrs).
	WithLevel(slog.LevelDebug).
	WithTimeFormat("2006-01-02T15:04:05.000Z07:00").
	WithColor(true)

func main() {
	log.Info("hello world")
}
```
Run it too with `go run main.go`:

![20250714_151226.png](https://github.com/sfmunoz/logit/blob/assets/20250714_151226.png)

## References

```
$ go doc slog
(...)
For a guide to writing a custom handler, see https://golang.org/s/slog-handler-guide
(...)
```

> [!TIP]
> So the must reading guide is the following one:\
> [https://golang.org/s/slog-handler-guide](https://golang.org/s/slog-handler-guide) → [https://github.com/golang/example/blob/master/slog-handler-guide/README.md](https://github.com/golang/example/blob/master/slog-handler-guide/README.md)

There are 4 versions of the indenthandler example in the README.md:

> [https://github.com/golang/example/tree/master/slog-handler-guide](https://github.com/golang/example/tree/master/slog-handler-guide)

Google: 'golang slog'

- https://go.dev/blog/slog
- https://pkg.go.dev/golang.org/x/exp/slog

slog related videos:

- [A consistent logging format for Go](https://www.youtube.com/watch?v=gd_Vyb5vEw0)
- [Golang: Structured Logging using slog](https://www.youtube.com/watch?v=gVL-Ilbj168)

## Others

Google: 'golang slog color'

- [https://github.com/lmittmann/tint](https://github.com/lmittmann/tint)
- [https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package](https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package)

Google: 'golang log packages'

- [https://betterstack.com/community/guides/logging/best-golang-logging-libraries/](https://betterstack.com/community/guides/logging/best-golang-logging-libraries/)
- [https://blog.logrocket.com/5-structured-logging-packages-for-go/](https://blog.logrocket.com/5-structured-logging-packages-for-go/)
- [https://www.highlight.io/blog/5-best-logging-libraries-for-go](https://www.highlight.io/blog/5-best-logging-libraries-for-go)

zap & zerolog:

- [https://github.com/uber-go/zap](https://github.com/uber-go/zap)
- [https://github.com/rs/zerolog](https://github.com/rs/zerolog)
