package main

import (
	"flag"
	"fmt"
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/format"
	"github.com/mroyme/dogstatsd-local/internal/format/json"
	"github.com/mroyme/dogstatsd-local/internal/format/pretty"
	"github.com/mroyme/dogstatsd-local/internal/format/raw"
	"github.com/mroyme/dogstatsd-local/internal/format/short"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"github.com/mroyme/dogstatsd-local/internal/server"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Bind address")
	port := flag.Int("port", 8125, "Listen port")
	out := flag.String("out", "pretty", "Output format: json|pretty|raw|short")
	rawTags := flag.String("tags", "", "Extra tags, comma delimited")
	maxNameWidth := flag.Int("max-name-width", 50,
		"Maximum length of name. Only used for 'pretty' format, increase if name is truncated."+
			" Values below 50 have no effect.")
	maxValueWidth := flag.Int("max-value-width", 15,
		"Maximum length of value. Only used for 'pretty' format, increase if value is truncated")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	extraTags := strings.Split(*rawTags, ",")
	logger := getLogger(*debug)

	var formatHandler format.Handler
	switch *out {
	case "json":
		logger.SetFormatter(log.JSONFormatter)
		formatHandler = &json.Handler{Logger: logger, ExtraTags: extraTags}
	case "short":
		formatHandler = &short.Handler{Logger: logger, ExtraTags: extraTags}
	case "raw":
		formatHandler = &raw.Handler{Logger: logger}
	default:
		theme := &pretty.CatppuccinAdaptiveTheme{
			Light: catppuccin.Latte,
			Dark:  catppuccin.Mocha,
		}
		formatHandler = &pretty.Handler{
			Logger:     logger,
			Theme:      theme,
			ExtraTags:  extraTags,
			NameWidth:  *maxNameWidth,
			ValueWidth: *maxValueWidth,
		}
	}

	asyncMessageHandler := messages.Handler{
		Logger:     logger,
		Out:        formatHandler.New(),
		PoolSize:   1000,
		BufferSize: 10000,
	}
	messageHandler := asyncMessageHandler.New()

	var wg sync.WaitGroup

	// create a new server and listen on a background goroutine
	addr := fmt.Sprintf("%s:%d", *host, *port)
	logger.Infof("listening over UDP at %s", addr)
	srv := server.NewServer(addr, messageHandler.Handle, logger)
	wg.Add(1)
	go func(srv server.Server) {
		defer wg.Done()
		if err := srv.Listen(); err != nil {
			logger.Fatal(err)
		}
	}(srv)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	if err := srv.Stop(); err != nil {
		logger.Error(err)
	}
	wg.Wait()
	messageHandler.Stop()
}

func getLogger(debug bool) *log.Logger {
	var logLevel log.Level
	var logReportCaller bool
	logTimeFormat := time.TimeOnly
	if debug {
		logLevel = log.DebugLevel
		logReportCaller = true
		logTimeFormat = time.RFC3339
	}

	return log.NewWithOptions(os.Stderr, log.Options{
		Level:           logLevel,
		ReportCaller:    logReportCaller,
		ReportTimestamp: true,
		TimeFormat:      logTimeFormat,
	})
}
