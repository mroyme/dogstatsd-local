package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"log"
	"os"
)

func NewJSONMessageHandler(extraTags []string) MessageHandler {
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

		jsonMsg := jsonMetric{
			Namespace:  metric.Namespace,
			Name:       metric.Name,
			Path:       fmt.Sprintf("%s.%s", metric.Namespace, metric.Name),
			Value:      metric.FloatValue,
			Extras:     metric.Extras,
			SampleRate: metric.SampleRate,
			Tags:       append(extraTags, metric.Tags...),
		}

		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(&jsonMsg); err != nil {
			log.Println("JSON serialize error:", err.Error())
			return nil
		}

		return nil
	}
}

type jsonMetric struct {
	Namespace  string   `json:"namespace"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Value      float64  `json:"value"`
	Extras     []string `json:"extras"`
	SampleRate float64  `json:"sample_rate"`
	Tags       []string `json:"tags"`
}
