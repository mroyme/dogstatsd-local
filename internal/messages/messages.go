package messages

import (
	"bytes"
)

type DogStatsDMetricType int

const (
	GaugeMetricType DogStatsDMetricType = iota
	CountMetricType
	SetMetricType
	TimerMetricType
	HistogramMetricType
	DistributionMetricType
)

func (d DogStatsDMetricType) String() string {
	switch d {
	case GaugeMetricType:
		return "gauge"
	case CountMetricType:
		return "count"
	case SetMetricType:
		return "set"
	case TimerMetricType:
		return "timer"
	case HistogramMetricType:
		return "histogram"
	case DistributionMetricType:
		return "distribution"
	}
	return "unknown"
}

type DogStatsDMessageType int

const (
	MetricMessageType DogStatsDMessageType = iota
	ServiceCheckMessageType
	EventMessageType
)

func (d DogStatsDMessageType) String() string {
	switch d {
	case MetricMessageType:
		return "metric"
	case EventMessageType:
		return "event"
	case ServiceCheckMessageType:
		return "service_check"
	}
	return "unknown"
}

type DogStatsDMessage interface {
	Type() DogStatsDMessageType
	Data() []byte
}

// ParseDogStatsDMessage parses a DogStatsDMessage, returning the correct message back
func ParseDogStatsDMessage(buf []byte) (DogStatsDMessage, error) {
	if bytes.HasPrefix(buf, []byte("_e{")) {
		return parseDogStatsDEventMessage(buf)
	}
	if bytes.HasPrefix(buf, []byte("_sc|")) {
		return parseDogStatsDServiceCheckMessage(buf)
	}
	return parseDogStatsDMetricMessage(buf)
}