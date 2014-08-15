glogger
=======

[![GoDoc](https://godoc.org/github.com/Xuyuanp/glogger?status.svg)](https://godoc.org/github.com/Xuyuanp/glogger)    
[![](https://travis-ci.org/Xuyuanp/glogger.svg?branch=master)](https://travis-ci.org/Xuyuanp/glogger)



A python-like logging library for go

## Be Carefull!!!

This is still a very early version under development, and far away from stable.

## Sample

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

### Config

```json
{
    "formatters": {
        "default": {
            "builder": "github.com/Xuyuanp/glogger/formatters.DefaultFormatter",
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
            "builder": "github.com/Xuyuanp/glogger/handlers.StreamHandler",
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
