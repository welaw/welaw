package test

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
)

type testEnv struct {
	repoPath string
}

func NewLogger(out io.Writer) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}
