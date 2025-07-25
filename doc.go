/*
Package logit implements a simple yet powerful logging/slog library

# Usage

Init module within an empty folder:

	$ go mod init example.com/demo

	$ cat go.mod
	module example.com/demo

	go 1.24.5

Download/install logit:

	$ go get -u github.com/sfmunoz/logit

Create main.go file:

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

Run it:

	$ go run main.go
	2025-07-14T15:10:26.591Z 0d00h00m00.000s [N] notice-msg app=my-app
	2025-07-14T15:10:26.591Z 0d00h00m00.000s [W] warn-msg app=my-app
	2025-07-14T15:10:26.591Z 0d00h00m00.000s [E] error-msg app=my-app

Detailed configuration:

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

Run it:

	$ go run main.go
	2025-07-14T15:11:42.757Z 0d00h00m00.000s [I] <demo/main.go:19> hello world app=my-app
*/
package logit
