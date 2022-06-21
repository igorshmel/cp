package zlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Red = iota + 31
	Green
	Yellow

	Magenta = 35
	Bold    = 1
)

func InitLogger(dev ...bool) zerolog.Logger {

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    true,
		TimeFormat: time.StampMicro,
	}

	if dev != nil && dev[0] {
		logWriter.FormatLevel = func(i interface{}) string {
			if l, ok := i.(string); ok {
				switch l {
				case "trace":
					l = Colorize("TRC", Magenta)
				case "debug":
					l = Colorize("DBG", Green)
				case "info":
					l = Colorize("INF", Green)
				case "warn":
					l = Colorize("WRN", Yellow)
				case "error":
					l = Colorize(Colorize("ERR", Red), Bold)
				case "fatal":
					l = Colorize(Colorize("FTL", Red), Bold)
				case "panic":
					l = Colorize(Colorize("PNC", Red), Bold)
				default:
					l = Colorize("???", Bold)
				}

				return fmt.Sprintf("| %s |", l)
			} else {
				if i == nil {
					return Colorize("???", Bold)
				} else {
					return fmt.Sprintf("| %s |", Colorize(strings.ToUpper(i.(string)), Bold)[0:3])
				}
			}
		}

		logWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}

		logWriter.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}

		return log.Logger.Output(logWriter)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
		return log.Logger
	}
}

func Colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
