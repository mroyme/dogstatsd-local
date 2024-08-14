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
	host := flag.String("host", "0.0.0.0", "bind address")
	port := flag.Int("port", 8125, "listen port")
	format := flag.String("format", "stdout", "output format: json|std|raw")
	rawTags := flag.String("tags", "", "extra tags: comma delimited")
	flag.Parse()

	extraTags := strings.Split(*rawTags, ",")
	var msgHandler handler.MsgHandler

	switch *format {
	case "json":
		msgHandler = handler.NewJSONDogStatsDMsgHandler(extraTags)
	case "human":
		msgHandler = handler.NewHumanDogStatsDMsgHandler(extraTags)
	default:
		msgHandler = handler.NewRawDogStatsDMsgHandler()
	}
	asyncHandler := handler.NewAsyncMsgHandler(msgHandler, 1000, 10000)

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
