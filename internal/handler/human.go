package handler

import (
	"fmt"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"log"
)

func NewShortMessageHandler(extraTags []string) MessageHandler {
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
