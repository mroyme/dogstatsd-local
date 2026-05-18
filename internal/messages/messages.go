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

type ServiceCheckStatus int

const (
	ServiceCheckOK       ServiceCheckStatus = 0
	ServiceCheckWarning  ServiceCheckStatus = 1
	ServiceCheckCritical ServiceCheckStatus = 2
	ServiceCheckUnknown  ServiceCheckStatus = 3
)

func (s ServiceCheckStatus) String() string {
	switch s {
	case ServiceCheckOK:
		return "OK"
	case ServiceCheckWarning:
		return "WARNING"
	case ServiceCheckCritical:
		return "CRITICAL"
	case ServiceCheckUnknown:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

type DogStatsDServiceCheck struct {
	Name      string
	Status    ServiceCheckStatus
	Message   string
	Tags      []string
	Timestamp time.Time
	Hostname  string
	RawData   []byte
}

func (d DogStatsDServiceCheck) Data() []byte {
	return d.RawData
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
	if bytes.HasPrefix(buf, []byte("_sc|")) {
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
			tags := strings.Split(piece[1:], ",")
			for i := range tags {
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

func parseDogStatsDServiceCheckMessage(buf []byte) (DogStatsDMessage, error) {
	// DogStatsD service check format:
	//   _sc|name|status|#tag1,tag2|d:timestamp|h:hostname|m:message
	// Per the spec, m: must be positioned last among metadata fields.
	// The message value may contain '|' characters, so we must extract it
	// before splitting the rest of the body on '|'.
	sc := DogStatsDServiceCheck{
		RawData:   buf,
		Timestamp: time.Now(),
		Tags:      make([]string, 0),
	}

	// Strip the "_sc|" prefix
	body := string(buf)
	if !strings.HasPrefix(body, "_sc|") {
		return nil, errors.New("invalid service check: missing _sc| prefix")
	}
	body = body[4:]

	// Extract m: message first since it must be last and can contain '|' characters.
	// We search for the last occurrence of "|m:" to find the message boundary.
	msgIdx := strings.LastIndex(body, "|m:")
	if msgIdx >= 0 {
		sc.Message = body[msgIdx+3:]
		body = body[:msgIdx]
	}

	pieces := strings.Split(body, "|")
	if len(pieces) < 2 {
		return nil, errors.New("invalid service check: missing name or status")
	}

	sc.Name = pieces[0]

	statusInt, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, errors.New("invalid service check: invalid status")
	}
	switch ServiceCheckStatus(statusInt) {
	case ServiceCheckOK, ServiceCheckWarning, ServiceCheckCritical, ServiceCheckUnknown:
		sc.Status = ServiceCheckStatus(statusInt)
	default:
		return nil, errors.New("invalid service check: status out of range (0-3)")
	}

	// Parse optional fields: tags, timestamp, hostname
	for _, piece := range pieces[2:] {
		if strings.HasPrefix(piece, "#") {
			tags := strings.Split(piece[1:], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			sc.Tags = append(sc.Tags, tags...)
			continue
		}
		if strings.HasPrefix(piece, "d:") {
			ts, err := strconv.ParseInt(piece[2:], 10, 64)
			if err == nil {
				sc.Timestamp = time.Unix(ts, 0)
			}
			continue
		}
		if strings.HasPrefix(piece, "h:") {
			sc.Hostname = piece[2:]
			continue
		}
	}

	return sc, nil
}
