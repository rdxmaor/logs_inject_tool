package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
	timeFormat    = "2006-01-02T15:04:05.999Z07:00"
)

var (
	pid int
)

func main() {

	zerolog.TimeFieldFormat = timeFormat
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat,
			FormatLevel:  LoganFormatLevel(false),
			FormatCaller: LoganFormatCaller(false),
		})

	pid = os.Getpid()

	log.Logger = log.With().
		Caller().
		Logger()

	signalChan := make(chan os.Signal, 1)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func(quitChannel chan os.Signal, cancelFunc context.CancelFunc) {

		<-quitChannel
		cancelFunc()

	}(signalChan, cancel)

	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for ctx.Err() != context.Canceled {

		for i := 0; i < 10000; i++ {

			log.Info().Msg("This is a test log!")
		}

		time.Sleep(1 * time.Second)
	}
}

func LoganFormatLevel(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = colorize("DEBUG", colorYellow, noColor)
			case "info":
				l = colorize("INFO", colorGreen, noColor)
			case "warn":
				l = colorize("WARN", colorRed, noColor)
			case "error":
				l = colorize(colorize("ERROR", colorRed, noColor), colorBold, noColor)
			case "fatal":
				l = colorize(colorize("FATAL", colorRed, noColor), colorBold, noColor)
			case "panic":
				l = colorize(colorize("PANIC", colorRed, noColor), colorBold, noColor)
			default:
				l = colorize(" ??? ", colorBold, noColor)
			}
		} else {
			l = strings.ToUpper(fmt.Sprintf("%s", i))[0:5]
		}
		return l
	}
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func LoganFormatCaller(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = filepath.Base(c)
				c = fmt.Sprintf("%s[%d]:[gid:%d] ", "logs_inject_tool", pid, getGID()) + colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
			}
		}
		return c
	}
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := ConvertStringToUint64(string(b))

	return n
}

func ConvertStringToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}
