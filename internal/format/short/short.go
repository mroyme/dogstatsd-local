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
	return func(msg messages.DogStatsDMessage) error {
		switch msg.Type() {
		case messages.MetricMessageType:
			metric, ok := msg.(messages.DogStatsDMetric)
			if !ok {
				return nil
			}

			str := fmt.Sprintf("metric:%s|%s.%s|%.2f", metric.MetricType.String(), metric.Namespace, metric.Name, metric.FloatValue)

			if metric.MetricType == messages.TimerMetricType {
				str += "ms"
			}

			for _, tag := range append(h.ExtraTags, metric.Tags...) {
				str += " " + tag
			}

			fmt.Println(str)

		case messages.ServiceCheckMessageType:
			sc, ok := msg.(messages.DogStatsDServiceCheck)
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

		case messages.EventMessageType:
			ev, ok := msg.(messages.DogStatsDEvent)
			if !ok {
				return nil
			}

			str := fmt.Sprintf("event:%s", ev.Title)
			if ev.Text != "" {
				str += "|" + ev.Text
			}
			if ev.Priority != "" {
				str += "|priority:" + string(ev.Priority)
			}
			if ev.AlertType != "" {
				str += "|alert:" + string(ev.AlertType)
			}
			if ev.Hostname != "" {
				str += "|host:" + ev.Hostname
			}

			for _, tag := range append(h.ExtraTags, ev.Tags...) {
				str += " " + tag
			}

			fmt.Println(str)
		}

		return nil
	}
}
