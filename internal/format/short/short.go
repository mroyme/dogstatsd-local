package short

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
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
		for _, tag := range append(h.ExtraTags, metric.Tags...) {
			str += " " + tag
		}

		fmt.Println(str)
		return nil
	}
}
