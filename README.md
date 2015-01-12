glogger
=======

[![GoDoc](https://godoc.org/github.com/Xuyuanp/glogger?status.svg)](https://godoc.org/github.com/Xuyuanp/glogger)    

[![Travis CI](https://travis-ci.org/Xuyuanp/glogger.svg?branch=master)](https://travis-ci.org/Xuyuanp/glogger)

A python-like logging library for go

## Getting Started

Install glogger package:

`go get github.com/Xuyuanp/glogger`

Create your `go` file. We'll call it `hello.go`.

```go
package main

import (
    "github.com/Xuyuanp/glogger"
    )

func main() {
    glogger.Default().Debug("Hello world!")
    glogger.Default().Info("Hello world again!")
)
```

Then run your code:

`go run hello.go`

You will now see the logging output in the terminal.

## Configuration Instruction

Config file is written in json format.

* `filters`: filter list. (nothing currently)
* `formatters`: formatter list.
    1. `builder`: formatter builder name.
        * `github.com/Xuyuanp/glogger.DefaultFormatter`: default formatter builder.
        * `github.com/Xuyuanp/glogger/formatters.RainbowFormatter`: format record colorized.
    2. `fmt`: the log message format. The following macros can be used with ${} (optional):
        * `name`: logger name
        * `levelno`: log level number
        * `levelname`: log level name
        * `time`: current time
        * `lfile`: file name with absolute path
        * `sfile`: file name
        * `line`: current code line
        * `func`: function name
        * `msg`: log message
        * other color macro for RainbowFormatter
    3. `timefmt`: format of time. (optional)
    4. `colors`: color map for RainbowFormatter. See the config sample above. (optional)
* `handlers`: handler list.
    1. `builder`: handler builder name, values:
        * `github.com/Xuyuanp/glogger.StreamHandler`: Output log message into stream. (default)
        * `github.com/Xuyuanp/glogger/handlers.FileHandler`: Output log message into file.
        * `github.com/Xuyuanp/glogger/handlers.RotatedFileHandler`: Output log message into file and auto-rotated.
        * `github.com/Xuyuanp/glogger/handlers.SMTPHandler`: Output log message via SMTP.
    1. `level`: log level, values: (optional)
        * `DEBUG` (default)
        * `INFO`
        * `WARNING`
        * `ERROR`
        * `CRITICAL`
    2. `filters`: filter name list. (optional)
    3. `formatter`: formatter name, DefaultFormatter if not supplied. (optional)
    4. `writer`: the output stream, for StreamHandler. (optional)
        * `stdout`: standard output (default)
        * `stderr`: standard error
    5. `filename`: file name, for FileHandler and RotatedFileHandler. (required)
    6. `autoRotated`: if enable log file auto-rotated. Boolean value, `true` or `false`. For RotatedFileHandler. (optional, `true` as default)
    7. `maxSize`: max size (byte) of log file. Integer value, 0 means unlimited. For RotatedFileHandler. (optional, `0` as default)
    8. `maxLine`: max line of log file. Integer value, 0 means unlimited. For RotatedFileHandler. (optional, `0` as default)
    9. `daily`: if auto-rotate log file daily. Boolean value, `true` or `false`. For RotatedFileHandler. (optional, `true` as default)
    10. `address`: email address to send log message from, for SMTPHandler. (required)
    11. `username`: SMTP server username, for SMTPHandler. (required)
    12. `password`: SMTP server password, for SMTPHandler. (required)
    13. `to`: target email address list, for SMTPHandler. (required)
    14. `subject`: email subject, for SMTPHandler. (required)
* `loggers`: logger list.
    1. `level`: log level. default is `DEBUG`. (optional)
    2. `filters`: filter name list. (optional)
    3. `handlers`: handler name list, default is StreamHandler. (optional)

## Further Sample

### Code

```go
package main

import (
    "github.com/Xuyuanp/glogger"
    _ "github.com/Xuyuanp/glogger/formatters"
    _ "github.com/Xuyuanp/glogger/handlers"
)

func init() {
    glogger.LoadConfigFromFile("log.conf")
}

func main() {
    logger := glogger.GetLogger("main")

    logger.Debug("This DEBUG message")
    logger.Info("This is INFO message")
    logger.Warning("This is WARNING message")
    logger.Error("This is ERROR message")
    logger.Critical("This is CRITICAL message")
}
```

### Config file: `log.conf`

```json
{
    "formatters": {
        "default": {
            "builder": "github.com/Xuyuanp/glogger.DefaultFormatter",
            "fmt": "${time} ${levelname} ${sfile}:${line} ${msg}",
            "timefmt": "2006-01-02 15:04:05"
        },
        "rainbow": {
            "builder": "github.com/Xuyuanp/glogger/formatters.RainbowFormatter",
            "fmt": "${log_color}[${time} ${levelname} ${sfile}:${line} ${func}] ${msg}",
            "timefmt": "01-02-2006 15:04:05",
            "colors": {
                "DEBUG":    "blue",
                "INFO":     "green",
                "WARNING":  "yellow",
                "ERROR":    "cyan",
                "CRITICAL": "red"
            }
        }
    },
    "handlers": {
        "file": {
            "builder": "github.com/Xuyuanp/glogger/handlers.FileHandler",
            "level": "INFO",
            "formatter": "default",
            "filename": "log/record.log"
        },
        "console": {
            "builder": "github.com/Xuyuanp/glogger.StreamHandler",
            "level": "DEBUG",
            "writer": "stdout",
            "formatter": "rainbow"
        }
    },
    "loggers": {
        "main": {
            "level": "DEBUG",
            "handlers": ["console", "file"]
        }
    }
}

```

### Result

![Imgur](http://i.imgur.com/xjWGUyC.png)

