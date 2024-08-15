package messages

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"
)

type DogStatsDMetricType int

const (
	GaugeMetricType DogStatsDMetricType = iota
	CounterMetricType
	SetMetricType
	TimerMetricType
	HistogramMetricType
)

func (d DogStatsDMetricType) String() string {
	switch d {
	case GaugeMetricType:
		return "gauge"
	case CounterMetricType:
		return "counter"
	case SetMetricType:
		return "set"
	case TimerMetricType:
		return "timer"
	case HistogramMetricType:
		return "histogram"
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

type DogStatsDMetric struct {
	Namespace     string
	Name          string
	MetricData    []byte
	Timestamp     time.Time
	MetricType    DogStatsDMetricType
	RawValue      string
	FloatValue    float64
	DurationValue time.Duration
	Extras        []string
	Tags          []string
	SampleRate    float64
}

func (d DogStatsDMetric) Data() []byte {
	return d.MetricData
}

func (d DogStatsDMetric) Type() DogStatsDMessageType {
	return MetricMessageType
}

type DogStatsDServiceCheck struct {
	data []byte
}

func (d DogStatsDServiceCheck) Data() []byte {
	return d.data
}

func (DogStatsDServiceCheck) Type() DogStatsDMessageType {
	return ServiceCheckMessageType
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
	if bytes.HasPrefix(buf, []byte("_sc{")) {
		return parseDogStatsDServiceCheckMessage(buf)
	}
	return parseDogStatsDMetricMessage(buf)
}

func parseDogStatsDMetricMessage(buf []byte) (DogStatsDMessage, error) {
	metric := DogStatsDMetric{
		Timestamp:  time.Now(),
		Tags:       make([]string, 0),
		SampleRate: 1.0,
	}

	// sample message: metric.name:value|type|@sample_rate|#tag1:value,tag2
	pieces := strings.Split(string(buf), "|")
	if len(pieces) < 2 {
		return nil, errors.New("invalid message: missing name, value, or type")
	}

	addrAndValue := strings.Split(pieces[0], ":")
	if len(addrAndValue) < 2 {
		return nil, errors.New("invalid message: missing name and value")
	}

	namespaceAndName := strings.SplitN(addrAndValue[0], ".", 2)
	if len(namespaceAndName) > 1 {
		metric.Namespace = namespaceAndName[0]
		metric.Name = namespaceAndName[1]
	} else {
		metric.Name = namespaceAndName[0]
	}

	metric.RawValue = addrAndValue[1]

	switch pieces[1] {
	case "c":
		metric.MetricType = CounterMetricType
	case "g":
		metric.MetricType = GaugeMetricType
	case "s":
		metric.MetricType = SetMetricType
	case "ms":
		metric.MetricType = TimerMetricType
	case "h":
		metric.MetricType = HistogramMetricType
	default:
		return nil, errors.New("invalid message: unknown metric type")
	}

	// all values are stored as a float
	floatValue, err := strconv.ParseFloat(metric.RawValue, 64)
	if err != nil {
		return nil, errors.New("invalid message: invalid value")
	}
	metric.FloatValue = floatValue

	if metric.MetricType == TimerMetricType {
		metric.DurationValue = time.Duration(metric.FloatValue) / time.Millisecond
	}

	// parse out sample rate, tags and any extras
	for _, piece := range pieces[2:] {
		if strings.HasPrefix(piece, "@") {
			sampleRate, err := strconv.ParseFloat(piece[1:], 64)
			if err != nil {
				return nil, errors.New("invalid sample rate")
			}
			metric.SampleRate = sampleRate
			continue
		}

		if strings.HasPrefix(piece, "#") {
			tags := strings.Split(piece[1:], ",")
			for i, _ := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			metric.Tags = append(metric.Tags, tags...)
			continue
		}

		metric.Extras = append(metric.Extras, piece)
	}

	return metric, nil
}

func parseDogStatsDEventMessage(_ []byte) (DogStatsDMessage, error) {
	return nil, errors.New("DogStatsD event messages not supported")
}

func parseDogStatsDServiceCheckMessage(_ []byte) (DogStatsDMessage, error) {
	return nil, errors.New("DogStatsD service check messages not supported")
}
