package messages

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

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
		return "WARN"
	case ServiceCheckCritical:
		return "CRIT"
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
			for part := range strings.SplitSeq(piece[1:], ",") {
				sc.Tags = append(sc.Tags, strings.TrimSpace(part))
			}
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
