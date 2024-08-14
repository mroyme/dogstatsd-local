package main

import (
	"flag"
	"fmt"
	"github.com/mroyme/dogstatsd-local/internal/handler"
	"github.com/mroyme/dogstatsd-local/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Bind address")
	port := flag.Int("port", 8125, "Listen port")
	format := flag.String("format", "pretty", "Output format: json|pretty|raw|short")
	rawTags := flag.String("tags", "", "Extra tags, comma delimited")
	maxNameWidth := flag.Int("max-name-width", 50,
		"Maximum length of name. Only used for 'pretty' format, increase if name is truncated")
	maxValueWidth := flag.Int("max-value-width", 15,
		"Maximum length of value. Only used for 'pretty' format, increase if value is truncated")
	flag.Parse()

	extraTags := strings.Split(*rawTags, ",")
	var msgHandler handler.MessageHandler

	switch *format {
	case "json":
		msgHandler = handler.NewJSONMessageHandler(extraTags)
	case "short":
		msgHandler = handler.NewShortMessageHandler(extraTags)
	case "raw":
		msgHandler = handler.NewRawMessageHandler()
	default:
		msgHandler = handler.NewPrettyMessageHandler(extraTags, *maxNameWidth, *maxValueWidth)
	}
	asyncHandler := handler.NewAsyncMessageHandler(msgHandler, 1000, 10000)

	var wg sync.WaitGroup

	// create a new server and listen on a background goroutine
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("listening over UDP at ", addr)
	srv := server.NewServer(addr, asyncHandler.Handler)
	wg.Add(1)
	go func(srv server.Server) {
		defer wg.Done()
		if err := srv.Listen(); err != nil {
			log.Fatal(err.Error())
		}
	}(srv)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	if err := srv.Stop(); err != nil {
		log.Println(err.Error())
	}
	wg.Wait()
	asyncHandler.Stop()
}
