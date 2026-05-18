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

		switch dMsg.Type() {
		case messages.MetricMessageType:
			metric, ok := dMsg.(messages.DogStatsDMetric)
			if !ok {
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

		case messages.ServiceCheckMessageType:
			sc, ok := dMsg.(messages.DogStatsDServiceCheck)
			if !ok {
				return nil
			}

			str := fmt.Sprintf("service_check:%s|%s", sc.Name, sc.Status.String())
			if sc.Message != "" {
				str += "|msg:" + sc.Message
			}
			if sc.Hostname != "" {
				str += "|host:" + sc.Hostname
			}

			for _, tag := range append(h.ExtraTags, sc.Tags...) {
				str += " " + tag
			}

			fmt.Println(str)
		}

		return nil
	}
}
