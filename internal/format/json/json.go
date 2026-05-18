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

		switch dMsg.Type() {
		case messages.MetricMessageType:
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
			}

		case messages.ServiceCheckMessageType:
			sc, ok := dMsg.(messages.DogStatsDServiceCheck)
			if !ok {
				h.Logger.Error("could not match message to service check")
				return nil
			}

			jsonMsg := jsonServiceCheck{
				Name:      sc.Name,
				Status:    sc.Status.String(),
				Message:   sc.Message,
				Hostname:  sc.Hostname,
				Tags:      append(h.ExtraTags, sc.Tags...),
				Timestamp: sc.Timestamp.Unix(),
			}

			enc := json.NewEncoder(os.Stdout)
			if err := enc.Encode(&jsonMsg); err != nil {
				h.Logger.Error("JSON serialize error:", err)
			}

		default:
			h.Logger.Error("unable to serialize message type to JSON")
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

type jsonServiceCheck struct {
	Name      string   `json:"name"`
	Status    string   `json:"status"`
	Message   string   `json:"message,omitempty"`
	Hostname  string   `json:"hostname,omitempty"`
	Tags      []string `json:"tags"`
	Timestamp int64    `json:"timestamp"`
}
