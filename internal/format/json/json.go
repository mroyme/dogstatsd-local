package json

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"os"
)

type Handler struct {
	Logger    *log.Logger
	ExtraTags []string
}

func (h *Handler) New() messages.OutputHandler {
	return func(msg []byte) error {
		dMsg, err := messages.ParseDogStatsDMessage(msg)
		if err != nil {
			h.Logger.Error(err)
		}

		if dMsg.Type() != messages.MetricMessageType {
			h.Logger.Error("unable to serialize non-metric messages to JSON yet")
			return nil
		}

		metric, ok := dMsg.(messages.DogStatsDMetric)
		if !ok {
			h.Logger.Error("could not match message to metric")
			return nil
		}

		jsonMsg := jsonMetric{
			Namespace:  metric.Namespace,
			Name:       metric.Name,
			Path:       fmt.Sprintf("%s.%s", metric.Namespace, metric.Name),
			Value:      metric.FloatValue,
			Extras:     metric.Extras,
			SampleRate: metric.SampleRate,
			Tags:       append(h.ExtraTags, metric.Tags...),
		}

		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(&jsonMsg); err != nil {
			h.Logger.Error("JSON serialize error:", err)
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
