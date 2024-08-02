package flags

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/samber/oops"
	slogmulti "github.com/samber/slog-multi"
	"github.com/urfave/cli/v3"
	"github.com/zackarysantana/velocity/internal/stats"
	"github.com/zackarysantana/velocity/internal/vlog"
)

type LoggerMode int

const (
	Debug LoggerMode = iota
	Quiet
	maxLoggerMode
)

func (m LoggerMode) String() string {
	return loggerModes[m]
}

func (m LoggerMode) Valid() bool {
	return 0 <= m && m < maxLoggerMode
}

func getLoggerMode(mode string) *LoggerMode {
	for k, v := range loggerModes {
		if v == mode {
			return &k
		}
	}
	return nil
}

func getAllLoggerModes() []string {
	modes := make([]string, 0, len(loggerModes))
	for _, v := range loggerModes {
		modes = append(modes, v)
	}
	return modes
}

var (
	loggerModes = map[LoggerMode]string{
		Debug: "debug",
		Quiet: "quiet",
	}
	LoggerModeFlag = &cli.StringFlag{
		Name:  "mode",
		Usage: fmt.Sprintf("set the mode to `%s`", strings.Join(getAllLoggerModes(), "|")),
		Validator: func(s string) error {
			if s == "" {
				return nil
			}
			if getLoggerMode(s) == nil {
				return oops.In("flags").Errorf("mode must be one of: %s", strings.Join(getAllLoggerModes(), ", "))
			}
			return nil
		},
	}
	VerboseFlag = &cli.BoolFlag{
		Name:  "verbose",
		Usage: "enable verbose output",
	}
)

func SetLogger(_ context.Context, cmd *cli.Command) error {
	level := slog.LevelInfo
	switch cmd.String("mode") {
	case "":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "error":
		level = slog.LevelError
	case "quiet":
		level = 12
	default:
		return oops.In("flags").Errorf("mode must be one of: %s", strings.Join(getAllLoggerModes(), ", "))
	}
	var stdLogger slog.Handler
	if cmd.Bool("verbose") {
		stdLogger = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		stdLogger = vlog.NewPlainHandler(os.Stdout, &vlog.Options{Level: level})
	}

	logger := slog.New(slogmulti.Fanout(
		stdLogger,
	))
	cmd.Metadata["logger"] = logger

	logger.Debug("Starting logger", "cmd", os.Args, "version", cmd.Version, "ip", stats.GetIP())
	return nil
}

func Logger(cmd *cli.Command) *slog.Logger {
	return cmd.Root().Metadata["logger"].(*slog.Logger)
}
