package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"log"
	"os"
	"sync"
)

type MsgHandler func([]byte) error

type AsyncMsgHandler interface {
	Handler([]byte) error
	Stop()
}

type handler struct {
	fn    MsgHandler
	msgCh chan []byte
	wg    sync.WaitGroup
}

func NewAsyncMsgHandler(fn MsgHandler, poolSize int, bufferSize int) AsyncMsgHandler {
	h := &handler{
		fn:    fn,
		msgCh: make(chan []byte, bufferSize),
	}

	h.wg.Add(poolSize)

	// build a pool of goroutines to listen for messages to process
	for i := 0; i < poolSize; i++ {
		go func() {
			defer h.wg.Done()
			for msg := range h.msgCh {
				err := h.fn(msg)
				if err != nil {
					return
				}
			}
		}()
	}

	return h
}

// Handler submits to the pool
func (a *handler) Handler(msg []byte) error {
	select {
	case a.msgCh <- msg:
		return nil
	default:
	}

	return errors.New("POOL_CAPACITY_EXCEEDED")
}

func (a *handler) Stop() {
	close(a.msgCh)
	a.wg.Wait()
}

type dogstatsdJsonMetric struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Path      string `json:"path"`

	Value      float64  `json:"value"`
	Extras     []string `json:"extras"`
	SampleRate float64  `json:"sample_rate"`
	Tags       []string `json:"tags"`
}

func NewJSONDogStatsDMsgHandler(extraTags []string) MsgHandler {
	return func(msg []byte) error {
		dMsg, err := messages.ParseDogStatsDMessage(msg)
		if err != nil {
			log.Println(err.Error())
		}

		if dMsg.Type() != messages.MetricMessageType {
			log.Println("Unable to serialize non metric messages to JSON yet")
			return nil
		}

		metric, ok := dMsg.(messages.DogStatsDMetric)
		if !ok {
			log.Fatalf("Programming error: invalid Type() = type matching")
		}

		jsonMsg := dogstatsdJsonMetric{
			Namespace:  metric.Namespace,
			Name:       metric.Name,
			Path:       fmt.Sprintf("%s.%s", metric.Namespace, metric.Name),
			Value:      metric.FloatValue,
			Extras:     metric.Extras,
			SampleRate: metric.SampleRate,
			Tags:       metric.Tags,
		}

		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(&jsonMsg); err != nil {
			log.Println("JSON serialize error:", err.Error())
			return nil
		}

		return nil
	}
}

func NewHumanDogStatsDMsgHandler(extraTags []string) MsgHandler {
	return func(msg []byte) error {
		dMsg, err := messages.ParseDogStatsDMessage(msg)
		if err != nil {
			log.Println(err.Error())
			return nil
		}

		metric, ok := dMsg.(messages.DogStatsDMetric)
		if dMsg.Type() != messages.MetricMessageType || !ok {
			return nil
		}

		tmpl := "metric:%s|%s.%s|%.2f"
		str := fmt.Sprintf(tmpl, metric.MetricType.String(), metric.Namespace, metric.Name, metric.FloatValue)

		if metric.MetricType == messages.TimerMetricType {
			str += "ms"
		}

		// iterate through tags
		for _, tag := range append(extraTags, metric.Tags...) {
			str += " " + tag
		}

		fmt.Println(str)
		return nil
	}
}

func NewRawDogStatsDMsgHandler() MsgHandler {
	return func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	}
}
