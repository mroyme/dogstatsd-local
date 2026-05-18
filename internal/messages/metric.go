package messages

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

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
		metric.MetricType = CountMetricType
	case "g":
		metric.MetricType = GaugeMetricType
	case "s":
		metric.MetricType = SetMetricType
	case "ms":
		metric.MetricType = TimerMetricType
	case "h":
		metric.MetricType = HistogramMetricType
	case "d":
		metric.MetricType = DistributionMetricType
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
		metric.DurationValue = time.Duration(metric.FloatValue * float64(time.Millisecond))
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
			for part := range strings.SplitSeq(piece[1:], ",") {
				metric.Tags = append(metric.Tags, strings.TrimSpace(part))
			}
			continue
		}

		metric.Extras = append(metric.Extras, piece)
	}

	return metric, nil
}
